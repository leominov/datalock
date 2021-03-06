package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/braintree/manners"
	"github.com/leominov/datalock/pkg/metrics"
	"github.com/leominov/datalock/pkg/server"
	"github.com/leominov/datalock/pkg/server/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

var (
	showVersion   = flag.Bool("version", false, "Show version and exit")
	blacklistPath = flag.String("blacklist.path", "", "RKN blacklist path")
)

func init() {
	metrics.InitMetrics()
}

func getVersion() string {
	if len(version.Version) == 0 {
		return "X.X"
	}
	return version.Version
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())

	if *showVersion {
		fmt.Println(version.Print("datalock"))
		return
	}

	log.Printf("Starting datalock %s...", getVersion())

	config := server.NewConfig()
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	if err := server.ParseTemplates(config); err != nil {
		log.Fatal(err)
	}

	log.Printf("HTTP service listening on %s", config.ListenAddr)

	s, err := server.New(config)
	if err != nil {
		log.Fatal(err)
	}

	err = s.LoadBlacklist(*blacklistPath)
	if err != nil {
		log.Fatal(err)
	}

	go s.Run()

	mux := http.NewServeMux()
	mux.Handle("/", handlers.IndexHandler(s))
	mux.Handle("/styleP.php", handlers.StyleHandler(s))
	mux.Handle("/player.php", handlers.PlayerHandler(s))
	mux.Handle("/plStat.php", handlers.NoContentHandler(s))
	mux.Handle("/autocomplete.php", handlers.ProxyHandler(s, false))
	mux.Handle("/tagautocomplete.php", handlers.ProxyHandler(s, false))
	// Examples:
	// /playls2/0/trans/14340/plist.txt
	// /playls2/0/trans/14340/list.xml
	mux.Handle("/playls2/", handlers.PlaylistHandler(s))
	mux.Handle("/crossdomain.xml", handlers.CrossdomainHandler(s))
	mux.Handle("/sitemap.xml", handlers.ProxyHandler(s, true))
	mux.Handle("/rss.php", handlers.ProxyHandler(s, true))
	// Static files
	mux.Handle("/uppod/", handlers.ProxyHandler(s, false))
	mux.Handle("/js/", handlers.ProxyHandler(s, false))
	mux.Handle("/sub/", handlers.ProxyHandler(s, false))
	mux.Handle("/favicon.ico", handlers.ProxyHandler(s, false))
	mux.Handle("/tpl/asset/js/", handlers.JavaScriptHandler(s))
	mux.Handle("/tpl/asset/", handlers.ProxyHandler(s, false))
	mux.Handle("/hls.min.js.map", handlers.NoContentHandler(s))
	// Interface helpers
	mux.Handle("/api/all_seasons", handlers.ApiAllSeasonsHandler(s))
	mux.Handle("/api/all_season_series", handlers.ApiAllSeasonSeriesHandler(s))
	mux.Handle("/api/all_series", handlers.ApiAllSeasonSeriesHandler(s))
	mux.Handle("/api/info_season", handlers.ApiInfoSeasonHandler(s))

	// http://seasonvar.ru/ajax.php?mode=new
	mux.Handle("/api/updated_series", handlers.ApiListSeriesHandler(s, "new"))
	// http://seasonvar.ru/ajax.php?mode=pop
	mux.Handle("/api/popular_series", handlers.ApiListSeriesHandler(s, "pop"))
	// http://seasonvar.ru/ajax.php?mode=newest
	mux.Handle("/api/new_series", handlers.ApiListSeriesHandler(s, "newest"))

	fs := http.FileServer(http.Dir(s.Config.PublicDir))
	mux.Handle("/public/", http.StripPrefix("/public", fs))
	mux.Handle("/robots.txt", fs)

	mux.Handle(s.Config.HealthzPath, handlers.HealthzHandler())
	mux.Handle(s.Config.MetricsPath, promhttp.Handler())

	httpServer := manners.NewServer()
	httpServer.Addr = config.ListenAddr
	httpServer.Handler = handlers.MiddlewareHandler(s, mux)

	errChan := make(chan error, 10)

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case signal := <-signalChan:
			switch signal {
			case syscall.SIGUSR1:
				log.Printf("Captured %v. Templates cache will be flushed on next HTTP-request", signal)
				s.MarkFlushTemplatesCache()
			default:
				log.Printf("Captured %v. Exiting...", signal)
				httpServer.BlockingClose()
				if err := s.Stop(); err != nil {
					log.Fatal(err)
				}
				log.Println("Bye")
				os.Exit(0)
			}
		}
	}
}
