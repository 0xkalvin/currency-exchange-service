package main

import (
	datasources "exchange-api/data-sources"
	"fmt"
	"net/http"
	"time"
)

func main() {
	httpHandler := http.NewServeMux()
	sqs := datasources.BuildSQSClient()

	setupHandlers(httpHandler, sqs)

	server := &http.Server{
		Addr:         ":3000",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      httpHandler,
	}

	fmt.Println("Server is up and kicking on port 3000")

	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Failed to start server")
		}

	}()

	setupGracefulShutdown(server)
}
