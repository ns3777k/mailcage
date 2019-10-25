package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReply(t *testing.T) {
	testCases := []struct {
		Name     string
		Reply    *Reply
		Expected []string
	}{
		{
			Name:     "Empty content 200",
			Reply:    &Reply{Status: 200, lines: []string{}, Done: nil},
			Expected: []string{"200\n"},
		},
		{
			Name:     "1 line reply 200",
			Reply:    &Reply{Status: 200, lines: []string{"Ok"}, Done: nil},
			Expected: []string{"200 Ok\r\n"},
		},
		{
			Name:     "2 lines reply 200",
			Reply:    &Reply{Status: 200, lines: []string{"Ok", "Still ok!"}, Done: nil},
			Expected: []string{"200-Ok\r\n", "200 Still ok!\r\n"},
		},
		{
			Name:     "3 lines reply 200",
			Reply:    &Reply{Status: 200, lines: []string{"Ok", "Still ok!", "OINK!"}, Done: nil},
			Expected: []string{"200-Ok\r\n", "200-Still ok!\r\n", "200 OINK!\r\n"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			l := tc.Reply.Lines()

			assert.Equal(t, tc.Expected, l)
		})
	}
}

func TestBuiltInReplies(t *testing.T) {
	testCases := []struct {
		Name           string
		Reply          *Reply
		ExpectedStatus int
		ExpectedLines  []string
	}{
		{
			Name:           "ReplyIdent is correct",
			Reply:          ReplyIdent("mailcage"),
			ExpectedStatus: 220,
			ExpectedLines:  []string{"mailcage"},
		},
		{
			Name:           "ReplyBye is correct",
			Reply:          ReplyBye(),
			ExpectedStatus: 221,
			ExpectedLines:  []string{"Bye"},
		},
		{
			Name:           "ReplyAuthOk is correct",
			Reply:          ReplyAuthOk(),
			ExpectedStatus: 235,
			ExpectedLines:  []string{"Authentication successful"},
		},
		{
			Name:           "ReplyOk without arguments is correct",
			Reply:          ReplyOk(),
			ExpectedStatus: 250,
			ExpectedLines:  []string{"Ok"},
		},
		{
			Name:           "ReplyOk with an argument is correct",
			Reply:          ReplyOk("mailcage"),
			ExpectedStatus: 250,
			ExpectedLines:  []string{"mailcage"},
		},
		{
			Name:           "ReplyOk with multiple arguments is correct",
			Reply:          ReplyOk("mailcage", "test"),
			ExpectedStatus: 250,
			ExpectedLines:  []string{"mailcage", "test"},
		},
		{
			Name:           "ReplySenderOk is correct",
			Reply:          ReplySenderOk("test"),
			ExpectedStatus: 250,
			ExpectedLines:  []string{"Sender test ok"},
		},
		{
			Name:           "ReplyRecipientOk is correct",
			Reply:          ReplyRecipientOk("test"),
			ExpectedStatus: 250,
			ExpectedLines:  []string{"Recipient test ok"},
		},
		{
			Name:           "ReplyAuthResponse is correct",
			Reply:          ReplyAuthResponse("test"),
			ExpectedStatus: 334,
			ExpectedLines:  []string{"test"},
		},
		{
			Name:           "ReplyDataResponse is correct",
			Reply:          ReplyDataResponse(),
			ExpectedStatus: 354,
			ExpectedLines:  []string{"End data with <CR><LF>.<CR><LF>"},
		},
		{
			Name:           "ReplyStorageFailed is correct",
			Reply:          ReplyStorageFailed("test"),
			ExpectedStatus: 452,
			ExpectedLines:  []string{"test"},
		},
		{
			Name:           "ReplyUnrecognisedCommand is correct",
			Reply:          ReplyUnrecognisedCommand(),
			ExpectedStatus: 500,
			ExpectedLines:  []string{"Unrecognized command"},
		},
		{
			Name:           "ReplyUnsupportedAuth is correct",
			Reply:          ReplyUnsupportedAuth(),
			ExpectedStatus: 504,
			ExpectedLines:  []string{"Unsupported authentication mechanism"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.ExpectedStatus, tc.Reply.Status)
			assert.Equal(t, tc.ExpectedLines, tc.Reply.lines)
		})
	}
}
