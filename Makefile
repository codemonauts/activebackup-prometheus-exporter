BINARY=activebackup-prometheus-exporter
build: clean
	go build -o $(BINARY) main.go ab_gsuite.go ab_business.go

clean:
	rm -f $(BINARY)
