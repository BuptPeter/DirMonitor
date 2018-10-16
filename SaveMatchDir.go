package main

import (
	"DirMonitor/model"
	"flag"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

func SortMapByValue(m map[string]int) model.DirList {
	p := make(model.DirList, len(m))
	i := 0
	for k, v := range m {
		p[i] = model.Dir{k, v}
	}
	sort.Sort(p)
	return p
}

func SaveDir(dir_list []model.Dir, file_path string) error {
	var (
		ans string
	)
	for _, i := range dir_list {
		ans = ans + i.Key + "\n"
	}
	file_byte := []byte(ans)
	err := ioutil.WriteFile(file_path, file_byte, 0644)
	if err != nil { //执行WriteFile失败
		log.Error("Do WriteFile failed: \n", err.Error())
		return err
	}
	return nil
}

func SaveMatchDir() {
	flag.Parse()
	if h {
		flag.Usage()
	}
	var (
		ans1, ans2 map[string]int
	)
	paths := CounTarDir(strings.Split(arg_paths, ","))
	log.Info("###All Paths Counted:",paths)
	for _, tar_dir := range paths {
		flag, info := IsMatch(tar_dir, num_file_limit, num_byte_limit)
		if flag == 1 {
			ans1[tar_dir+" FileNum:"+strconv.Itoa(info["file"])+", Bytes:"+strconv.Itoa(info["bytes"])+", AvgByte:"+strconv.Itoa(info["avg"])] = info["file"]
		} else if flag == 2 {
			ans2[tar_dir+" FileNum:"+strconv.Itoa(info["file"])+", Bytes:"+strconv.Itoa(info["bytes"])+", AvgByte:"+strconv.Itoa(info["avg"])] = info["file"]
		}
	}
	log.Info("###All Paths Info:",ans1)
	log.Info("###All Paths Info(after sort):",SortMapByValue(ans1))

	err := SaveDir(SortMapByValue(ans1), "MatchDirList1.txt")
	if err != nil {
		log.Error("Do SaveDir failed: \n", err.Error())
	}
	err = SaveDir(SortMapByValue(ans2), "MatchDirList2.txt")
	if err != nil {
		log.Error("Do SaveDir failed: \n", err.Error())
	}
}
