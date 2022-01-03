package main

import (
	datasources "exchange-api/data-sources"
	"exchange-api/utils/logger"
	"os"

	"net/http"
	"time"
)

var PORT = os.Getenv("PORT")
var LOG_LEVEL = os.Getenv("LOG_LEVEL")

func main() {
	httpHandler := http.NewServeMux()
	sqs := datasources.BuildSQSClient()
	log := logger.NewLogger()

	setupHandlers(httpHandler, sqs)

	server := &http.Server{
		Addr:         ":" + PORT,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      httpHandler,
	}

	log.Info(map[string]interface{}{
		"message":          "Server is up and kicking",
		"port":             PORT,
		"server_log_level": LOG_LEVEL,
	})

	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Error(map[string]interface{}{
				"message":       "Failed to start server",
				"error_message": err.Error(),
			})
		}

	}()

	setupGracefulShutdown(server)
}
