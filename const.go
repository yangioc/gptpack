package gptapi

const (

	// 低解析度圖像
	// 限制 Size: 512x512 消耗Token: 85
	ImageDetailMode_Low string = "low"

	// 高解析度圖像
	// 在低解析度的限制下
	// 為每個額外的 512x512 size 額外消耗 170 token 圖像建立資料
	// 限制 Size 上限:  768x2000
	ImageDetailMode_High string = "high"

	// 圖片容量大小限制
	ImageSizeLimit int = 20 * 1024 * 1024

	MessageFinishType_Stop          string = "stop"           // API 傳回完整訊息，或透過stop參數提供的停止序列之一終止的訊息
	MessageFinishType_Length        string = "length"         // 由於max_tokens參數或令牌限制，模型輸出不完整
	MessageFinishType_FunctionCall  string = "function_call"  // 模型決定呼叫一個函數
	MessageFinishType_ContentFilter string = "content_filter" // 由於我們的內容過濾器中的標記而省略了內容
	MessageFinishType_Null          string = "null"           // API 回應仍在進行中或不完整

	MessageContentRole_System    string = "system"    // 提示內文,用來影響AI的行為與回應風格
	MessageContentRole_User      string = "user"      // 使用者輸入內容
	MessageContentRole_Assistant string = "assistant" // 模型的實際回應
	MessageContentRole_Tool      string = "tool"      // 工具處理訊息

	MessageContentType_Text  string = "text"
	MessageContentType_Image string = "image_url"

	Url_Batches             string = "https://api.openai.com/v1/batches"
	Url_ListBatch           string = "https://api.openai.com/v1/batches"                   // 查詢已存在的批次任務
	Url_RetrieveBatch       string = "https://api.openai.com/v1/batches/{batch_id}"        // 查詢指定批次任務
	Url_CancelBatch         string = "https://api.openai.com/v1/batches/{batch_id}/cancel" // 取消指定批次任務
	Url_Completions         string = "https://api.openai.com/v1/chat/completions"          // 模型演算
	Url_UploadFiles         string = "https://api.openai.com/v1/files"                     // 上傳檔案
	Url_ListFiles           string = "https://api.openai.com/v1/files"                     // 取得檔案列表
	Url_RetrieveFile        string = "https://api.openai.com/v1/files/{file_id}"           // 檢索檔案資訊
	Url_DeleteFile          string = "https://api.openai.com/v1/files/{file_id}"           // 刪除檔案
	Url_RetrueveFileContent string = "https://api.openai.com/v1/files/{file_id}/content"   // 檢索檔案內文
)

// 批次處理目的標籤
const (

	// 檔案上傳大小限制
	BatchFileSizeLimit int64 = 1024 * 1024 * 1024 // 1G

	// 用途: 用於微調模型。文件中的資料將用於訓練或微調自訂模型。
	// 說明: 通常包含大量的訓練數據，如問答對或文字數據，用於改善模型在特定任務上的表現。
	BatchPurpose_Fine_tune string = "fine-tune"
	// 用途: 用於建置或訓練虛擬助理。文件中的資料將用於訓練智慧助理以執行特定任務或對話。
	// 說明: 包含對話資料或互動範例，用於提升虛擬助理的對話能力。
	BatchPurpose_Assistants string = "assistants"
	// 用途: 用於批次處理。文件中的資料將用於批次處理請求，這些請求可能涉及大規模的模型呼叫。
	// 說明: 適用於需要處理大量資料或請求的場景，批次上傳檔案並進行後續處理。
	BatchPurpose_Batch string = "batch"
	// 用途: 用於用戶資料。文件中的資料將用於儲存或處理與使用者相關的資訊。
	// 說明: 可用於記錄使用者互動、歷史資料或個人化資訊。
	BatchPurpose_User_data string = "user_data"
	// 	用途: 用於儲存回應資料。文件中的數據將用於記錄或分析模型產生的回應。
	// 說明: 包含模型產生的回應數據，用於後續分析或評估。
	BatchPurpose_Responses string = "responses"
	// 用途: 用於影像資料處理。文件中的數據將用於與視覺相關的任務，如圖像識別或分析。
	// 說明: 包含圖像或視覺數據，用於訓練或評估視覺模型。
	BatchPurpose_Vision string = "vision"

	// Batch 任務類型 用於指定如何解析內文
	BatchType_Completions string = "completions"
	BatchType_Embeddings  string = "embeddings"
)

// 'fine-tune', 'assistants', 'batch', 'user_data', 'responses', 'vision'
var (
	SupImage = []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}
)
