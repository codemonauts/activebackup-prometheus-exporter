package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var dataDir string

func init() {
	flag.StringVar(&dataDir, "dir", "/volume1/@ActiveBackup", "Path to the ABB folder")
}

func main() {
	flag.Parse()

	if _, err := os.Stat(path.Join(dataDir, "activity.db")); os.IsNotExist(err) {
		log.Fatal("No activity.db in the given directory. Check your path")
	}

	if _, err := os.Stat(path.Join(dataDir, "config.db")); os.IsNotExist(err) {
		log.Fatal("No config.db in the given directory. Check your path")
	}

	abbCollector := newABBCollector(dataDir)
	prometheus.MustRegister(abbCollector)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8000", nil)
}
