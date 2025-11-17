package response

import "time"

// UserListItem 用户列表项
type UserListItem struct {
	ID         int64     `json:"id"`
	UserName   string    `json:"userName"`
	NickName   string    `json:"nickName"`
	Avatar     string    `json:"avatar"`
	Status     string    `json:"status"`
	UserGender string    `json:"userGender"`
	UserPhone  string    `json:"userPhone"`
	UserEmail  string    `json:"userEmail"`
	UserRoles  []string  `json:"userRoles"`
	CreateBy   string    `json:"createBy"`
	CreateTime time.Time `json:"createTime"`
	UpdateBy   string    `json:"updateBy"`
	UpdateTime time.Time `json:"updateTime"`
}

