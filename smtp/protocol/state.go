// Based on https://github.com/mailhog/smtp/commit/f5aa8ab40b86f969f7d18186948e5eddcb93c5b5
package protocol

// State represents the state of an SMTP conversation
type State int

// SMTP message conversation states
const (
	INVALID   = State(-1)
	ESTABLISH = State(iota)
	AUTHPLAIN
	AUTHLOGIN
	AUTHLOGIN2
	AUTHCRAMMD5
	MAIL
	RCPT
	DATA
	DONE
)

// stateMap provides string representations of SMTP conversation states
//nolint:gochecknoglobals
var stateMap = map[State]string{
	INVALID:     "INVALID",
	ESTABLISH:   "ESTABLISH",
	AUTHPLAIN:   "AUTHPLAIN",
	AUTHLOGIN:   "AUTHLOGIN",
	AUTHLOGIN2:  "AUTHLOGIN2",
	AUTHCRAMMD5: "AUTHCRAMMD5",
	MAIL:        "MAIL",
	RCPT:        "RCPT",
	DATA:        "DATA",
	DONE:        "DONE",
}
