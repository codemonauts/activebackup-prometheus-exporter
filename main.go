package main

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


func main() {

	abbCollector := newABBCollector()
	prometheus.MustRegister(abbCollector)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8000", nil)
}
