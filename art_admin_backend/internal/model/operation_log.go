package model

import (
	"time"

	"gorm.io/gorm"
)

// OperationLog 操作日志模型
type OperationLog struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Module        string         `gorm:"type:varchar(50);index;not null" json:"module"`     // 操作模块
	BusinessType  string         `gorm:"type:varchar(20);index" json:"businessType"`        // 业务类型(新增/修改/删除/查询/导出/导入等)
	RequestMethod string         `gorm:"type:varchar(10)" json:"requestMethod"`             // 请求方式(GET/POST/PUT/DELETE)
	RequestURL    string         `gorm:"type:varchar(500)" json:"requestUrl"`               // 请求URL
	OperatorName  string         `gorm:"type:varchar(50);index" json:"operatorName"`        // 操作人员
	OperatorIP    string         `gorm:"type:varchar(50);index" json:"operatorIp"`          // 操作IP
	OperatorAddr  string         `gorm:"type:varchar(200)" json:"operatorAddr"`             // 操作地点
	RequestParam  string         `gorm:"type:text" json:"requestParam"`                     // 请求参数
	ResponseData  string         `gorm:"type:text" json:"responseData"`                     // 响应数据
	Status        int            `gorm:"type:tinyint(1);default:1;index" json:"status"`     // 操作状态(1-成功 0-失败)
	ErrorMsg      string         `gorm:"type:varchar(2000)" json:"errorMsg"`                // 错误消息
	CostTime      int64          `gorm:"type:bigint" json:"costTime"`                       // 消耗时间(毫秒)
	UserAgent     string         `gorm:"type:varchar(500)" json:"userAgent"`                // 用户代理
	OperationTime time.Time      `gorm:"type:datetime;index;not null" json:"operationTime"` // 操作时间
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (OperationLog) TableName() string {
	return "sys_operation_log"
}
