package main

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

func IsMatch(path string, files_limit int, byte_limit int) (int, map[string]int) {

	var info map[string]int
	num_file,num_byte,err := GetDirInfo(path)
	if err != nil{
		log.Error("Do GetDirInfo failed: \n", err.Error())
		return -1,nil
	}
	info["files"] = num_file
	if num_file < files_limit { //若文件数下文件数小于files_limit
		return 0, info
	}

	info["bytes"] = num_byte
	avg := num_byte / num_file
	info["avg"] = avg
	if avg < byte_limit { //#若平均每个文件小于byte_limit
		return 1, info
	} else {
		return 2, info
	}

}

func GetDirInfo(path string)( int, int,error){
	cmd_get_files := "getfattr -d -n ceph.dir.rfiles " + path

	var stdErr0, stdOut0 bytes.Buffer
	cmd0 := exec.Command(cmd_get_files)
	cmd0.Stderr = &stdErr0
	cmd0.Stdout = &stdOut0
	err := cmd0.Run()
	if err != nil { //执行getfattr失败,返回错误信息及标准输出
		log.Error("Do cmd failed: \n", err.Error()+stdErr0.String())
		return -1,-1,err
	}
	s := strings.Split(stdOut0.String(), "\"")[1]
	rfiles, err := strconv.Atoi(s)
	if err != nil {
		log.Error("get num_file failed: \n", err.Error())
		return -1,-1,err
	}

	var stdErr1, stdOut1 bytes.Buffer
	cmd_get_bytes := "getfattr -d -n ceph.dir.rbytes " + path
	cmd1 := exec.Command(cmd_get_bytes)
	cmd1.Stderr = &stdErr1
	cmd1.Stdout = &stdOut1
	err = cmd1.Run()
	if err != nil { //执行getfattr失败
		log.Error("Do cmd failed: \n", err.Error()+stdErr1.String())
		return -1,-1,err
	}
	s = strings.Split(stdOut1.String(), "\"")[1]
	rbytes, err := strconv.Atoi(s)
	if err != nil {
		log.Error("get num_file failed: \n", err.Error())
		return -1,-1,err
	}
	return rfiles,rbytes,nil

}

func GetAllDir(path string, level int, tar_level int) []string {
	var ans []string
	if level == tar_level {
		ch_dir := listFile(path)
		return ch_dir
	}
	paths := listFile(path)
	for _, i := range paths {
		ch_dir := listFile(i)
		ans = append(ans, ch_dir...)

	}
	return ans
}

func listFile(myfolder string) []string {
	var ans []string
	files, _ := ioutil.ReadDir(myfolder)
	for _, fd := range files {
		if fd.IsDir() {
			ans = append(ans, fd.Name())
		}
	}
	return ans
}

func CounTarDir(paths []string) []string {
	//paths = []string{"/mnt/cephfs/algor-api/user/*/*","/mnt/cephfs/algor-api/dataset/user/*/*","/mnt/cephfs/algor-api/dataset/public/*"}
	var ans []string
	for _, i := range paths {
		tar_dir := GetAllDir(i, 0, strings.Count(i, "*"))
		ans = append(ans, tar_dir...)
	}
	return ans
}
