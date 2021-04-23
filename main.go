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
		logFile      = flag.String("logFile", "out", "Specify a log file.")
		smtpHost     = flag.String("smtpHost", "smtp.gmail.com", "Specify the smtp host name to send messages.")
		smtpPort     = flag.Int("smtpPort", 587, "Specify the smtp port to send messages.")
	)
	flag.Parse()

	if *apiURL != "" {
		log.Printf("apiUrl: %s.\n", *apiURL)
	}
	if *logFile != "" {
		log.Printf("logFile: %s.\n", *logFile)
	}

	messagBrokrService := NewMessageBrokerService(
		*smtpHost,
		*smtpPort,
		*fromEmail,
		*fromEmailKey,
		*toEmail,
	)
	bestBuyService := NewBestBuyService(
		*apiURL,
		messagBrokrService,
	)
	bestBuyService.Run()
}
