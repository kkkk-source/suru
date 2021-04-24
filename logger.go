package main

import (
	"log"
)

type LoggerService struct {
	logs chan string
}

func NewLoggerService() *LoggerService {
	logs := make(chan string)
	return &LoggerService{
		logs: logs,
	}
}

func (l *LoggerService) SendMessage(msg string) {
	l.logs <- msg
}

func (l *LoggerService) Logger() {
	for {
		select {
		case msg := <-l.logs:
			log.Printf(msg)
		}
	}
}
