package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/tcnksm/go-httpstat"
)

const (
	URL           = "http://35.243.94.80/" // リクエスト先URL
	RequestRate   = 100                    // 秒間リクエスト数
	ClientTimeout = 5 * time.Second        // クライアントのタイムアウト
)

func main() {
	client := &http.Client{
		Timeout: ClientTimeout,
		Transport: &http.Transport{
			//MaxIdleConns:    10,
			//IdleConnTimeout: 30 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   1 * time.Second,
				KeepAlive: 15 * time.Second,
			}).DialContext,
		},
	}

	ticker := time.NewTicker(time.Second / time.Duration(RequestRate))
	defer ticker.Stop()

	var wg sync.WaitGroup

	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				makeRequest(client)
			}()
		}
	}
}

func makeRequest(client *http.Client) {
	var result httpstat.Result

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Printf("Request creation failed: %v\n", err)
		return
	}

	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return
	}

	result.End(time.Now())

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-OK HTTP status: %d\n", resp.StatusCode)
		fmt.Printf("DNS Lookup:    %d µs\n", result.DNSLookup.Microseconds())
		fmt.Printf("TCP Connection:%d µs\n", result.TCPConnection.Microseconds())
		fmt.Printf("TLS Handshake: %d µs\n", result.TLSHandshake.Microseconds())
		fmt.Printf("Server Processing: %d ms\n", result.ServerProcessing.Milliseconds())

		fmt.Printf("Name Lookup:    %d µs\n", result.NameLookup.Microseconds())
		fmt.Printf("Connect:        %d µs\n", result.Connect.Microseconds())
		fmt.Printf("Pre Transfer:   %d µs\n", result.Pretransfer.Microseconds())
		fmt.Printf("Start Transfer: %d µs\n", result.StartTransfer.Microseconds())
	}

	// エラーの詳細を表示（成功時には何も出力しない）
	fmt.Println("--------------------")
}
