package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	dataDir string
	port    int
)

func init() {
	flag.StringVar(&dataDir, "dir", "/volume1", "Path to the volume containing the ActiveBackup folders")
	flag.IntVar(&port, "port", 9763, "Listening port")
}

func main() {
	flag.Parse()

	count := 0

	if _, err := os.Stat(path.Join(dataDir, "@ActiveBackup/activity.db")); os.IsNotExist(err) {
		log.Println("You don't have 'ActiveBackup for Business' installed")
	} else {
		count++
		log.Println("Found an @ActiveBackup folder. Loading the 'ActiveBackup for Business' module")
		abbCollector, err := newABBCollector(dataDir)
		if err != nil {
			log.Println(err)
		}
		prometheus.MustRegister(abbCollector)
	}

	if _, err := os.Stat(path.Join(dataDir, "@ActiveBackup-GSuite/db/config.sqlite")); os.IsNotExist(err) {
		log.Println("You don't have 'ActiveBackup for GSuite' installed")
	} else {
		log.Println("Found an @ActiveBackup-GSuite folder. Loading the 'ActiveBackup for GSuite' module")
		gsuiteCollector, err := newGSuiteCollector(dataDir)
		if err != nil {
			log.Println(err)
		}
		count++
		prometheus.MustRegister(gsuiteCollector)
	}

	if count == 0 {
		log.Fatal("It appears that you don't have any ActiveBackup software installed. Exiting.")
		os.Exit(1)
	}

	http.Handle("/metrics", promhttp.Handler())

	fmt.Printf("Listening on :%d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
