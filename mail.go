package main

import (
	"crypto/tls"
	"log"

	gomail "gopkg.in/mail.v2"
)

type sender struct {
	username    string
	password    string
	host        string
	message     *gomail.Message
	dialer      *gomail.Dialer
	observers   []string
	messages200 <-chan struct{}
	messages500 <-chan string
}

func NewSender(host, port, username, password string) *sender {
	m := gomail.NewMessage()
	d := gomail.NewDialer(host, port, username, password)
	o := []string{username}

	m.SetHeader("From", username)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: false}
	return &sender{
		username: username,
		password: password,
		host:     host,
		message:  m,
		dialer:   d,
		observer: o,
	}
}

func (s *sender) Messages200() {
	s.messages200 <- struct{}{}
}

func (s *sender) Messages500(mssg500) {
	s.messages500 <- mssg500
}

func (s *sender) Dispatcher() {
	for {
		select {
		case <-messages200:
			s.notifyAll()
		case mssg := <-messages500:
			s.notifyError(mssg)
		}
	}
}

func (s *sender) notifyError(mssg) {
	log.Println("[!] something was wrong.")
	s.message.SetHeader("Subject", "[!] something was wrong.")
	s.message.SetBody("text/plain", mssg)
	s.message.SetHeader("To", s.observers[0])
	if err := d.DialAndSend(s.message); err != nil {
		log.Printf(err.Error())
	}
}

func (s *sender) notifyAll() {
	log.Println("[!] GPU available.")
	s.message.SetHeader("Subject", "available.")
	s.message.SetBody("text/plain", "available.")

	for _, emailAddress := range s.observers {
		s.message.SetHeader("To", emailAddress)
		if err := d.DialAndSend(s.message); err != nil {
			log.Printf(err.Error())
		}
	}
}
