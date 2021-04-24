package main

import (
	"flag"
	"log"
)

func main() {
	log.Println("starting suru ...")

	var (
		apiURL       = flag.String("apiUrl", "", "Specify the url to request to.")
		fromEmail    = flag.String("fromEmail", "", "Set the email sender.")
		fromEmailKey = flag.String("fromEmailKey", "", "Set the email sender password.")
		toEmail      = flag.String("toEmail", "", "Set the email receiver.")
		smtpHost     = flag.String("smtpHost", "smtp.gmail.com", "Specify the smtp host name to send messages.")
		smtpPort     = flag.Int("smtpPort", 587, "Specify the smtp port to send messages.")
	)
	flag.Parse()

	if *apiURL != "" {
		log.Printf("apiUrl: %s.\n", *apiURL)
	}

	loggerService := NewLoggerService()
	go loggerService.Logger()

	messagBrokrService := NewMessageBrokerService(
		*smtpHost,
		*smtpPort,
		*fromEmail,
		*fromEmailKey,
		*toEmail,
		loggerService,
	)
	bestBuyService := NewBestBuyService(
		*apiURL,
		messagBrokrService,
		loggerService,
	)
	bestBuyService.Run()
}
