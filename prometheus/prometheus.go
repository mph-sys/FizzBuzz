package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	counterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "fizzbuzz_processed_ops_total",
		Help: "The total number of processed events by status",
	}, []string{"job", "status"})
)

func Start(prometheusBindAddr string) {
	log.Printf("start prometheus server on port %s", prometheusBindAddr)

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		// default collectors
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),

		counterVec,
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(prometheusBindAddr, nil)
}

// Increment total counter of job requested
func IncRequest(job string) {
	counterVec.WithLabelValues(job, "request").Inc()
}

// Increment total counter of job by status (e.g "success", "error"...)
func IncStats(job, status string) {
	counterVec.WithLabelValues(job, status).Inc()
}
