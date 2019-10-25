package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocol_InitialState(t *testing.T) {
	p := NewProtocol()

	assert.NotNil(t, p)
	assert.NotNil(t, p.Message)
	assert.IsType(t, &Protocol{}, p)
	assert.IsType(t, &Message{}, p.Message)
	assert.Equal(t, "mailcage.example", p.Hostname)
	assert.Equal(t, "ESMTP MailCage", p.Ident)
	assert.Equal(t, INVALID, p.State)
}

func TestProtocol_LoggingHandler(t *testing.T) {
	p := NewProtocol()
	handlerCalled := false
	p.LogHandler = func(message string, args ...interface{}) {
		handlerCalled = true

		assert.Equal(t, "[PROTO: %s] Test message %s %s", message)
		assert.Equal(t, []interface{}{"INVALID", "test arg 1", "test arg 2"}, args)
	}

	p.logf("Test message %s %s", "test arg 1", "test arg 2")
	assert.True(t, handlerCalled)
}

func TestProtocol_StartState(t *testing.T) {
	p := NewProtocol()
	reply := p.Start()

	assert.Equal(t, ESTABLISH, p.State)
	assert.NotNil(t, reply)
	assert.Equal(t, 220, reply.Status)
	assert.Equal(t, reply.Lines(), []string{"220 mailcage.example ESMTP MailCage\r\n"})
}

func TestProtocol_ChangeIdentInReply(t *testing.T) {
	p := NewProtocol()
	p.Ident = "NEW_IDENT"
	reply := p.Start()

	assert.NotNil(t, reply)
	assert.Equal(t, 220, reply.Status)
	assert.Equal(t, reply.Lines(), []string{"220 mailcage.example NEW_IDENT\r\n"})
}

func TestProtocol_ChangeHostnameInReply(t *testing.T) {
	p := NewProtocol()
	p.Hostname = "mailcage.host"
	reply := p.Start()

	assert.NotNil(t, reply)
	assert.Equal(t, 220, reply.Status)
	assert.Equal(t, reply.Lines(), []string{"220 mailcage.host ESMTP MailCage\r\n"})
}

func TestProtocol_ProcessCommand_InvalidCommand(t *testing.T) {
	p := NewProtocol()
	reply := p.ProcessCommand("INVALID CMD")

	assert.NotNil(t, reply)
	assert.Equal(t, 500, reply.Status)
	assert.Equal(t, reply.Lines(), []string{"500 Unrecognized command\r\n"})
}

func TestProtocol_ProcessCommand_HeloCommand(t *testing.T) {
	p := NewProtocol()
	p.Start()

	reply := p.ProcessCommand("HELO localhost")
	assert.Equal(t, MAIL, p.State)
	assert.NotNil(t, reply)
	assert.Equal(t, 250, reply.Status)
	assert.Equal(t, reply.Lines(), []string{"250 Hello localhost\r\n"})
}

func TestProtocol_ProcessCommand_InvalidCommandAfterHelo(t *testing.T) {
	p := NewProtocol()
	p.Start()
	p.ProcessCommand("HELO localhost")

	reply := p.ProcessCommand("INVALID CMD")
	assert.NotNil(t, reply)
	assert.Equal(t, 500, reply.Status)
	assert.Equal(t, reply.Lines(), []string{"500 Unrecognized command\r\n"})
}
