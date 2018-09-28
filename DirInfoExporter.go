package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
)
var(
	ceph_dir_rfiles = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})

	ceph_dir_rbytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)

	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)
)

type(
	CephDirMonitor struct {
		Zone         string
		scrapeChan   chan scrapeRequest
		CephDirRfilesDesc *prometheus.Desc
		CephDirRbytesDesc *prometheus.Desc
	}
	scrapeRequest struct {
		results chan<- prometheus.Metric
		done    chan struct{}
	}
)
func(c *CephDirMonitor)GetDirInfo2map(paths []string)(
	cephDirInfoRfiles map[string] int ,cephDirInfoRbytes map[string]int){
		for _,path := range paths{
			rfiles,rbytes,err := GetDirInfo(path)
			if err!=nil{
				log.Error("Do GetDirInfo failed,path:",path,err.Error())
			}
			cephDirInfoRfiles[path],cephDirInfoRbytes[path] = rfiles,rbytes
		}
		return
}
func (c *CephDirMonitor) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.CephDirRfilesDesc
	ch <- c.CephDirRbytesDesc
}
func (c *CephDirMonitor) Scrape(ch chan<- prometheus.Metric) {
	paths := CounTarDir(strings.Split(arg_paths, ","))
	ceph_dir_rfiles,ceph_dir_rbytes := c.GetDirInfo2map(paths)
	for path, rfiles := range ceph_dir_rfiles {
		ch <- prometheus.MustNewConstMetric(
			c.CephDirRfilesDesc,
			prometheus.CounterValue,
			float64(rfiles),
			path,
		)
	}
	for path, rbytes := range ceph_dir_rbytes {
		ch <- prometheus.MustNewConstMetric(
			c.CephDirRbytesDesc,
			prometheus.GaugeValue,
			float64(rbytes),
			path,
		)
	}
}
func (c *CephDirMonitor) Collect(ch chan<- prometheus.Metric) {
	req := scrapeRequest{results: ch, done: make(chan struct{})}
	c.scrapeChan <- req
	<-req.done
}

func (c *CephDirMonitor) start() {
	for req := range c.scrapeChan {
		ch := req.results
		c.Scrape(ch)
		req.done <- struct{}{}
	}
}
func NewCephDirMonitor(zone string) *CephDirMonitor {
	m:= &CephDirMonitor{
		Zone: zone,
		CephDirRfilesDesc: prometheus.NewDesc(
			"ceph_dir_rfiles",
			"Number of File in the Path.",
			[]string{"path"},
			prometheus.Labels{"cluster": zone},
		),
		CephDirRbytesDesc: prometheus.NewDesc(
			"ceph_dir_rbytes",
			"Bytes of Path.",
			[]string{"path"},
			prometheus.Labels{"cluster": zone},
		),
	}
	go m.start()
	return m
}



func DirInfoExporter() {
	pc := NewCephDirMonitor("ht-01")
	prometheus.MustRegister(pc)

	//http.Handle("/metrics", prometheus.Handler())

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
