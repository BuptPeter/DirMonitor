package model

type UserInfo struct {
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
