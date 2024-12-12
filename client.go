package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:      10,
			IdleConnTimeout:   30 * time.Second,
			DisableKeepAlives: false,
		},
	}

	url := "http://localhost:8080"

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
		time.Sleep(1 * time.Second)
	}
}
