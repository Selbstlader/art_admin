package response

import "time"

// AppUserListItem APP用户列表项
type AppUserListItem struct {
	ID            int64      `json:"id"`
	UserName      string     `json:"userName"`
	NickName      string     `json:"nickName"`
	Phone         string     `json:"phone"`
	Email         string     `json:"email"`
	Avatar        string     `json:"avatar"`
	UserType      string     `json:"userType"`
	Status        string     `json:"status"`
	LastLoginTime *time.Time `json:"lastLoginTime"`
	LastLoginIP   string     `json:"lastLoginIp"`
	CreateBy      string     `json:"createBy"`
	CreateTime    time.Time  `json:"createTime"`
	UpdateBy      string     `json:"updateBy"`
	UpdateTime    time.Time  `json:"updateTime"`
}
