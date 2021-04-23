package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GPU struct {
	Name      string `json:"name"`
	OnSale    bool   `json:"onSale"`
	Active    bool   `json:"active"`
	Orderable string `json:"orderable"`
}

const (
	timeToSleep = 200 * time.Millisecond
	apiKey      = "avS2W2GXy5rTERtEXAFdOjKO"
	url         = "https://api.bestbuy.com/v1/products/6439402.json?show=name,onSale,active,inStoreAvailabilityUpdateDate,orderable&apiKey=" + apiKey
)

var logs = make(chan interface{})

func init() {
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
			resp, err := http.Get(url)
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
