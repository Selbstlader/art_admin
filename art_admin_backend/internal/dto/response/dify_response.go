package response

// DifyDatasetListResponse 知识库列表响应
type DifyDatasetListResponse struct {
	Data    []DifyDataset `json:"data"`
	HasMore bool          `json:"has_more"`
	Limit   int           `json:"limit"`
	Total   int           `json:"total"`
	Page    int           `json:"page"`
}

// DifyDataset 知识库信息
type DifyDataset struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Description            string `json:"description"`
	Provider               string `json:"provider"`
	Permission             string `json:"permission"`
	DataSourceType         string `json:"data_source_type"`
	IndexingTechnique      string `json:"indexing_technique"`
	AppCount               int    `json:"app_count"`
	DocumentCount          int    `json:"document_count"`
	WordCount              int    `json:"word_count"`
	CreatedBy              string `json:"created_by"`
	CreatedAt              int64  `json:"created_at"`
	UpdatedBy              string `json:"updated_by"`
	UpdatedAt              int64  `json:"updated_at"`
	EmbeddingModel         string `json:"embedding_model"`
	EmbeddingModelProvider string `json:"embedding_model_provider"`
	EmbeddingAvailable     bool   `json:"embedding_available"`
}

// DifyDatasetDetailResponse 知识库详情响应
type DifyDatasetDetailResponse struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Description            string `json:"description"`
	Provider               string `json:"provider"`
	Permission             string `json:"permission"`
	DataSourceType         string `json:"data_source_type"`
	IndexingTechnique      string `json:"indexing_technique"`
	AppCount               int    `json:"app_count"`
	DocumentCount          int    `json:"document_count"`
	WordCount              int    `json:"word_count"`
	CreatedBy              string `json:"created_by"`
	CreatedAt              int64  `json:"created_at"`
	UpdatedBy              string `json:"updated_by"`
	UpdatedAt              int64  `json:"updated_at"`
	EmbeddingModel         string `json:"embedding_model"`
	EmbeddingModelProvider string `json:"embedding_model_provider"`
	EmbeddingAvailable     bool   `json:"embedding_available"`
}

// DifyUploadFileResponse 上传文件响应
type DifyUploadFileResponse struct {
	Document DifyDocument `json:"document"`
	Batch    string       `json:"batch"`
}

// DifyDocument 文档信息
type DifyDocument struct {
	ID                 string                 `json:"id"`
	Position           int                    `json:"position"`
	DataSourceType     string                 `json:"data_source_type"`
	DataSourceInfo     map[string]interface{} `json:"data_source_info"`
	DatasetProcessRule map[string]interface{} `json:"dataset_process_rule"`
	Name               string                 `json:"name"`
	CreatedFrom        string                 `json:"created_from"`
	CreatedBy          string                 `json:"created_by"`
	CreatedAt          int64                  `json:"created_at"`
	Tokens             int                    `json:"tokens"`
	IndexingStatus     string                 `json:"indexing_status"`
	Error              string                 `json:"error,omitempty"`
	Enabled            bool                   `json:"enabled"`
	DisabledAt         int64                  `json:"disabled_at,omitempty"`
	DisabledBy         string                 `json:"disabled_by,omitempty"`
	Archived           bool                   `json:"archived"`
	DisplayStatus      string                 `json:"display_status"`
	WordCount          int                    `json:"word_count"`
	HitCount           int                    `json:"hit_count"`
	DocForm            string                 `json:"doc_form"`
}

// DifyChatResponse AI对话响应(非流式)
type DifyChatResponse struct {
	Event          string                 `json:"event"`
	MessageID      string                 `json:"message_id"`
	ConversationID string                 `json:"conversation_id"`
	Mode           string                 `json:"mode"`
	Answer         string                 `json:"answer"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedAt      int64                  `json:"created_at"`
}

// DifyChatStreamEvent 流式响应事件
type DifyChatStreamEvent struct {
	Event          string                 `json:"event"` // 事件类型
	TaskID         string                 `json:"task_id,omitempty"`
	MessageID      string                 `json:"message_id,omitempty"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	Answer         string                 `json:"answer,omitempty"` // 消息内容
	CreatedAt      int64                  `json:"created_at,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// DifyErrorResponse 错误响应
type DifyErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// DifyConversationListResponse 会话列表响应
type DifyConversationListResponse struct {
	Limit   int                `json:"limit"`
	HasMore bool               `json:"has_more"`
	Data    []DifyConversation `json:"data"`
}

// DifyConversation 会话信息
type DifyConversation struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Inputs       map[string]interface{} `json:"inputs"`
	Status       string                 `json:"status"`
	Introduction string                 `json:"introduction"`
	CreatedAt    int64                  `json:"created_at"`
	UpdatedAt    int64                  `json:"updated_at"`
}

// DifyMessageListResponse 消息历史响应
type DifyMessageListResponse struct {
	Limit   int           `json:"limit"`
	HasMore bool          `json:"has_more"`
	Data    []DifyMessage `json:"data"`
}

// DifyMessage 消息信息
type DifyMessage struct {
	ID                 string                  `json:"id"`
	ConversationID     string                  `json:"conversation_id"`
	Inputs             map[string]interface{}  `json:"inputs"`
	Query              string                  `json:"query"`
	Answer             string                  `json:"answer"`
	MessageFiles       []DifyMessageFile       `json:"message_files"`
	Feedback           *DifyFeedback           `json:"feedback,omitempty"`
	RetrieverResources []DifyRetrieverResource `json:"retriever_resources"`
	CreatedAt          int64                   `json:"created_at"`
}

// DifyMessageFile 消息文件
type DifyMessageFile struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	BelongsTo string `json:"belongs_to"`
}

// DifyFeedback 反馈信息
type DifyFeedback struct {
	Rating string `json:"rating"`
}

// DifyRetrieverResource 检索资源
type DifyRetrieverResource struct {
	Position     int     `json:"position"`
	DatasetID    string  `json:"dataset_id"`
	DatasetName  string  `json:"dataset_name"`
	DocumentID   string  `json:"document_id"`
	DocumentName string  `json:"document_name"`
	SegmentID    string  `json:"segment_id"`
	Score        float64 `json:"score"`
	Content      string  `json:"content"`
}

// DifySuggestedQuestionsResponse 建议问题响应
type DifySuggestedQuestionsResponse struct {
	Result string   `json:"result"`
	Data   []string `json:"data"`
}

// DifyStopResponse 停止响应
type DifyStopResponse struct {
	Result string `json:"result"`
}
