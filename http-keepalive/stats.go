package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tcnksm/go-httpstat"
)

const (
	URL           = "http://35.243.94.80/"
	SleepTime     = 1000 * time.Millisecond
	ClientTimeout = 50 * time.Second
	TryCount      = 100
)

func main() {
	client := &http.Client{
		Timeout: ClientTimeout,
		Transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
	}

	fmt.Printf("Sleep duration: %d ms\n", SleepTime.Milliseconds())

	for i := 1; i <= TryCount; i++ {
		var result httpstat.Result

		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			fmt.Printf("Request %d creation failed: %v\n", i, err)
			continue
		}

		ctx := httpstat.WithHTTPStat(req.Context(), &result)
		req = req.WithContext(ctx)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Request %d failed: %v\n", i, err)
			continue
		}

		_, err = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Failed to read response %d: %v\n", i, err)
			continue
		}

		result.End(time.Now())

		fmt.Printf("Request %d:\n", i)
		fmt.Printf("DNS Lookup:    %d µs\n", result.DNSLookup.Microseconds())
		fmt.Printf("TCP Connection:%d µs\n", result.TCPConnection.Microseconds())
		fmt.Printf("TLS Handshake: %d µs\n", result.TLSHandshake.Microseconds())
		fmt.Printf("Server Processing: %d ms\n", result.ServerProcessing.Milliseconds())

		fmt.Printf("Name Lookup:    %d µs\n", result.NameLookup.Microseconds())
		fmt.Printf("Connect:        %d µs\n", result.Connect.Microseconds())
		fmt.Printf("Pre Transfer:   %d µs\n", result.Pretransfer.Microseconds())
		fmt.Printf("Start Transfer: %d µs\n", result.StartTransfer.Microseconds())
		fmt.Println("--------------------")

		time.Sleep(SleepTime)
	}
}
