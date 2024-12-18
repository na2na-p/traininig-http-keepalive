package main

import (
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request: %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func main() {
	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr:        ":8080",
		ReadTimeout: 1 * time.Second,
		//WriteTimeout: 3 * time.Second,
		IdleTimeout: 3 * time.Second,
	}

	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
