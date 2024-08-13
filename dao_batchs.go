package gptapi

// 批次處理請求
type BatchRequest struct {
	InputFileID      string `json:"input_file_id"`     // 上傳的請求 jsonl 檔案Id
	Endpoint         string `json:"endpoint"`          // 批次指令要執行的目標網址
	CompletionWindow string `json:"completion_window"` // 完成時間當前限制 24H 後
}

// BatchInfo 代表一個批次對象
type BatchInfo struct {
	ID               string            `json:"id"`                       // 批次的唯一標識符
	Object           string            `json:"object"`                   // 物件的類型，應該是 "batch"
	Endpoint         string            `json:"endpoint"`                 // 使用的 OpenAI API 端點
	Errors           BatchErrors       `json:"errors,omitempty"`         // 批次中的錯誤信息
	InputFileID      string            `json:"input_file_id"`            // 輸入文件的 ID
	CompletionWindow string            `json:"completion_window"`        // 批次處理的時間範圍
	Status           string            `json:"status"`                   // 批次的當前狀態
	OutputFileID     string            `json:"output_file_id"`           // 成功請求的輸出文件 ID
	ErrorFileID      string            `json:"error_file_id"`            // 錯誤請求的輸出文件 ID
	CreatedAt        int64             `json:"created_at"`               // 批次創建時間的 Unix 時間戳
	InProgressAt     int64             `json:"in_progress_at,omitempty"` // 批次開始處理的 Unix 時間戳
	ExpiresAt        int64             `json:"expires_at,omitempty"`     // 批次過期的 Unix 時間戳
	FinalizingAt     int64             `json:"finalizing_at,omitempty"`  // 批次開始最終處理的 Unix 時間戳
	CompletedAt      int64             `json:"completed_at,omitempty"`   // 批次完成的 Unix 時間戳
	FailedAt         int64             `json:"failed_at,omitempty"`      // 批次失敗的 Unix 時間戳
	ExpiredAt        int64             `json:"expired_at,omitempty"`     // 批次過期的 Unix 時間戳
	CancellingAt     int64             `json:"cancelling_at,omitempty"`  // 批次開始取消的 Unix 時間戳
	CancelledAt      int64             `json:"cancelled_at,omitempty"`   // 批次取消的 Unix 時間戳
	RequestCounts    BatchTokenCounts  `json:"request_counts,omitempty"` // 批次中的請求數量統計
	Metadata         map[string]string `json:"metadata,omitempty"`       // 附加的元數據
}

// 批次處理結果資訊
type BatchOutput struct {
	ID       string             `json:"id"`
	CustomID string             `json:"custom_id"`
	Response BatchOutputResData `json:"response"` // 使用指標以處理可能的 null
	Error    BatchErrorData     `json:"error"`    // 使用指標以處理可能的 null
}

// 批次處理結果 任務資訊
type BatchOutputResData struct {
	StatusCode int         `json:"status_code"`
	RequestID  string      `json:"request_id"` // 任務ID
	Body       interface{} `json:"body"`       // 回應內文
}

// BatchErrors 代表批次中的錯誤信息
type BatchErrors struct {
	Object string           `json:"object"` // 物件的類型，應該是 "list"
	Data   []BatchErrorData `json:"data"`   // 錯誤的詳細信息
}

// BatchErrorData 代表單個錯誤對象
type BatchErrorData struct {
	Code    string `json:"code"`            // 錯誤代碼
	Message string `json:"message"`         // 錯誤信息
	Param   string `json:"param,omitempty"` // 導致錯誤的參數名稱（如果適用）
	Line    int    `json:"line,omitempty"`  // 輸入文件中的行號（如果適用）
}

// BatchTokenCounts 代表批次中的請求數量統計
type BatchTokenCounts struct {
	Total     int `json:"total"`     // 批次中的總請求數量
	Completed int `json:"completed"` // 已成功完成的請求數量
	Failed    int `json:"failed"`    // 失敗的請求數量
}

// 檔案查詢結果
type FileInfo struct {
	ID        string `json:"id"`
	Object    string `json:"object"` // 固定為 "file"
	Bytes     int    `json:"bytes"`
	CreatedAt int    `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

// 批次檔案規定格式 jsonl
type Record struct {
	CustomID string             `json:"custom_id"` // 自定義請求名稱
	Method   string             `json:"method"`    // http 傳輸方式
	URL      string             `json:"url"`       // api 路徑
	Body     completionsRequest `json:"body"`      // api 內容
}
