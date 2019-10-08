// Based on https://github.com/mailhog/MailHog-Server/commit/e7f979c845de6e9404a91875861a69cca2847f58

package smtp

import (
	"crypto/tls"
	"github.com/ns3777k/mailcage/storage"
	"github.com/pkg/errors"
	gosmtp "net/smtp"
	"strconv"
)

type OutgoingServer struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Mechanism string `json:"mechanism"`
}

type Mailer struct {
	hostname string
	servers  map[string]*OutgoingServer
}

func NewMailer(hostname string, servers map[string]*OutgoingServer) *Mailer {
	return &Mailer{
		hostname: hostname,
		servers:  servers,
	}
}

func (m *Mailer) createAuth(outgoingServer *OutgoingServer) (gosmtp.Auth, error) {
	var auth gosmtp.Auth
	var err error

	if outgoingServer.Username == "" && outgoingServer.Password == "" {
		return auth, err
	}

	switch outgoingServer.Mechanism {
	case "CRAMMD5":
		auth = gosmtp.CRAMMD5Auth(outgoingServer.Username, outgoingServer.Password)
	case "PLAIN":
		auth = gosmtp.PlainAuth("", outgoingServer.Username, outgoingServer.Password, outgoingServer.Host)
	default:
		err = errors.New("invalid authentication mechanism")
	}

	return auth, err
}

func (m *Mailer) send(outgoingServer *OutgoingServer, message []byte) error {
	client, err := gosmtp.Dial(outgoingServer.Host + ":" + strconv.Itoa(outgoingServer.Port))
	if err != nil {
		return err
	}

	defer client.Close()

	if err := client.Hello(m.hostname); err != nil {
		return err
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: outgoingServer.Host}
		if err := client.StartTLS(config); err != nil {
			return err
		}
	}

	auth, err := m.createAuth(outgoingServer)
	if err != nil {
		return err
	}

	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail("nobody@" + m.hostname); err != nil {
		return err
	}

	if err = client.Rcpt(outgoingServer.Email); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write(message); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func (m *Mailer) Send(server string, message *storage.Message) error {
	if _, ok := m.servers[server]; !ok {
		return errors.New("outgoing server not found")
	}

	bytes := make([]byte, 0)
	for h, l := range message.Content.Headers {
		for _, v := range l {
			bytes = append(bytes, []byte(h+": "+v+"\r\n")...)
		}
	}
	bytes = append(bytes, []byte("\r\n"+message.Content.Body)...)

	return m.send(m.servers[server], bytes)
}
