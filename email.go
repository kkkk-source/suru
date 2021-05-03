package main

import (
	"crypto/tls"

	gomail "gopkg.in/mail.v2"
)

type emailService struct {
	sender    string
	senderKey string
	receiver  string
	smtHost   string
	smtPort   int
	dialer    *gomail.Dialer
	message   *gomail.Message
}

func NewEmailService(sender, senderKey, receiver, smtHost string, smtPort int) *Sender {
	dialer := gomail.NewDialer(smtHost, smtPort, sender, senderKey)
	dialer.TLSConfig = &tls.Config{ServerName: smtHost, InsecureSkipVerify: false}
	message := gomail.NewMessage()
	message.SetHeader("From", sender)

	return &emailService{
		sender:    sender,
		senderKey: senderKey,
		receiver:  receiver,
		smtHost:   smtHost,
		smtPort:   smtPort,
		dialer:    dialer,
		message:   message,
	}
}

func (e *emailService) Sender() string {
	return e.sender
}

func (e *emailService) Receiver() string {
	return e.receiver
}

func (e *emailService) Send(msgSubject, msgText string) error {
	e.message.SetHeader("To", e.receiver)
	e.message.SetHeader("Subject", msgSubject)
	e.message.SetBody("text/plain", msgText)
	if err := e.dialer.DialAndSend(e.message); err != nil {
		return err
	}
}
