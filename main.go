package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
)

const timeToSleep = 10000 * time.Millisecond

type GPU struct {
	Name   string `json:"name"`
	OnSale bool   `json:"onSale"`
}

var (
	apiURL string
	strs   = make(chan string)
	gpus   = make(chan GPU)
	emails = make(chan struct{})
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	go sender()
	go logger()
}

func logger() {
	for {
		select {
		case gpu := <-gpus:
			log.Printf("name: %s ---> onSale: %t\n", gpu.Name, gpu.OnSale)
		case str := <-strs:
			log.Println(str)
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
	m.SetHeader("Subject", "!!! GPU AVAILABLE !!!")
	// TODO: send the product url
	m.SetBody("text/plain", "!!! GPU AVAILABLE !!!")

	d := gomail.NewDialer("smtp.gmail.com", 587, emailSender, passwSender)
	d.TLSConfig = &tls.Config{ServerName: "smtp.gmail.com", InsecureSkipVerify: false}

	for {
		select {
		case <-emails:
			if err := d.DialAndSend(m); err != nil {
				strs <- err.Error()
				continue
			}
			strs <- "email sent successfully"
		}
	}
}

func main() {
	apiURL = os.Getenv("API_URL")
	var gpu GPU
	for {
		func() {
			resp, err := http.Get(apiURL)
			if err != nil {
				strs <- err.Error()
				return
			}
			defer resp.Body.Close()

			err = json.NewDecoder(resp.Body).Decode(&gpu)
			if err != nil {
				strs <- err.Error()
				return
			}

			if gpu.OnSale {
				emails <- struct{}{}
			}
			gpus <- gpu
		}()
		time.Sleep(timeToSleep)
	}
}
