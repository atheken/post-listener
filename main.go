package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	done := make(chan error)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "9700"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			done <- errors.New("invalid method used, POST is required")
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			done <- err
			return
		}
		defer r.Body.Close()
		fmt.Println(string(body))
		close(done)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%s", port)}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
			os.Exit(1)
		}
	}()

	err := <-done
	server.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	} else {
		os.Exit(0)
	}
}
