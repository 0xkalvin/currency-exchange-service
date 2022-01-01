package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupGracefulShutdown(server *http.Server) {
	shutdownSignalListener := make(chan os.Signal, 1)
	signal.Notify(shutdownSignalListener, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-shutdownSignalListener

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Server is gracefully shutting down...")

	err := server.Shutdown(ctx)

	if err != nil {
		fmt.Println("Failed to shutdown server gracefully")
	}

	fmt.Println("Exiting process...")
}
