package data

import (
	"fmt"
	"net/smtp"
)

// Mail ...
type Mail struct {
	sender  string
	to      string
	subject string
	body    string
}

// NewMail ...
func NewMail(to string, subject string, body string) *Mail {
	emailAddress := "<<YOUR EMAIL ADDRESS >>"
	return &Mail{emailAddress, to, subject, body}
}

func (mail *Mail) buildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", "Amalhanaja Account Activation")
	message += fmt.Sprintf("To: %s\r\n", mail.to)

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

//SendMessage ...
func (mail *Mail) SendMessage() error {
	host := "smtp.gmail.com"
	port := 587
	server := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", mail.sender, "<<YOUR PASSWORD>>", host)
	return smtp.SendMail(server, auth, mail.sender, []string{mail.to}, []byte(mail.buildMessage()))
}
