package balancer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	target string
	rate   float32

	stop chan struct{}
	mux  sync.RWMutex
}

func NewClient(rate float32, target_addr *url.URL) *Client {
	return &Client{
		rate:   rate,
		target: fmt.Sprintf("http://%s", target_addr.Host),
	}
}

func (client *Client) Run(period time.Duration) {
	t := time.NewTicker(period)
	for {
		select {
		case <-t.C:
			log.Printf("Client: Sending request to %s...\n", client.target)
			response, err := http.Get(client.target)
			if err != nil {
				log.Printf("Client Error: %s\n", err)
			}
			res, _ := ioutil.ReadAll(response.Body)
			log.Printf("Client: Received response from %s - %s\n", response.Request.URL, res)
		case <-client.stop:
			t.Stop()
			return
		}
	}
}

func (client *Client) Start() {
	client.mux.Lock()
	defer client.mux.Unlock()

	period := time.Duration(float32(time.Second) / client.rate)

	if client.stop == nil {
		client.stop = make(chan struct{})
		go client.Run(period)

		log.Printf("Client: Requesting %s every %s\n", client.target, period)
	}
}

func (client *Client) Stop() {
	client.mux.Lock()
	defer client.mux.Unlock()

	if client.stop != nil {
		client.stop <- struct{}{}
		close(client.stop)
		client.stop = nil
		log.Printf("Client: Stopped")
	}
}
