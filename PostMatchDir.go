package main

import (
	"DirMonitor/model"
	"bytes"
	"encoding/json"
	"flag"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"unsafe"
)

func Post2Platfm(path model.DatadataPost) {
	url := "http://www.baidu.com"
	data_post, err := json.Marshal(path)
	if err != nil {
		log.Error("Do json.Marshal failed: \n", err.Error())
	}
	reader := bytes.NewReader(data_post)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Error("Do http.NewRequest failed: \n", err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Error("Do client.Do failed: \n", err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Do ioutil.ReadAll failed: \n", err.Error())
		return
	}
	str := (*string)(unsafe.Pointer(&respBytes))
	log.Info(str)
}
func PostMatchDir() {
	flag.Parse()
	if h {
		flag.Usage()
	}
	var ans model.DatadataPost
	paths := CounTarDir(strings.Split(arg_paths, ","))
	for _, tar_dir := range paths {
		flag, info := IsMatch(tar_dir, num_file_limit, num_byte_limit)
		if flag == 1 {
			item := model.DataPath{TarPath: tar_dir, NumFile: info["files"], SumBytes: info["bytes"], AvgBytes: info["avg"]}
			ans.DataFileLimit = append(ans.DataFileLimit, item)
		} else if flag == 2 {
			item := model.DataPath{TarPath: tar_dir, NumFile: info["files"], SumBytes: info["bytes"], AvgBytes: info["avg"]}
			ans.DataFileLimit = append(ans.DataFileLimit, item)
		}
	}
	Post2Platfm(ans)
}
