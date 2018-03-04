package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
)

const (
	namespace = "datalock"
)

var (
	once sync.Once

	HttpRequestsErrorCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "http_requests_error_count"),
		Help: "How many errors while requesting remote data.",
	})

	HttpRequestsTotalCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "http_requests_total_count"),
		Help: "How many requesting remote data.",
	})

	TemplateExecuteErrorCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "template_execute_error_count"),
		Help: "How many errors while executing template.",
	})

	SeasonIDErrorCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "season_id_error_count"),
		Help: "How many errors while getting season id.",
	})

	SerialIDErrorCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "serial_id_error_count"),
		Help: "How many errors while getting serial id.",
	})

	SeasonTitleErrorCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "season_title_error_count"),
		Help: "How many errors while getting season title.",
	})

	SeasonKeywordsErrorCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "season_keywords_error_count"),
		Help: "How many errors while getting season keywords.",
	})

	SeasonDescriptionErrorCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "season_description_error_count"),
		Help: "How many errors while getting season description.",
	})
)

func InitMetrics() {
	once.Do(func() {
		prometheus.MustRegister(version.NewCollector(namespace))
		prometheus.MustRegister(HttpRequestsErrorCount)
		prometheus.MustRegister(HttpRequestsTotalCount)
		prometheus.MustRegister(TemplateExecuteErrorCount)
		prometheus.MustRegister(SeasonIDErrorCount)
		prometheus.MustRegister(SerialIDErrorCount)
		prometheus.MustRegister(SeasonTitleErrorCount)
		prometheus.MustRegister(SeasonKeywordsErrorCount)
		prometheus.MustRegister(SeasonDescriptionErrorCount)
	})
}
