package main

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"testing"
)

func init() {
	//prometheus.
}

func Test_DirInfoExporter(t *testing.T) {
	//DirInfoExporter()
	log.Info("staring...")
	pc := NewCephDirMonitor("ht-01")
	log.Info("NewCephDirMonitor success.")
	prometheus.MustRegister(pc)
	log.Info("MustRegister success.")

	http.Handle("/metrics", prometheus.Handler())
	log.Info("Handle success.")

	http.HandleFunc("/q", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Named Process Exporter</title></head>
			<body>
			<h1>Named Process Exporter</h1>
			<p><a href="` + "/metrics" + `">Metrics</a></p>
			</body>
			</html>`))
	})
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Unable to setup HTTP server: %v", err)
	}
}


