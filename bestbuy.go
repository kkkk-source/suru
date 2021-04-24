package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	timeToSleep = 200 * time.Millisecond
	maxAttempts = 10
)

type item struct {
	Name   string `json:"name"`
	OnSale bool   `json:"onSale"`
}

type bestBuyService struct {
	apiURL string
	broker *MessageBroker
	logger *LoggerService
}

func NewBestBuyService(apiURL string, messageBroker *MessageBroker, logger *LoggerService) *bestBuyService {
	return &bestBuyService{
		broker: messageBroker,
		apiURL: apiURL,
		logger: logger,
	}
}

func (b *bestBuyService) Run() {
	go b.broker.Dispatcher()
	conscutiveNotOKExecutes := 0
	var item item
	for {
		func() {
			resp, err := http.Get(b.apiURL)
			if err != nil {
				b.broker.SendMessage(err.Error())
				b.logger.SendMessage(err.Error())
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				conscutiveNotOKExecutes++
			} else {
				conscutiveNotOKExecutes = 0
			}

			if conscutiveNotOKExecutes >= maxAttempts {
				msg := fmt.Sprintf("status code %d", resp.StatusCode)
				b.broker.SendMessage(msg)
				b.logger.SendMessage(msg)
				return
			}

			err = json.NewDecoder(resp.Body).Decode(&item)
			if err != nil {
				b.broker.SendMessage(err.Error())
				b.logger.SendMessage(err.Error())
				return
			}

			if item.OnSale {
				msg := "item in stock"
				b.broker.SendMessage(msg)
				b.logger.SendMessage(msg)
			}

			b.logger.SendMessage(fmt.Sprintf("%s - in stock: %t", item.Name, item.OnSale))
		}()
		time.Sleep(timeToSleep)
	}
}
