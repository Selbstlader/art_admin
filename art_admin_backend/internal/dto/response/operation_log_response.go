package response

import "time"

// OperationLogListItem 操作日志列表项
type OperationLogListItem struct {
	ID            int64     `json:"id"`
	Module        string    `json:"module"`
	BusinessType  string    `json:"businessType"`
	RequestMethod string    `json:"requestMethod"`
	RequestURL    string    `json:"requestUrl"`
	OperatorName  string    `json:"operatorName"`
	OperatorIP    string    `json:"operatorIp"`
	OperatorAddr  string    `json:"operatorAddr"`
	RequestParam  string    `json:"requestParam"`
	ResponseData  string    `json:"responseData"`
	Status        int       `json:"status"`
	ErrorMsg      string    `json:"errorMsg"`
	CostTime      int64     `json:"costTime"`
	UserAgent     string    `json:"userAgent"`
	OperationTime time.Time `json:"operationTime"`
}

// OperationLogDetail 操作日志详情
type OperationLogDetail struct {
	OperationLogListItem
}
