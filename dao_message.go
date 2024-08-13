package gptapi

import "encoding/json"

type IMessage interface{}

type IToolChoice interface {
	Contents() string // 工具選擇的內文
}

type IContent interface{}

// 提示與影響AI風格訊息
type SystemMessage struct {
	Name    string   `json:"name,omitempty"` // 用來區分相同 role 下不同的參與者
	Role    string   `json:"role"`           // 訊息來源角色
	Content IContent `json:"content"`        // 內文
}

type UserMessage struct {
	Name    string   `json:"name,omitempty"` // 用來區分相同 role 下不同的參與者
	Role    string   `json:"role"`           // 訊息來源角色
	Content IContent `json:"content"`        // 內文
}

func (self *UserMessage) Contents() string { return "" }

// AI 回應訊息
type AssistantMessage struct {
	Name      string      `json:"name,omitempty"`    // 用來區分相同 role 下不同的參與者
	Role      string      `json:"role"`              // 訊息來源角色
	Content   string      `json:"content"`           // 內文
	Refusal   string      `json:"refusal,omitempty"` // 拒絕回應原因
	ToolCalls []ToolCalls `json:"tool_calls"`        // 本次回應AI使用到的 tool 調用資訊
}

// Tool 處理訊息
type ToolMessage struct {
	Role       string   `json:"role"`         // 訊息來源角色
	Content    IContent `json:"content"`      // 內文
	ToolCallId string   `json:"tool_call_id"` // 回應 tool 的索引Id
}

/////// IContent 實作區塊

// IContent實作 單一文字api訊息結構
type ContentText struct {
	Role    string `json:"role"`              // 訊息來源角色 EX: "system","user"
	Content string `json:"content,omitempty"` // 內文
}

func (self *ContentText) Contents() string {
	js, _ := json.Marshal(self)
	return string(js)
}

// IContent實作 包含圖片內文
type ContentImage struct {
	Type     string            `json:"type"`                // 文類型
	Text     string            `json:"text,omitempty"`      // 內文
	ImageURL *ContentImageData `json:"image_url,omitempty"` // 圖片內容
}

func (self *ContentImage) Contents() string {
	js, _ := json.Marshal(self)
	return string(js)
}

// 圖片內文結構
type ContentImageData struct {
	URL    string `json:"url"`    // 圖片"網址"或"base64編碼圖片"
	Detail string `json:"detail"` // 圖像模式
}

/////// IToolChoice 實作區塊

type ToolChoiceString string

func (self *ToolChoiceString) Contents() string {
	return string(*self)
}

type ToolChoiceObject struct {
	Type               string             `json:"type"`
	ToolChoiceFunction ToolChoiceFunction `json:"function"`
}

type ToolChoiceFunction struct {
	Name string `json:"name"`
}

///////

// Tool 定義工具的結構
type Tool struct {
	Type         string       `json:"type"`     // 工具類型目前只有 "function"
	ToolFunction ToolFunction `json:"function"` // 函數定義
	Strict       bool         `json:"strict"`   // 生成函數呼叫時是否啟用嚴格的架構遵循。當 strict 為 true 時，僅支援 JSON 模式的子集。
}

// ToolFunction 定義工具的函數結構
type ToolFunction struct {
	Name        string             `json:"name"`        // 函數名稱
	Description string             `json:"description"` // 函數描述
	Parameters  FunctionParameters `json:"parameters"`  // 函數參數
}

// FunctionParameters 定義函數參數的結構
type FunctionParameters struct {
	Type                 string               `json:"type"`                           // 函數參數類型
	Properties           map[string]Parameter `json:"properties"`                     // 參數目錄表
	Required             []string             `json:"required"`                       // 必要的參數名稱
	AdditionalProperties bool                 `json:"additionalProperties,omitempty"` // 是否比對為定義參數
}

// Parameter 定義函數參數的具體屬性
type Parameter struct {
	Type        string   `json:"type"`           // 資料類型
	Description string   `json:"description"`    // 資料描述
	Enum        []string `json:"enum,omitempty"` // 資料限定序列
}

// Usage 定義使用情況的結構體
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// 工具呼叫先關參數
type ToolCalls struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Function ToolCallsFunction `json:"function"`
}

type ToolCallsFunction struct {
	Name      string `json:"name"`      // 方法名稱
	Arguments string `json:"arguments"` // 輸入參數 json string
}

// IToolChoice實作 定義選擇的結構體
type ToolChoice struct {
	Message      AssistantMessage `json:"message"`       // 回應內容
	FinishReason string           `json:"finish_reason"` // 完成原因
	Index        int              `json:"index"`         // 索引值
}
