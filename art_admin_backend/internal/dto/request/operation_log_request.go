package request

// OperationLogListRequest 操作日志列表请求
type OperationLogListRequest struct {
	Current      int    `form:"current" binding:"required,min=1" example:"1"`
	Size         int    `form:"size" binding:"required,min=1,max=1000" example:"10"`
	Module       string `form:"module" json:"module"`             // 操作模块
	BusinessType string `form:"businessType" json:"businessType"` // 业务类型
	OperatorName string `form:"operatorName" json:"operatorName"` // 操作人员
	Status       *int   `form:"status" json:"status"`             // 操作状态
	StartTime    string `form:"startTime" json:"startTime"`       // 开始时间
	EndTime      string `form:"endTime" json:"endTime"`           // 结束时间
}

// CreateOperationLogRequest 创建操作日志请求
type CreateOperationLogRequest struct {
	Module        string `json:"module" binding:"required"`       // 操作模块
	BusinessType  string `json:"businessType" binding:"required"` // 业务类型
	RequestMethod string `json:"requestMethod"`                   // 请求方式
	RequestURL    string `json:"requestUrl"`                      // 请求URL
	OperatorName  string `json:"operatorName"`                    // 操作人员
	OperatorIP    string `json:"operatorIp"`                      // 操作IP
	OperatorAddr  string `json:"operatorAddr"`                    // 操作地点
	RequestParam  string `json:"requestParam"`                    // 请求参数
	ResponseData  string `json:"responseData"`                    // 响应数据
	Status        int    `json:"status"`                          // 操作状态
	ErrorMsg      string `json:"errorMsg"`                        // 错误消息
	CostTime      int64  `json:"costTime"`                        // 消耗时间
	UserAgent     string `json:"userAgent"`                       // 用户代理
}
