package main

import (
	"log"
	"path"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

type gsuiteCollector struct {
	logDB                 *sqlx.DB
	confDB                *sqlx.DB
	statusMetric          *prometheus.Desc
	errorCodeMetric       *prometheus.Desc
	startTimeMetric       *prometheus.Desc
	endTimeMetric         *prometheus.Desc
	transferedBytesMetric *prometheus.Desc

	driveSuccessMetric    *prometheus.Desc
	driveWarningMetric    *prometheus.Desc
	driveErrorMetric      *prometheus.Desc
	driveTransferedMetric *prometheus.Desc

	teamdriveSuccessMetric    *prometheus.Desc
	teamdriveWarningMetric    *prometheus.Desc
	teamdriveErrorMetric      *prometheus.Desc
	teamdriveTransferedMetric *prometheus.Desc

	mailSuccessMetric    *prometheus.Desc
	mailWarningMetric    *prometheus.Desc
	mailErrorMetric      *prometheus.Desc
	mailTransferedMetric *prometheus.Desc

	contactSuccessMetric    *prometheus.Desc
	contactWarningMetric    *prometheus.Desc
	contactErrorMetric      *prometheus.Desc
	contactTransferedMetric *prometheus.Desc

	calendarSuccessMetric    *prometheus.Desc
	calendarWarningMetric    *prometheus.Desc
	calendarErrorMetric      *prometheus.Desc
	calendarTransferedMetric *prometheus.Desc
}

type taskName struct {
	ID   int    `db:"task_id"`
	Name string `db:"task_name"`
}

type jobResult struct {
	TaskID         int `db:"task_id"`
	Status         int `db:"execution_status"`
	ErrorCode      int `db:"error_code"`
	TimeStart      int `db:"start_run_time"`
	TimeEnd        int `db:"end_run_time"`
	TransferedSize int `db:"transfered_size"`

	DriveSuccess    int `db:"drive_success_count"`
	DriveWarning    int `db:"drive_warning_count"`
	DriveError      int `db:"drive_error_count"`
	DriveTransfered int `db:"drive_transfered_size"`

	TeamdriveSuccess    int `db:"teamdrive_success_count"`
	TeamdriveWarning    int `db:"teamdrive_warning_count"`
	TeamdriveError      int `db:"teamdrive_error_count"`
	TeamdriveTransfered int `db:"teamdrive_transfered_size"`

	MailSuccess    int `db:"mail_success_count"`
	MailWarning    int `db:"mail_warning_count"`
	MailError      int `db:"mail_error_count"`
	MailTransfered int `db:"mail_transfered_size"`

	ContactSuccess    int `db:"contact_success_count"`
	ContactWarning    int `db:"contact_warning_count"`
	ContactError      int `db:"contact_error_count"`
	ContactTransfered int `db:"contact_transfered_size"`

	CalendarSuccess    int `db:"calendar_success_count"`
	CalendarWarning    int `db:"calendar_warning_count"`
	CalendarError      int `db:"calendar_error_count"`
	CalendarTransfered int `db:"calendar_transfered_size"`
}

func newGSuiteCollector(dataDir string) (*gsuiteCollector, error) {
	logDB, err := sqlx.Connect("sqlite3", path.Join(dataDir, "@ActiveBackup-GSuite/db/log.sqlite"))
	if err != nil {
		return nil, err
	}
	confDB, err := sqlx.Connect("sqlite3", path.Join(dataDir, "@ActiveBackup-GSuite/db/config.sqlite"))
	if err != nil {
		return nil, err
	}

	return &gsuiteCollector{
		logDB:  logDB,
		confDB: confDB,
		statusMetric: prometheus.NewDesc("ab_gsuite_device_result_status",
			"Status of the latest task",
			[]string{"task_name"},
			nil,
		),
		errorCodeMetric: prometheus.NewDesc("ab_gsuite_error_code",
			"Error code of the latest task",
			[]string{"task_name"},
			nil,
		),
		startTimeMetric: prometheus.NewDesc("ab_gsuite_time_start",
			"Start time for the latest task",
			[]string{"task_name"},
			nil,
		),
		endTimeMetric: prometheus.NewDesc("ab_gsuite_time_end",
			"End time for the latest task",
			[]string{"task_name"},
			nil,
		),
		transferedBytesMetric: prometheus.NewDesc("ab_gsuite_transfered_bytes",
			"Total amount of transfered bytes for the complete task",
			[]string{"device_name"},
			nil,
		),
		// Drive
		driveSuccessMetric: prometheus.NewDesc("ab_gsuite_drive_success",
			"Amount of succesfuly backed up drives",
			[]string{"task_name"},
			nil,
		),
		driveWarningMetric: prometheus.NewDesc("ab_gsuite_drive_warning",
			"Amount of backed up drives with warnings",
			[]string{"task_name"},
			nil,
		),
		driveErrorMetric: prometheus.NewDesc("ab_gsuite_drive_error",
			"Amount of backed up drives with errors",
			[]string{"task_name"},
			nil,
		),
		driveTransferedMetric: prometheus.NewDesc("ab_gsuite_drive_transfered_bytes",
			"Amount of transfered bytes for the drives",
			[]string{"task_name"},
			nil,
		),
		// Teamdrive
		teamdriveSuccessMetric: prometheus.NewDesc("ab_gsuite_teamdrive_success",
			"Amount of succesfuly backed up teamdrives",
			[]string{"task_name"},
			nil,
		),
		teamdriveWarningMetric: prometheus.NewDesc("ab_gsuite_teamdrive_warning",
			"Amount of backed up teamdrives with warnings",
			[]string{"task_name"},
			nil,
		),
		teamdriveErrorMetric: prometheus.NewDesc("ab_gsuite_teamdrive_error",
			"Amount of backed up teamdrives with errors",
			[]string{"task_name"},
			nil,
		),
		teamdriveTransferedMetric: prometheus.NewDesc("ab_gsuite_teamdrive_transfered_bytes",
			"Amount of transfered bytes for the teamdrives",
			[]string{"task_name"},
			nil,
		),
		// Mail
		mailSuccessMetric: prometheus.NewDesc("ab_gsuite_mail_success",
			"Amount of succesfuly backed up mail accounts",
			[]string{"task_name"},
			nil,
		),
		mailWarningMetric: prometheus.NewDesc("ab_gsuite_mail_warning",
			"Amount of backed up mail accounts with warnings",
			[]string{"task_name"},
			nil,
		),
		mailErrorMetric: prometheus.NewDesc("ab_gsuite_mail_error",
			"Amount of backed up mail account with errors",
			[]string{"task_name"},
			nil,
		),
		mailTransferedMetric: prometheus.NewDesc("ab_gsuite_mail_transfered_bytes",
			"Amount of transfered bytes for all mails",
			[]string{"task_name"},
			nil,
		),
		// Contact
		contactSuccessMetric: prometheus.NewDesc("ab_gsuite_contact_success",
			"Amount of succesfuly backed up contacts",
			[]string{"task_name"},
			nil,
		),
		contactWarningMetric: prometheus.NewDesc("ab_gsuite_contact_warning",
			"Amount of backed up contacts with warnings",
			[]string{"task_name"},
			nil,
		),
		contactErrorMetric: prometheus.NewDesc("ab_gsuite_contact_error",
			"Amount of backed up contacts with errors",
			[]string{"task_name"},
			nil,
		),
		contactTransferedMetric: prometheus.NewDesc("ab_gsuite_contact_transfered_bytes",
			"Amount of transfered bytes for all contacts",
			[]string{"task_name"},
			nil,
		),
		// Calendar
		calendarSuccessMetric: prometheus.NewDesc("ab_gsuite_calendar_success",
			"Amount of succesfuly backed up calendars",
			[]string{"task_name"},
			nil,
		),
		calendarWarningMetric: prometheus.NewDesc("ab_gsuite_calendar_warning",
			"Amount of backed up calendars with warnings",
			[]string{"task_name"},
			nil,
		),
		calendarErrorMetric: prometheus.NewDesc("ab_gsuite_calendar_error",
			"Amount of backed up calendars with errors",
			[]string{"task_name"},
			nil,
		),
		calendarTransferedMetric: prometheus.NewDesc("ab_gsuite_calendar_transfered_bytes",
			"Amount of transfered bytes for all calendars",
			[]string{"task_name"},
			nil,
		),
	}, nil
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *gsuiteCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.statusMetric
	ch <- collector.errorCodeMetric
	ch <- collector.startTimeMetric
	ch <- collector.endTimeMetric
	ch <- collector.transferedBytesMetric

	ch <- collector.driveSuccessMetric
	ch <- collector.driveWarningMetric
	ch <- collector.driveErrorMetric
	ch <- collector.driveTransferedMetric

	ch <- collector.teamdriveSuccessMetric
	ch <- collector.teamdriveWarningMetric
	ch <- collector.teamdriveErrorMetric
	ch <- collector.teamdriveTransferedMetric

	ch <- collector.mailSuccessMetric
	ch <- collector.mailWarningMetric
	ch <- collector.mailErrorMetric
	ch <- collector.mailTransferedMetric

	ch <- collector.contactSuccessMetric
	ch <- collector.contactWarningMetric
	ch <- collector.contactErrorMetric
	ch <- collector.contactTransferedMetric

	ch <- collector.calendarSuccessMetric
	ch <- collector.calendarWarningMetric
	ch <- collector.calendarErrorMetric
	ch <- collector.calendarTransferedMetric
}

func taskNameLookup(lookup []taskName, id int) string {
	for _, entry := range lookup {
		if entry.ID == id {
			return entry.Name
		}
	}

	return string(id)

}

//Collect implements required collect function for all promehteus collectors
func (collector *gsuiteCollector) Collect(ch chan<- prometheus.Metric) {
	taskNames := []taskName{}
	err := collector.confDB.Select(&taskNames, "SELECT task_id, task_name FROM task_info_table;")
	if err != nil {
		log.Fatal(err)
	}

	results := []jobResult{}
	err = collector.logDB.Select(&results, ` SELECT task_id, execution_status, transfered_size, start_run_time, end_run_time, error_code,
		drive_success_count,drive_warning_count,drive_error_count,drive_transfered_size,
		teamdrive_success_count,teamdrive_warning_count,teamdrive_error_count,teamdrive_transfered_size,
		mail_success_count,mail_warning_count,mail_error_count,mail_transfered_size,
		contact_success_count,contact_warning_count,contact_error_count,contact_transfered_size,
		calendar_success_count,calendar_warning_count,calendar_error_count,calendar_transfered_size 
		FROM job_log_table
		WHERE job_execution_id IN (SELECT max(job_execution_id) FROM job_log_table GROUP BY task_id);`)

	if err != nil {
		log.Fatal(err)
	}

	for _, res := range results {
		labels := []string{taskNameLookup(taskNames, res.TaskID)}

		ch <- prometheus.MustNewConstMetric(collector.statusMetric, prometheus.GaugeValue, float64(res.Status), labels...)
		ch <- prometheus.MustNewConstMetric(collector.errorCodeMetric, prometheus.GaugeValue, float64(res.ErrorCode), labels...)
		ch <- prometheus.MustNewConstMetric(collector.startTimeMetric, prometheus.GaugeValue, float64(res.TimeStart), labels...)
		ch <- prometheus.MustNewConstMetric(collector.endTimeMetric, prometheus.GaugeValue, float64(res.TimeEnd), labels...)
		ch <- prometheus.MustNewConstMetric(collector.transferedBytesMetric, prometheus.GaugeValue, float64(res.TransferedSize), labels...)

		ch <- prometheus.MustNewConstMetric(collector.driveSuccessMetric, prometheus.GaugeValue, float64(res.DriveSuccess), labels...)
		ch <- prometheus.MustNewConstMetric(collector.driveErrorMetric, prometheus.GaugeValue, float64(res.DriveError), labels...)
		ch <- prometheus.MustNewConstMetric(collector.driveWarningMetric, prometheus.GaugeValue, float64(res.DriveWarning), labels...)
		ch <- prometheus.MustNewConstMetric(collector.driveTransferedMetric, prometheus.GaugeValue, float64(res.DriveTransfered), labels...)

		ch <- prometheus.MustNewConstMetric(collector.teamdriveSuccessMetric, prometheus.GaugeValue, float64(res.TeamdriveSuccess), labels...)
		ch <- prometheus.MustNewConstMetric(collector.teamdriveErrorMetric, prometheus.GaugeValue, float64(res.TeamdriveError), labels...)
		ch <- prometheus.MustNewConstMetric(collector.teamdriveWarningMetric, prometheus.GaugeValue, float64(res.TeamdriveWarning), labels...)
		ch <- prometheus.MustNewConstMetric(collector.teamdriveTransferedMetric, prometheus.GaugeValue, float64(res.TeamdriveTransfered), labels...)

		ch <- prometheus.MustNewConstMetric(collector.mailSuccessMetric, prometheus.GaugeValue, float64(res.MailSuccess), labels...)
		ch <- prometheus.MustNewConstMetric(collector.mailErrorMetric, prometheus.GaugeValue, float64(res.MailError), labels...)
		ch <- prometheus.MustNewConstMetric(collector.mailWarningMetric, prometheus.GaugeValue, float64(res.MailWarning), labels...)
		ch <- prometheus.MustNewConstMetric(collector.mailTransferedMetric, prometheus.GaugeValue, float64(res.MailTransfered), labels...)

		ch <- prometheus.MustNewConstMetric(collector.contactSuccessMetric, prometheus.GaugeValue, float64(res.ContactSuccess), labels...)
		ch <- prometheus.MustNewConstMetric(collector.contactErrorMetric, prometheus.GaugeValue, float64(res.ContactWarning), labels...)
		ch <- prometheus.MustNewConstMetric(collector.contactWarningMetric, prometheus.GaugeValue, float64(res.ContactWarning), labels...)
		ch <- prometheus.MustNewConstMetric(collector.contactTransferedMetric, prometheus.GaugeValue, float64(res.ContactTransfered), labels...)

		ch <- prometheus.MustNewConstMetric(collector.calendarSuccessMetric, prometheus.GaugeValue, float64(res.CalendarSuccess), labels...)
		ch <- prometheus.MustNewConstMetric(collector.calendarErrorMetric, prometheus.GaugeValue, float64(res.CalendarError), labels...)
		ch <- prometheus.MustNewConstMetric(collector.calendarWarningMetric, prometheus.GaugeValue, float64(res.CalendarWarning), labels...)
		ch <- prometheus.MustNewConstMetric(collector.calendarTransferedMetric, prometheus.GaugeValue, float64(res.CalendarTransfered), labels...)

	}

}
