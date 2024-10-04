package mailer

import (
	"net/smtp"

	"github.com/jordan-wright/email"

	"SnipAndNeat/app/config"
)

//smtp_server=infra-01.postfix.qa.lamoda.tech for tests
//smtp_server=mailgw.hosts.lamoda.ru

type EmailSender interface {
	SendText(sm *SmtpMessage) error
	Send(e *email.Email) error
}

type SmtpMessage struct {
	To      []string
	From    string
	Body    string
	Subject string
	ReplyTo string
}

type Mailer struct {
	Host string
	Port string
	Auth smtp.Auth
}

func NewMailer(cfg *config.Config) EmailSender {
	getMailerAuth := func() smtp.Auth {
		if cfg.EmailServer.UseAuth {
			return smtp.PlainAuth(
				cfg.EmailServer.Identity,
				cfg.EmailServer.Username,
				cfg.EmailServer.Password,
				cfg.EmailServer.Host)
		}
		return nil
	}

	return &Mailer{
		Host: cfg.EmailServer.Host,
		Port: cfg.EmailServer.Port,
		Auth: getMailerAuth(),
	}
}
