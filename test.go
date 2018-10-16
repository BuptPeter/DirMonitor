package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
)
type UserInfo_T struct {
	Email     string `json:"email"`
	Enabled   bool   `json:"enabled"`
	Erp       string `json:"erp"`
	GroupList []struct {
		CPU            int    `json:"cpu"`
		Erp            string `json:"erp"`
		GroupType      int    `json:"group_type"`
		ID             string `json:"id"`
		Memory         int    `json:"memory"`
		Name           string `json:"name"`
		NotebookStatus int    `json:"notebook_status"`
		NvidiaGpu      int    `json:"nvidia_gpu"`
		Office         string `json:"office"`
		SyncStatus     int    `json:"sync_status"`
		UseGroupID     string `json:"use_group_id"`
	} `json:"groupList"`
	ID           string `json:"id"`
	IsAdmin      bool   `json:"isAdmin"`
	IsDelete     int    `json:"isDelete"`
	IsGroupAdmin int    `json:"isGroupAdmin"`
	RealName     string `json:"realName"`
}

func init() {

}

func main_test() {
	var user_info UserInfo_T
	url:="http://authintf.jd.com/access-auth/getUserById?userId=91fe6b5e-0f07-11e8-914e-6eb2663d5452"
	resp,err:=http.Get(url)
	if err != nil {
		log.Error("Do http Get failed:", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	//body := []byte(`{"email":"wangyiying@jd.com","enabled":true,"erp":"wangyiying6","groupList":[{"cpu":5000,"erp":"zhuertao","group_type":1,"id":"b0044583-6f7e-11e8-b52e-fa163ee59f29","memory":4804,"name":"aiplatform","notebook_status":1,"nvidia_gpu":200,"office":"AI平台与研究部-基础平台部","sync_status":0,"use_group_id":""}],"id":"91fe6b5e-0f07-11e8-914e-6eb2663d5452","isAdmin":false,"isDelete":0,"isGroupAdmin":0,"realName":"王艺颖"}`)
	log.Info(string(body))
	json.NewDecoder(resp.Body).Decode(&user_info)

	if err := json.Unmarshal(body,&user_info); err != nil {
		fmt.Println(err)
		return
	}
	log.Info(user_info)
	log.Info("GetUserInfo success:",user_info.RealName,user_info.Erp,user_info.Email)
}
