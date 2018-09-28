package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
)

var (
	h                              bool
	arg_paths, mode                    string
	num_file_limit, num_byte_limit int
)

func usage() {
	fmt.Fprintf(os.Stderr, `FindMatchDir Version: FindMatchDir/0.1
Usage: DirMonitor [-h] [-m mode] [-p paths] [-NF num_file_limit] [-NB num_byte_limit] 
Options:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `Example:
DirMonitor -m save -p /mnt/cephfs/algor-api/user/*/* /mnt/cephfs/algor-api/dataset/user/*/* /mnt/cephfs/algor-api/dataset/public/* -NF 10000 -NB 1048576
`)

}

func init() {
	flag.BoolVar(&h, "h", false, "show help")
	flag.StringVar(&mode, "m", "", "运行模式，save:将统计结果把存在本地，post:将统计结果post到平台,exporter:用以Prometheus抓取目录信息。")
	flag.StringVar(&arg_paths, "p", "", "列举操作的目录，可使用通配符表示，多个以逗号隔开。")
	flag.IntVar(&num_file_limit, "NF", 0, "设置的目录下文件数上限。")
	flag.IntVar(&num_byte_limit, "NB", 0, "设置的整个目录大小上限（单位：Byte）。")
	flag.Usage = usage

	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true                    // 显示完整时间
	customFormatter.TimestampFormat = "2006-01-02 15:04:05" // 时间格式
	customFormatter.DisableTimestamp = false                // 禁止显示时间
	customFormatter.DisableColors = false                   // 禁止颜色显示
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}
func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	if mode == "save" {
		SaveMatchDir()
	} else if mode == "post" {
		PostMatchDir()
	} else if mode == "exporter" {
		DirInfoExporter()
	} else {
		log.Error("mode参数输入有误。\n")
		flag.Usage()
	}
}
