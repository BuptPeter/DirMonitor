package main

import (
	"DirMonitor/model"
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)
var user_info_cache map[string]model.UserInfo
var dir_all_cache []string
func IsMatch(path string, files_limit int, byte_limit int) (int, map[string]int) {

	var info map[string]int
	info = make(map[string]int)
	num_file, num_byte, err := GetDirInfo(path)
	if err != nil {
		log.Error("Do GetDirInfo failed: \n", err.Error())
		return -1, nil
	}
	log.Info("GetDirInfo success,info:",num_file, num_byte)
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
func GetUserInfo(path string)( user_info model.UserInfo,err error) {

	if ! strings.Contains(path,"user"){
		return user_info,nil
	}
	sub_str_path:= strings.Split(path,"/")
	uuid:= sub_str_path[len(sub_str_path)-2]
	log.Info("GetUserInfo:",path,"\t uuid:",uuid)
	url := "http://authintf.jd.com/access-auth/getUserById?userId=" + uuid
	resp, err := http.Get(url)
	if err != nil {
		log.Error("Do http Get failed:", err.Error())
		return user_info,nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	//json.NewDecoder(resp.Body).Decode(&user_info)
	if err := json.Unmarshal(body,&user_info); err != nil {
		log.Info("Json Unmarshal failed:",err.Error())
		return user_info,nil
	}
	log.Info("GetUserInfo success,","Name:",user_info.RealName,"\tErp Num:",user_info.Erp,"\tEmail:",user_info.Email)
	return user_info,nil
}
func PutInfo2Cache(){
	paths := CounTarDir(strings.Split(arg_paths, ","))
	user_info_cache = make(map[string]model.UserInfo)
	for _,path := range(paths)  {
		user_info,err:= GetUserInfo(path)
		if err != nil{
			log.Error("GetUserInfo failed:",err.Error())
		}
		user_info_cache[path] = user_info
	}



}
func GetDirInfo(path string) (int, int, error) {
	log.Info("GetDirInfo path:" + path)

	cmd_get_files := "getfattr -d -n ceph.dir.rfiles " + path

	var stdErr0, stdOut0 bytes.Buffer
	cmd0 := exec.Command("/bin/bash","-c",cmd_get_files)
	cmd0.Stderr = &stdErr0
	cmd0.Stdout = &stdOut0
	err := cmd0.Run()
	if err != nil { //执行getfattr失败,返回错误信息及标准输出
		log.Error("Do cmd failed: \n", err.Error()+stdErr0.String())
		return -1, -1, err
	}
	s := strings.Split(stdOut0.String(), "\"")[1]
	rfiles, err := strconv.Atoi(s)
	if err != nil {
		log.Error("get num_file failed: \n", err.Error())
		return -1, -1, err
	}

	var stdErr1, stdOut1 bytes.Buffer
	cmd_get_bytes := "getfattr -d -n ceph.dir.rbytes " + path
	cmd1 := exec.Command("/bin/bash","-c",cmd_get_bytes)
	cmd1.Stderr = &stdErr1
	cmd1.Stdout = &stdOut1
	err = cmd1.Run()
	if err != nil { //执行getfattr失败
		log.Error("Do cmd failed: \n", err.Error()+stdErr1.String())
		return -1, -1, err
	}
	s = strings.Split(stdOut1.String(), "\"")[1]
	rbytes, err := strconv.Atoi(s)
	if err != nil {
		log.Error("get num_file failed: \n", err.Error())
		return -1, -1, err
	}
	return rfiles, rbytes, nil

}

func GetAllDir(path string, level int, tar_level int) []string {
	var ans []string
	if level == tar_level-1 {
		ch_dir := listFile(path)
		return ch_dir
	}
	paths := listFile(path)
	//log.Info("ch_dir:" + path)
	for _, i := range paths {
		//log.Info("Path:",i)
		ch_dir := GetAllDir(i,level+1,tar_level)
		//log.Info("Ch_dir:",ch_dir)
		ans = append(ans, ch_dir...)
	}
	return ans
}

func listFile(myfolder string) []string {
	var ans []string
	files, _ := ioutil.ReadDir(myfolder)
	for _, fd := range files {
		if fd.IsDir() {
			ans = append(ans, myfolder+"/"+fd.Name())
		}
	}
	//log.Info("Root Files:",myfolder,"Files:",files,"  List Files :",ans)
	return ans
}

func CounTarDir(paths []string) []string {
	//paths = []string{"/mnt/cephfs/algor-api/user/*/*","/mnt/cephfs/algor-api/dataset/user/*/*","/mnt/cephfs/algor-api/dataset/public/*"}
	var ans []string
	log.Info("input attr - paths:",paths)
	for _, i := range paths {
		path:= strings.Split(i, "*")[0]
		path = path[0:len(path)-1]
		tar_level:=strings.Count(i, "*")
		log.Info("GetAllDir - path:",path)
		tar_dir := GetAllDir(path, 0,tar_level)
		ans = append(ans, tar_dir...)
		log.Info("Path:",path,"  tar_level:",tar_level,"  All Ch_Dir:",ans)
	}
	return ans
}

//func main(){
//	//files, _ := filepath.Glob("*")
//	files := listFile("/Users/wangyanlei3/go")
//	fmt.Print(files)
//	paths := CounTarDir([]string{"/Users/wangyanlei3/*/*","/dev/*"})
//	log.Info(paths)
//	for k, i := range paths {
//		log.Info(k, i)
//	}
//}