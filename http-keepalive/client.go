package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {

	client := &http.Client{
		Timeout: 50 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
	}

	url := "http://localhost:49755"

	for i := 1; i <= 3; i++ {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Request %d failed: %v\n", i, err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Failed to read response %d: %v\n", i, err)
			continue
		}

		fmt.Printf("Response %d: %s\n", i, body)
		time.Sleep(2 * time.Second)
	}
}
