// Based on https://github.com/mailhog/mhsendmail/commit/90a69f19300e9f489a7be0b7fc6b74cdf66c88ec
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Configuration struct {
	Recipients []string
	Addr       string
	Sender     string
}

func getHostname() string {
	host, err := os.Hostname()
	if err != nil {
		return "localhost"
	}

	return host
}

func getUsername() string {
	username := "nobody"
	currentUser, err := user.Current()
	if err == nil && currentUser != nil && len(currentUser.Username) > 0 {
		username = currentUser.Username
	}

	return username
}

func main() {
	config := new(Configuration)
	app := kingpin.New(filepath.Base(os.Args[0]), "MailCage Sendmail")
	app.HelpFlag.Short('h')

	app.Flag("smtp-addr", "SMTP server address").
		Default("localhost:1025").
		Envar("SMTP_ADDR").
		StringVar(&config.Addr)

	app.Flag("smtp-sender", "Sender").
		Default(getUsername() + "@" + getHostname()).
		Envar("SMTP_SENDER").
		StringVar(&config.Sender)

	app.Flag("long-i", "Ignored. This flag exists for sendmail compatibility").
		Short('i').
		Default("").
		String()

	app.Flag("long-o", "Ignored. This flag exists for sendmail compatibility").
		Short('o').
		Default("").
		String()

	app.Flag("long-t", "Ignored. This flag exists for sendmail compatibility").
		Short('t').
		Default("").
		String()

	app.Arg("recipients", "Recipients").StringsVar(&config.Recipients)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "error parsing commandline arguments"))
		app.Usage(os.Args[1:])
		os.Exit(2)
	}

	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "error reading stdin"))
		os.Exit(11)
	}

	msg, err := mail.ReadMessage(bytes.NewReader(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "error parsing message body"))
		os.Exit(11)
	}

	if len(config.Recipients) == 0 {
		// We only need to parse the message to get a recipient if none where
		// provided on the command line.
		config.Recipients = append(config.Recipients, msg.Header.Get("To"))
	}

	if err = smtp.SendMail(config.Addr, nil, config.Sender, config.Recipients, body); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "error sending mail"))
		os.Exit(11)
	}
}
