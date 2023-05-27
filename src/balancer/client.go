package balancer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	target string
	rate   float32
}

func NewClient(rate float32, target_addr *url.URL) *Client {
	return &Client{
		rate:   rate,
		target: fmt.Sprintf("http://%s", target_addr.Host),
	}
}

func (client *Client) Start() {
	period := time.Duration(float32(time.Second) / client.rate)
	log.Printf("Client: Requesting %s every %s\n", client.target, period)
	for {
		log.Printf("Client: Sending request to %s...\n", client.target)
		response, err := http.Get(client.target)
		if err != nil {
			log.Printf("Client Error: %s\n", err)
			return
		}
		res, _ := ioutil.ReadAll(response.Body)
		log.Printf("Client: Received response from %s - %s\n", response.Request.URL, res)
		time.Sleep(period)
	}
}
