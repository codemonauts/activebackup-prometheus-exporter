# activebackup-prometheus-exporter

This is intendet to run on a Synology NAS and then export metrics from various **ActiveBackup** tools which then can
be scraped by Prometheus.

Currently supported Backup Tools are:
  * ActiveBackup for Business
  * ActiveBackup for GSuite

# Usage
`-dir` is the path to the directory where the tools have their data folders.

```
./activebackup-prometheus-exporter -dir /volume1/
```
