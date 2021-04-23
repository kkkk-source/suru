package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
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
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiURL = os.Getenv("apiUrl")
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
				// notifyme
			}

			logs <- gpu
		}()
		time.Sleep(timeToSleep)
	}
}
