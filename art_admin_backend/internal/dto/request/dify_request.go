package request

// DifyDatasetListRequest 获取知识库列表请求
type DifyDatasetListRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`          // 页码
	Limit   int    `form:"limit" binding:"omitempty,min=1,max=100"` // 每页数量
	Keyword string `form:"keyword"`                                 // 搜索关键词
}

// DifyUploadFileRequest 上传文件到知识库请求
type DifyUploadFileRequest struct {
	DatasetID string `form:"dataset_id" binding:"required"` // 知识库ID
}

// DifyChatRequest AI对话请求
type DifyChatRequest struct {
	Query            string                 `json:"query" binding:"required"`     // 用户问题
	ConversationID   string                 `json:"conversation_id,omitempty"`    // 会话ID
	User             string                 `json:"user" binding:"required"`      // 用户标识
	ResponseMode     string                 `json:"response_mode,omitempty"`      // 响应模式: blocking/streaming (内部设置)
	Inputs           map[string]interface{} `json:"inputs,omitempty"`             // 输入变量
	Files            []DifyFileInfo         `json:"files,omitempty"`              // 文件列表
	AutoGenerateName bool                   `json:"auto_generate_name,omitempty"` // 是否自动生成标题
}

// DifyFileInfo 文件信息
type DifyFileInfo struct {
	Type           string `json:"type"`                     // 文件类型: image/document
	TransferMethod string `json:"transfer_method"`          // 传输方式: remote_url/local_file
	URL            string `json:"url,omitempty"`            // 远程URL
	UploadFileID   string `json:"upload_file_id,omitempty"` // 上传文件ID
}
