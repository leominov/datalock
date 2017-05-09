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
	"github.com/leominov/datalock/metrics"
	"github.com/leominov/datalock/seasonvar"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	metrics.InitMetrics()
}

func main() {
	log.Println("Starting datalock...")

	if err := handlers.ParseTemplates(); err != nil {
		log.Fatal(err)
	}

	config := seasonvar.NewConfig()
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	log.Printf("HTTP service listening on %s", config.ListenAddr)

	s := seasonvar.New(config)

	mux := http.NewServeMux()
	mux.Handle("/", handlers.IndexHandler(s))
	mux.Handle("/me", handlers.MeHandler(s))
	mux.Handle(s.Config.HealthzPath, handlers.HealthzHandler())
	mux.Handle(s.Config.MetricsPath, promhttp.Handler())
	mux.Handle("/js/", handlers.ProxyHandler(s))
	mux.Handle("/tpl/asset/js/", handlers.JavaScriptHandler(s))
	mux.Handle("/styleP.php", handlers.ProxyHandler(s))
	mux.Handle("/player.php", handlers.ProxyHandler(s))
	mux.Handle("/playls2/", handlers.PlaylistHandler(s))

	httpServer := manners.NewServer()
	httpServer.Addr = config.ListenAddr
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
			httpServer.BlockingClose()
			os.Exit(0)
		}
	}
}
