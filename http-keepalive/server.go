package main

import (
	"context"
	"fmt"
	"golang.org/x/xerrors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		IdleTimeout: 5 * time.Second,
	}

	fmt.Println("Starting server on :8080")
	//if err := server.ListenAndServe(); err != nil {
	//	fmt.Printf("Server error: %v\n", err)
	//}

	startServerErr := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			startServerErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	fmt.Println("Server is ready to handle requests at :8080")

	select {
	case <-quit:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic(xerrors.Errorf("failed to graceful shutdown: %w", err))
	}
	fmt.Println("successfully graceful shutdown server")
}
