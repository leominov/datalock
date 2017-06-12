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
	"github.com/leominov/datalock/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	metrics.InitMetrics()
}

func main() {
	log.Println("Starting datalock...")

	config := server.NewConfig()
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	if err := handlers.ParseTemplates(config); err != nil {
		log.Fatal(err)
	}

	log.Printf("HTTP service listening on %s", config.ListenAddr)

	s := server.New(config)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", handlers.IndexHandler(s))
	mux.Handle("/me", handlers.MeHandler(s))
	mux.Handle("/all_seasons", handlers.AllSeasonsHandler(s))
	mux.Handle("/js/", handlers.ProxyHandler(s))
	mux.Handle("/tpl/asset/js/", handlers.JavaScriptHandler(s))
	mux.Handle("/styleP.php", handlers.StyleHandler(s))
	mux.Handle("/player.php", handlers.ProxyHandler(s))
	mux.Handle("/playls2/", handlers.PlaylistHandler(s))
	mux.Handle("/autocomplete.php", handlers.ProxyHandler(s))

	fs := http.FileServer(http.Dir(s.Config.PublicDir))
	mux.Handle("/public/", http.StripPrefix("/public", fs))

	mux.Handle(s.Config.HealthzPath, handlers.HealthzHandler())
	mux.Handle(s.Config.MetricsPath, promhttp.Handler())

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
		case signal := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", signal))
			httpServer.BlockingClose()
			if err := s.Stop(); err != nil {
				log.Fatal(err)
			}
			log.Println("Bye")
			os.Exit(0)
		}
	}
}
