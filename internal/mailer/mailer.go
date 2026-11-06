package mailer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/a-h/templ"
	"github.com/wneessen/go-mail"
)

type Config struct {
	Host     string
	Username string
	Password string
	BaseMail string
}

type Mailer struct {
	client   *mail.Client
	baseMail string
}

func New(cfg Config) Mailer {
	client, err := mail.NewClient(
		cfg.Host,
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(cfg.Username),
		mail.WithPassword(cfg.Password),
	)
	if err != nil {
		log.Fatalf("failed to create mail client: %v", err)
	}

	return Mailer{
		client:   client,
		baseMail: cfg.BaseMail,
	}
}

func (m Mailer) Send(recipient, subject string, template templ.Component) error {
	message := mail.NewMsg()
	err := message.From(m.baseMail)
	if err != nil {
		return fmt.Errorf("failed to set FROM address: %s", err)
	}
	err = message.To(recipient)
	if err != nil {
		return fmt.Errorf("failed to set TO address: %s", err)
	}
	html, err := templ.ToGoHTML(context.Background(), template)
	if err != nil {
		return fmt.Errorf("failed to convert templ.Component to string: %s", err)
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, string(html))

	for range 3 {
		err = m.client.DialAndSend(message)
		if err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return err
}
