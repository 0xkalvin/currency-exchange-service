package main

import (
	"exchange-api/utils/logger"

	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var log = logger.NewLogger()

func setupGracefulShutdown(server *http.Server) {
	shutdownSignalListener := make(chan os.Signal, 1)
	signal.Notify(shutdownSignalListener, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-shutdownSignalListener

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info(map[string]interface{}{
		"message": "Server is gracefully shutting down...",
	})

	err := server.Shutdown(ctx)

	if err != nil {
		log.Error(map[string]interface{}{
			"message":       "Failed to shutdown server gracefully.",
			"error_message": err.Error(),
		})
	}

	log.Info(map[string]interface{}{
		"message": "Exiting process...",
	})

}
