package main

import (
	"log"
	"path"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

type abbCollector struct {
	activityDB_Path       string
	statusMetric          *prometheus.Desc
	startTimeMetric       *prometheus.Desc
	endTimeMetric         *prometheus.Desc
	transferedBytesMetric *prometheus.Desc
}

type deviceResult struct {
	Status          int
	DeviceName      string `db:"device_name"`
	TimeStart       int    `db:"time_start"`
	TimeEnd         int    `db:"time_end"`
	TransferedBytes int    `db:"transfered_bytes"`
}

func newABBCollector(dataDir string) (*abbCollector, error) {

	// TODO: Keine connection reinreichen sondern nur filename damit man innendrin die connection auf und zu machen kann
	activityDB_Path := path.Join(dataDir, "@ActiveBackup/activity.db")

	return &abbCollector{
		activityDB_Path: activityDB_Path,
		statusMetric: prometheus.NewDesc("ab_business_device_result_status",
			"Status of the latest device backup",
			[]string{"device_name"},
			nil,
		),
		startTimeMetric: prometheus.NewDesc("ab_business_device_result_time_start",
			"Start time for the latest device backup",
			[]string{"device_name"},
			nil,
		),
		endTimeMetric: prometheus.NewDesc("ab_business_device_result_time_end",
			"End time for the latest device backup",
			[]string{"device_name"},
			nil,
		),
		transferedBytesMetric: prometheus.NewDesc("ab_business_device_result_transfered_bytes",
			"Amount of transfered bytes for the latest device backup",
			[]string{"device_name"},
			nil,
		),
	}, nil
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *abbCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.statusMetric
	ch <- collector.startTimeMetric
	ch <- collector.endTimeMetric
	ch <- collector.transferedBytesMetric
}

// Collect implements required collect function for all promehteus collectors
func (collector *abbCollector) Collect(ch chan<- prometheus.Metric) {
	results := []deviceResult{}
	activityDB, err := sqlx.Connect("sqlite3", collector.activityDB_Path)
	err = activityDB.Select(&results, `SELECT status, device_name, time_start,time_end,transfered_bytes FROM device_result_table
			WHERE device_result_id IN (SELECT max(device_result_id) FROM device_result_table WHERE time_end != 0 GROUP BY device_name);`)
	if err != nil {
		log.Fatalln(err)

	}
	for _, res := range results {
		labels := []string{res.DeviceName}
		ch <- prometheus.MustNewConstMetric(collector.statusMetric, prometheus.GaugeValue, float64(res.Status), labels...)
		ch <- prometheus.MustNewConstMetric(collector.startTimeMetric, prometheus.GaugeValue, float64(res.TimeStart), labels...)
		ch <- prometheus.MustNewConstMetric(collector.endTimeMetric, prometheus.GaugeValue, float64(res.TimeEnd), labels...)
		ch <- prometheus.MustNewConstMetric(collector.transferedBytesMetric, prometheus.GaugeValue, float64(res.TransferedBytes), labels...)

	}

}
