// Based on https://github.com/mailhog/MailHog-Server/commit/1a4b117f17c5cdb879b81ebbbdb3c976df5ce37b
// http://www.rfc-editor.org/rfc/rfc5321.txt
package smtp

import (
	"io"
	"strings"

	"github.com/ns3777k/mailcage/smtp/protocol"
	"github.com/ns3777k/mailcage/storage"
	"github.com/rs/zerolog"
)

// Session represents a SMTP session using net.TCPConn
type Session struct {
	conn          io.ReadWriteCloser
	proto         *protocol.Protocol
	storage       storage.Storage
	remoteAddress string
	line          string
	logger        zerolog.Logger

	reader io.Reader
	writer io.Writer
}

func (s *Session) validateAuthentication(mechanism string, args ...string) (errorReply *protocol.Reply, ok bool) {
	return nil, true
}

func (s *Session) validateRecipient(to string) bool {
	return true
}

func (s *Session) validateSender(from string) bool {
	return true
}

func (s *Session) acceptMessage(rawMessage *protocol.Message) (string, error) {
	message := RawMessageToStorage(rawMessage, s.proto.Hostname)

	s.logger.Debug().Interface("message", rawMessage).Msg("storing smtp message to database")
	s.logger.Info().Str("id", message.ID).Str("from", message.From.Mailbox).Msg("storing message")

	return s.storage.Store(message)
}

// Read reads from the underlying net.TCPConn
func (s *Session) Read() bool {
	buf := make([]byte, 1024)
	n, err := s.reader.Read(buf)

	if n == 0 {
		s.logger.Debug().Msg("connection closed by remote host")
		io.Closer(s.conn).Close() // not sure this is necessary?
		return false
	}

	if err != nil {
		s.logger.Err(err).Msg("error reading from socket")
		return false
	}

	text := string(buf[0:n])
	logText := strings.Replace(text, "\n", "\\n", -1)
	logText = strings.Replace(logText, "\r", "\\r", -1)
	s.logger.Debug().Int("bytes", n).Str("content", logText).Msg("received bytes")

	s.line += text

	for strings.Contains(s.line, "\r\n") {
		line, reply := s.proto.Parse(s.line)
		s.line = line

		if reply != nil {
			s.Write(reply)
			if reply.Status == 221 {
				io.Closer(s.conn).Close()
				return false
			}
		}
	}

	return true
}

// Write writes a reply to the underlying net.TCPConn
func (s *Session) Write(reply *protocol.Reply) {
	lines := reply.Lines()
	for _, l := range lines {
		logText := strings.Replace(l, "\n", "\\n", -1)
		logText = strings.Replace(logText, "\r", "\\r", -1)
		s.logger.Debug().Int("bytes", len(l)).Str("content", logText).Msg("sent bytes")
		s.writer.Write([]byte(l)) //nolint:errcheck
	}
}
