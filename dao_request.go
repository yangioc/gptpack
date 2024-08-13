package gptapi

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// 錯誤詳細信息結構體
type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"` // 參數可能為 null
	Code    string `json:"code"`  // 錯誤代碼可能為 null
}

// Completions Request 請求結構
type completionsRequest struct {
	Model      string      `json:"model"`
	Messages   []IMessage  `json:"messages"`
	MaxTokens  int         `json:"max_tokens"`            // 最大 token 使用數量 (每個token大約能回傳4的文字的內文)
	Tools      []Tool      `json:"tools,omitempty"`       // 模型可能呼叫的工具列表。目前，僅支援函數。使用它來提供模型可以為其產生 JSON 輸入的函數列表。最多支援 128 個功能。
	ToolChoice IToolChoice `json:"tool_choice,omitempty"` //
}

// Completions Response 回應結構
type CompletionsResponse struct {
	ID                 string       `json:"id"`
	Object             string       `json:"object"`
	Created            int          `json:"created"` // 完成時間
	Model              string       `json:"model"`   // 本次請求指定模型
	Usage              Usage        `json:"usage"`   // token 使用紀錄
	Choices            []ToolChoice `json:"choices"` // 模型完成後返回的清單
	Service_tier       string       `json:"service_tier,omitempty"`
	System_fingerprint string       `json:"system_fingerprint"`
}

func (self *completionsRequest) AddMessage(message IMessage) {
	self.Messages = append(self.Messages, message)
}

func (self *completionsRequest) AddTools(tools []Tool) {
	self.Tools = append(self.Tools, tools...)
}

// 啟動批次任務回應
type CreateBatchResponse struct {
	BatchInfo
}

// 查詢批次列表回應
type ListBatchResponse struct {
	Object  string      `json:"object"`   // 固定為 "list"
	Data    []BatchInfo `json:"data"`     // 批次數據列表
	FirstID string      `json:"first_id"` // 第一筆批次ID
	LastID  string      `json:"last_id"`  // 最後一筆批次ID
	HasMore bool        `json:"has_more"` // 後續是否還有資料
}

// 查詢批次任務回應
type RetrieveBatchResponse struct {
	BatchInfo
}

// 取消批次任務回應
type CancelBatchResponse struct {
	BatchInfo
}

// 檔案上傳回應
type FileUploadResponse struct {
	FileInfo
}

// 檔案查詢列表回應
type ListFileResponse struct {
	Object string     `json:"object"` // 固定為 "list"
	Data   []FileInfo `json:"data"`   // 文件信息列表
}

// 檔案資訊查詢回應
type RetrieveFileResponse struct {
	FileInfo
}

// Delete檔案回應
type DeleteFileResponse struct {
	Id      string `json:"id"`      // 檔案名稱
	Object  string `json:"object"`  // 物件類型
	Deleted bool   `json:"deleted"` // 是否刪除
}

type RetrieveFileContentResponse struct {
	Data []BatchOutput // 每列資料
}
