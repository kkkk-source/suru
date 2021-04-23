package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
)

type GPU struct {
	Name   string `json:"name"`
	OnSale bool   `json:"onSale"`
}

const timeToSleep = 10 * time.Second

func main() {
	log.Println("starting suru ...")

	var (
		apiURL       = flag.String("apiUrl", "", "Specify the url to request to.")
		fromEmail    = flag.String("fromEmail", "", "Set the email messages sender.")
		fromEmailKey = flag.String("fromEmailKey", "", "Set the email messages sender password.")
		toEmail      = flag.String("toEmail", "", "Set the email messages receiver.")
		logFile      = flag.String("logFile", "out", "Specify a log file.")
		smtpHost     = flag.String("smtpHost", "smtp.gmail.com", "Specify the smtp host name to send messages.")
		smtpPort     = flag.Int("smtpPort", 587, "Specify the smtp port to send messages.")
	)
	flag.Parse()

	if *apiURL != "" {
		log.Printf("[!] apiUrl: %s.\n", *apiURL)
	}
	if *fromEmail != "" {
		log.Printf("[!] fromEmail: %s.\n", *fromEmail)
	}
	if *fromEmailKey != "" {
		log.Println("[!] fromEmailKey: fixed.")
	}
	if *toEmail != "" {
		log.Printf("[!] toEmail: %s.\n", *toEmail)
	}
	if *logFile != "" {
		log.Printf("[!] logFile: %s.\n", *logFile)
	}
	if *smtpHost != "" {
		log.Printf("[!] smtpHost: %s.\n", *smtpHost)
	}

	emailService := NewSender(
		*smtpHost,
		*smtpPort,
		*fromEmail,
		*fromEmailKey,
		*toEmail,
	)

	go emailService.Dispatcher()
	var gpu GPU

	for {
		func() {
			resp, err := http.Get(*apiURL)
			if err != nil {
				emailService.Messages("something was wrong")
				log.Println(err.Error())
				return
			}
			defer resp.Body.Close()

			err = json.NewDecoder(resp.Body).Decode(&gpu)
			if err != nil {
				emailService.Messages("something was wrong")
				log.Println(err.Error())
				return
			}

			if !gpu.OnSale {
				emailService.Messages("gpu in stock")
			}

			log.Printf("name: %s onSale: %t\n", gpu.Name, gpu.OnSale)
		}()
		time.Sleep(timeToSleep)
	}
}
