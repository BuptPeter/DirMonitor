package main

import (
	"DirMonitor/model"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
)

var (
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

)

type (
	CephDirMonitor struct {
		Zone              string
		scrapeChan        chan scrapeRequest
		CephDirRfilesDesc *prometheus.Desc
		CephDirRbytesDesc *prometheus.Desc
	}
	scrapeRequest struct {
		results chan<- prometheus.Metric
		done    chan struct{}
	}
)
func (c *CephDirMonitor) GetDirInfo2map(paths []string) (cephDirInfoRfiles map[string]int, cephDirInfoRbytes map[string]int) {
	cephDirInfoRfiles = make(map[string]int)
	cephDirInfoRbytes = make(map[string]int)
	for _, path := range paths {
		rfiles, rbytes, err := GetDirInfo(path)
		if err != nil {
			log.Error("Do GetDirInfo failed,path:", path, err.Error())
		}
		cephDirInfoRfiles[path], cephDirInfoRbytes[path] = rfiles, rbytes
	}
	return cephDirInfoRfiles, cephDirInfoRbytes
}
func (c *CephDirMonitor) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.CephDirRfilesDesc
	ch <- c.CephDirRbytesDesc
}
func (c *CephDirMonitor) Scrape(ch chan<- prometheus.Metric) {
	var user_info model.UserInfo
	log.Info("start Scrape...  ")
	paths := CounTarDir(strings.Split(arg_paths, ","))
	log.Info(paths)
	ceph_dir_rfiles, ceph_dir_rbytes := c.GetDirInfo2map(paths)
	for path, rfiles := range ceph_dir_rfiles {
		if _,ok := user_info_cache[path]; ok {
			user_info = user_info_cache[path]
		}else {
			info,err := GetUserInfo(path)
			if err != nil{
				log.Error("GetUserInfo failed:",err.Error())
			}
			user_info = info
			user_info_cache[path] = user_info
		}
		ch <- prometheus.MustNewConstMetric(
			c.CephDirRfilesDesc,
			prometheus.GaugeValue,
			float64(rfiles),
			path,user_info.RealName,user_info.Erp,user_info.Email,

		)
	}
	for path, rbytes := range ceph_dir_rbytes {
		if _,ok := user_info_cache[path]; ok {
			user_info = user_info_cache[path]
		}else {
			info,err := GetUserInfo(path)
			if err != nil{
				log.Error("GetUserInfo failed:",err.Error())
			}
			user_info = info
			user_info_cache[path] = user_info
		}
		ch <- prometheus.MustNewConstMetric(
			c.CephDirRbytesDesc,
			prometheus.GaugeValue,
			float64(rbytes),
			path,user_info.RealName,user_info.Erp,user_info.Email,
		)
	}
}
func (c *CephDirMonitor) Collect(ch chan<- prometheus.Metric) {
	log.Info("Receive request from Web Server...")
	req := scrapeRequest{results: ch, done: make(chan struct{})}
	log.Info("make scrapeRequst successfully...")
	c.scrapeChan <- req
	log.Info("Put request to scrapeChan successfully...")
	<-req.done
	log.Info("This requst handled successfully...")
}
func (c *CephDirMonitor) start() {
	for req := range c.scrapeChan {
		log.Info("Handling req from scrapeChan")
		ch := req.results
		c.Scrape(ch)
		req.done <- struct{}{}
	}
	log.Info("Stoping handle rep from scrapeChan")
}
func NewCephDirMonitor(zone string) *CephDirMonitor {
	m := &CephDirMonitor{
		Zone: zone,
		scrapeChan: make(chan scrapeRequest),
		CephDirRfilesDesc: prometheus.NewDesc(
			"ceph_dir_rfiles",
			"Number of File in the Path.",
			[]string{"path","name","erp","email"},
			prometheus.Labels{"cluster": zone},
		),
		CephDirRbytesDesc: prometheus.NewDesc(
			"ceph_dir_rbytes",
			"Bytes of Path.",
			[]string{"path","name","erp","email"},
			prometheus.Labels{"cluster": zone},
		),
	}
	go m.start()
	log.Info("NewCephDirMonitor:Listening scrape request...")
	return m
}

func DirInfoExporter() {
	log.Info("Working as exporter mode...\nStarting to scrape user_info to cache,plase wait a minute...")
	PutInfo2Cache()
	pc := NewCephDirMonitor(arg_cluster)
	prometheus.MustRegister(pc)
	http.Handle("/metrics", prometheus.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Named Process Exporter</title></head>
			<body>
			<h1>Named Process Exporter</h1>
			<p><a href="` + "/metrics" + `">Metrics</a></p>
			</body>
			</html>`))
	})
	log.Info("Scrape data success,Starting web Service...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Unable to setup HTTP server: %v", err)
	}
}