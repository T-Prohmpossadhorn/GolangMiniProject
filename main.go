package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/T-Prohmpossadhorn/GolangMiniProject/internal"
	"github.com/T-Prohmpossadhorn/GolangMiniProject/internal/config"
)

func main() {
	log.Println("rf-financialtransaction", "go", runtime.Version())

	svc, err := internal.New(context.Background(), config.Options{
		ListenAddress:     os.Getenv("LISTEN_ADDRESS"),
		ListenAddressHTTP: os.Getenv("LISTEN_ADDRESS_HTTP"),
		FilePath:          os.Getenv("FILE_PATH"),
	})

	if err != nil {
		log.Fatal("internal.New(): %w", err)
	}

	shutdownOnSignal(svc)
}

func waitForShutdownSignal() string {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until signaled
	sig := <-c

	return sig.String()
}

func shutdownOnSignal(svc *internal.Service) {
	signalName := waitForShutdownSignal()
	log.Println("Received signal, starting shutdown", "signal", signalName)

	if svc.Shutdown() {
		log.Println("Shutdown complete")
	} else {
		log.Println("Shutdown timed out")
	}
}
