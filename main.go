package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
)

const timeToSleep = 200 * time.Millisecond

type GPU struct {
	Name      string `json:"name"`
	OnSale    bool   `json:"onSale"`
	Active    bool   `json:"active"`
	Orderable string `json:"orderable"`
}

var (
	apiURL string
	logs   = make(chan interface{})
	emails = make(chan struct{})
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiURL = os.Getenv("API_URL")
	go sender()
	go recorder()
}

func recorder() {
	for {
		select {
		case log := <-logs:
			fmt.Printf("%+v\n", log)
		}
	}
}

func sender() {
	emailReceiver := os.Getenv("TO_EMAIL")
	emailSender := os.Getenv("FROM_EMAIL")
	passwSender := os.Getenv("FROM_EMAIL_PASSWORD")

	m := gomail.NewMessage()
	m.SetHeader("From", emailSender)
	m.SetHeader("To", emailReceiver)
	m.SetHeader("Subject", "go: gpu available")
	m.SetBody("text/plain", "go: gpu available")

	d := gomail.NewDialer("smtp.gmail.com", 587, emailSender, passwSender)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{ServerName: "smtp.gmail.com", InsecureSkipVerify: false}

	for {
		select {
		case <-emails:
			if err := d.DialAndSend(m); err != nil {
				logs <- err.Error()
				continue
			}
			logs <- "Email Sent Successfully"
		}
	}
}

func main() {
	var gpu GPU
	for {
		func() {
			resp, err := http.Get(apiURL)
			if err != nil {
				logs <- err.Error()
				return
			}
			defer resp.Body.Close()

			err = json.NewDecoder(resp.Body).Decode(&gpu)
			if err != nil {
				logs <- err.Error()
				return
			}

			if gpu.OnSale {
				emails <- struct{}{}
			}
			logs <- gpu
		}()
		time.Sleep(timeToSleep)
	}
}
