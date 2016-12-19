package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/braintree/manners"
	"github.com/leominov/datalock/handlers"
	"github.com/leominov/datalock/seasonvar"
)

const (
	DefaultHTTPAddr = "127.0.0.1:7000"
)

func main() {
	log.Println("Starting datalock...")

	httpAddr := os.Getenv("DATALOCK_LISTEN_ADDR")
	if httpAddr == "" {
		httpAddr = DefaultHTTPAddr
	}

	log.Printf("HTTP service listening on %s", httpAddr)

	seasonvar := seasonvar.New()

	mux := http.NewServeMux()
	mux.Handle("/", handlers.IndexHandle(seasonvar))
	mux.Handle("/player", handlers.PlayerHandle(seasonvar))

	httpServer := manners.NewServer()
	httpServer.Addr = httpAddr
	httpServer.Handler = handlers.LoggingHandler(mux)

	errChan := make(chan error, 10)

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			os.Exit(0)
		}
	}
}
