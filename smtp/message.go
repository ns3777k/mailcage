// Based on https://github.com/mailhog/data/commit/024d554958b5bea5db220bfd84922a584d878ded
package smtp

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"mime"
	"strings"
	"time"

	"github.com/ns3777k/mailcage/smtp/protocol"
	"github.com/ns3777k/mailcage/storage"
)

// GenerateMessageID generates a new message ID
func GenerateMessageID(hostname string) (string, error) {
	size := 32

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	rs := base64.URLEncoding.EncodeToString(rb)

	return rs + "@" + hostname, nil
}

func isMIME(content *storage.Content) bool {
	header, ok := content.Headers["Content-Type"]
	if !ok {
		return false
	}
	return strings.HasPrefix(header[0], "multipart/")
}

func parseMIMEBody(content *storage.Content) *storage.MIMEBody {
	var parts []*storage.Content

	if hdr, ok := content.Headers["Content-Type"]; ok {
		if len(hdr) > 0 {
			boundary := extractBoundary(hdr[0])
			var p []string
			if len(boundary) > 0 {
				p = strings.Split(content.Body, "--"+boundary)
			}

			for _, s := range p {
				if len(s) > 0 {
					part := ContentFromString(strings.Trim(s, "\r\n"))
					if isMIME(part) {
						part.MIME = parseMIMEBody(part)
					}
					parts = append(parts, part)
				}
			}
		}
	}

	return &storage.MIMEBody{Parts: parts}
}

func RawMessageToStorage(rawMessage *protocol.Message, hostname string) *storage.Message {
	arr := make([]*storage.Path, len(rawMessage.To))

	for i, path := range rawMessage.To {
		arr[i] = PathFromString(path)
	}

	id, _ := GenerateMessageID(hostname)
	msg := &storage.Message{
		ID:        id,
		From:      PathFromString(rawMessage.From),
		To:        arr,
		Content:   ContentFromString(rawMessage.Data),
		CreatedAt: time.Now(),
		Raw: &storage.RawMessage{
			From: rawMessage.From,
			To:   rawMessage.To,
			Data: rawMessage.Data,
			Helo: rawMessage.Helo,
		},
	}

	if isMIME(msg.Content) {
		msg.MIME = parseMIMEBody(msg.Content)
	}

	// find headers
	var hasMessageID bool
	var receivedHeaderName string
	var returnPathHeaderName string

	for k := range msg.Content.Headers {
		if strings.ToLower(k) == "message-id" {
			hasMessageID = true
			continue
		}
		if strings.ToLower(k) == "received" {
			receivedHeaderName = k
			continue
		}
		if strings.ToLower(k) == "return-path" {
			returnPathHeaderName = k
			continue
		}
	}

	if !hasMessageID {
		msg.Content.Headers["Message-ID"] = []string{id}
	}

	from := fmt.Sprintf(
		"from %s by %s (MailCage)\r\n          id %s; %s",
		rawMessage.Helo,
		hostname,
		id,
		time.Now().Format(time.RFC1123Z),
	)

	if len(receivedHeaderName) > 0 {
		msg.Content.Headers[receivedHeaderName] = append(msg.Content.Headers[receivedHeaderName], from)
	} else {
		msg.Content.Headers["Received"] = []string{from}
	}

	if len(returnPathHeaderName) > 0 {
		msg.Content.Headers[returnPathHeaderName] = append(msg.Content.Headers[returnPathHeaderName], "<"+rawMessage.From+">")
	} else {
		msg.Content.Headers["Return-Path"] = []string{"<" + rawMessage.From + ">"}
	}

	return msg
}

// PathFromString parses a forward-path or reverse-path into its parts
func PathFromString(path string) *storage.Path {
	var relays []string
	email := path
	if strings.Contains(path, ":") {
		x := strings.SplitN(path, ":", 2)
		r, e := x[0], x[1]
		email = e
		relays = strings.Split(r, ",")
	}
	mailbox, domain := "", ""
	if strings.Contains(email, "@") {
		x := strings.SplitN(email, "@", 2)
		mailbox, domain = x[0], x[1]
	} else {
		mailbox = email
	}

	return &storage.Path{
		Relays:  relays,
		Mailbox: mailbox,
		Domain:  domain,
		Params:  "", // FIXME?
	}
}

// ContentFromString parses SMTP content into separate headers and body
func ContentFromString(data string) *storage.Content {
	x := strings.SplitN(data, "\r\n\r\n", 2)
	h := make(map[string][]string)

	// FIXME this fails if the message content has no headers - specifically,
	// if it doesn't contain \r\n\r\n

	if len(x) == 2 {
		headers, body := x[0], x[1]
		hdrs := strings.Split(headers, "\r\n")
		var lastHdr = ""
		for _, hdr := range hdrs {
			if lastHdr != "" && (strings.HasPrefix(hdr, " ") || strings.HasPrefix(hdr, "\t")) {
				h[lastHdr][len(h[lastHdr])-1] = h[lastHdr][len(h[lastHdr])-1] + hdr
			} else if strings.Contains(hdr, ": ") {
				y := strings.SplitN(hdr, ": ", 2)
				key, value := y[0], y[1]
				// TODO multiple header fields
				h[key] = []string{value}
				lastHdr = key
			}
		}
		return &storage.Content{
			Size:    len(data),
			Headers: h,
			Body:    body,
		}
	}
	return &storage.Content{
		Size:    len(data),
		Headers: h,
		Body:    x[0],
	}
}

// extractBoundary extract boundary string in contentType.
// It returns empty string if no valid boundary found
func extractBoundary(contentType string) string {
	_, params, err := mime.ParseMediaType(contentType)
	if err == nil {
		return params["boundary"]
	}
	return ""
}
