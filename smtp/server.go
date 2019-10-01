package smtp

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/ns3777k/mailcage/smtp/protocol"
	"github.com/ns3777k/mailcage/storage"
	"github.com/rs/zerolog"
)

type Server struct {
	logger  zerolog.Logger
	storage storage.Storage
	opts    *ServerOptions
}

type ServerOptions struct {
	ListenAddr string
	Hostname   string
}

func NewServer(opts *ServerOptions, logger zerolog.Logger, storage storage.Storage) *Server {
	return &Server{logger: logger, opts: opts, storage: storage}
}

func (s *Server) serve(remoteAddress string, conn io.ReadWriteCloser) {
	defer conn.Close()

	protoDebugLogger := s.logger.With().Str("smtp", "protocol").Logger()

	proto := protocol.NewProtocol()
	proto.Hostname = s.opts.Hostname
	session := &Session{
		conn:          conn,
		proto:         proto,
		storage:       s.storage,
		remoteAddress: remoteAddress,
		reader:        io.Reader(conn),
		writer:        io.Writer(conn),
		logger:        s.logger.With().Str("smtp", "session").Logger(),
	}
	proto.LogHandler = func(message string, args ...interface{}) {
		protoDebugLogger.Debug().Msg(fmt.Sprintf(message, args...))
	}
	proto.MessageReceivedHandler = session.acceptMessage
	proto.ValidateSenderHandler = session.validateSender
	proto.ValidateRecipientHandler = session.validateRecipient
	proto.ValidateAuthenticationHandler = session.validateAuthentication
	proto.GetAuthenticationMechanismsHandler = func() []string {
		return []string{"PLAIN"}
	}

	s.logger.Info().Str("ip", remoteAddress).Msg("starting session")
	session.Write(proto.Start())
	for session.Read() {
	}
	s.logger.Info().Str("ip", remoteAddress).Msg("ending session")
}

func (s *Server) Run(ctx context.Context) error {
	var connections sync.WaitGroup
	s.logger.Info().Msg("starting to listen on " + s.opts.ListenAddr)

	addr, err := net.ResolveTCPAddr("tcp", s.opts.ListenAddr)
	if err != nil {
		return err
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	defer func() {
		if err := ln.Close(); err != nil {
			s.logger.Err(err).Msg("failed to close tcp listener")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			ln.Close()
			connections.Wait()
			return nil
		default:
			ln.SetDeadline(time.Now().Add(1e9)) //nolint:errcheck

			conn, err := ln.AcceptTCP()
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			if err != nil {
				s.logger.Err(err).Msg("error accepting connection")
			}

			remoteAddr := conn.RemoteAddr().String()
			s.logger.Debug().Str("ip", remoteAddr).Msg("incoming connection")

			connections.Add(1)

			go func() {
				s.serve(remoteAddr, io.ReadWriteCloser(conn))
				connections.Done()
			}()
		}
	}
}
