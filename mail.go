package main

import (
	"crypto/tls"
	"log"

	gomail "gopkg.in/mail.v2"
)

type sender struct {
	username string
	password string
	to       string
	host     string
	port     int
	message  *gomail.Message
	dialer   *gomail.Dialer
	messages chan string
}

func NewSender(host string, port int, username, password, to string) *sender {
	m := gomail.NewMessage()
	d := gomail.NewDialer(host, port, username, password)
	c := make(chan string)

	m.SetHeader("From", username)
	d.TLSConfig = &tls.Config{ServerName: host, InsecureSkipVerify: false}
	return &sender{
		username: username,
		password: password,
		host:     host,
		port:     port,
		to:       to,
		messages: c,
		message:  m,
		dialer:   d,
	}
}

func (s *sender) Messages(mssg string) {
	s.messages <- mssg
}

func (s *sender) Dispatcher() {
	for {
		select {
		case mssg := <-s.messages:
			s.notify(mssg)
		}
	}
}

func (s *sender) notify(mssg string) {
	log.Println("[!] sending an email.")
	s.message.SetHeader("To", s.to)
	s.message.SetHeader("Subject", mssg)
	s.message.SetBody("text/plain", mssg)
	if err := s.dialer.DialAndSend(s.message); err != nil {
		log.Printf(err.Error())
	}
}
