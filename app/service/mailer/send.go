package mailer

import (
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
)

func (mailer *Mailer) Send(e *email.Email) error {
	addr := mailer.Host + ":" + mailer.Port
	return e.Send(addr, mailer.Auth)
}

func (mailer *Mailer) SendText(sm *SmtpMessage) error {
	mess := []string{}
	addr := mailer.Host + ":" + mailer.Port

	mess = append(mess, "To: "+sm.To[0])
	mess = append(mess, "From: "+sm.From)
	mess = append(mess, "Subject: "+sm.Subject)
	mess = append(mess, "Reply-To: "+sm.ReplyTo)
	mess = append(mess, "MIME-version: 1.0")
	mess = append(mess, "Content-Type: text/plain; charset=utf-8")
	mess = append(mess, "")
	mess = append(mess, sm.Body)

	fullMsg := []byte(strings.Join(mess, "\r\n"))
	if err := smtp.SendMail(addr, mailer.Auth, sm.From, sm.To, fullMsg); err != nil {
		return err
	}
	return nil
}
