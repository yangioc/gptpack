package gptapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// gpt-4o vs gpt-4o-mini 模型 token 價值不同, mini 消耗更多 token 但價格更便宜 4o VS 4o-mini 價格大約為 1:5
	model = "gpt-4o-mini"
)

func NewCompletionsRequest(maxToken int) completionsRequest {
	return completionsRequest{
		Model:     model,
		MaxTokens: maxToken,
	}
}

// 模型任務
func CompletionsRequest(apiKey string, reqBody completionsRequest) (*CompletionsResponse, error) {
	// Marshal the request body to JSON
	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, Url_Completions, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check if the response status is not OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", respBody)
	}

	var response CompletionsResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, fmt.Errorf("無法解析回應: %v", err)
	}
	return &response, nil
}

// 模型任務 以串流方式回應
func CompletionsStreamingRequest(apiKey string, reqBody completionsRequest, output chan<- string) error {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("Error marshalling request body: %v", err)
	}

	// 創建 HTTP 請求
	req, err := http.NewRequest(http.MethodPost, Url_Completions, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 發送請求並接收回應
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Non-OK HTTP status: %v", resp.Status)
	}

	// 讀取並處理流式數據
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			output <- line
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	close(output)
	return nil
}

// ///// 批次任務

// 建立新批次處理
func CreateBatchRequest(apiKey, inputFileId string) (*CreateBatchResponse, error) {
	requestBody := BatchRequest{
		InputFileID:      inputFileId, // Replace with your actual file ID
		Endpoint:         "/v1/chat/completions",
		CompletionWindow: "24h",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Error marshaling request data: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, Url_Batches, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	if resp.StatusCode != 200 {
		errRes := ErrorResponse{}
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}

		return nil, errors.New(errRes.Error.Message)

	} else {
		res := CreateBatchResponse{}
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, err
		}

		return &res, nil
	}
}

// 取得已存在批次任務
//
// @after 指定batchId 查詢會從此ID後的像是開始輸出
// @limit 指定輸出項目數量 default: 20 range: 1~100
func ListBatchRequest(apiKey, after string, limit int) (*ListBatchResponse, error) {

	if limit < 1 || 100 < limit {
		limit = 20
	}

	apiUrl, err := url.Parse(Url_ListBatch)
	if err != nil {
		return nil, err
	}

	// 創建查詢參數
	params := url.Values{}
	if after != "" {
		params.Add("after", "batch_10") // 添加 `after` 參數
	}
	params.Add("limit", strconv.Itoa(limit)) // 添加 `limit` 參數

	// 設置查詢參數到 URL 中
	apiUrl.RawQuery = params.Encode()
	req, err := http.NewRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		errRes := ErrorResponse{}
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}

		return nil, errors.New(errRes.Error.Message)

	} else {
		res := ListBatchResponse{}
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, err
		}

		return &res, nil
	}
}

// 查詢指定批次任務
func RetrieveBatchRequest(apiKey, batchId string) (*RetrieveBatchResponse, error) {
	url := strings.Replace(Url_RetrieveBatch, "{batch_id}", batchId, -1)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		errRes := ErrorResponse{}
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}

		return nil, errors.New(errRes.Error.Message)

	} else {
		res := RetrieveBatchResponse{}
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, err
		}

		return &res, nil
	}
}

// 取消指定批次任務
func CancelBatchRequest(apiKey, batchId string) (*CancelBatchResponse, error) {
	url := strings.ReplaceAll(Url_CancelBatch, "{batch_id}", batchId)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		errRes := ErrorResponse{}
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}

		return nil, errors.New(errRes.Error.Message)

	} else {
		res := CancelBatchResponse{}
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, err
		}

		return &res, nil
	}
}

// ///// 檔案處理任務

// 上傳檔案
func UploadFileRequest(apiKey, jsonlPath, purpose string) (*FileUploadResponse, error) {

	// apiKey := "YOUR_OPENAI_API_KEY"        // 替换为你的 OpenAI API 金钥
	// filePath := "path/to/your/file.jsonl"  // 替换为你要上传的文件路径
	// purpose := "fine-tune"                  // 替换为文件的用途

	// 打开文件
	file, err := os.Open(jsonlPath)
	if err != nil {
		return nil, fmt.Errorf("無法打開文件: %v", err)
	}
	defer file.Close()

	fs, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("檔案狀態異常: %v", err)
	} else if fs.Size() >= BatchFileSizeLimit { // 檔案大小檢查
		return nil, fmt.Errorf("檔案過大: %v", err)
	} else if filepath.Ext(fs.Name()) != ".jsonl" { // 只支援 .jsonl格式
		return nil, fmt.Errorf("[UploadFile] Error filetype filePath: %s", jsonlPath)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fs.Name())
	if err != nil {
		return nil, fmt.Errorf("無法建立文件字段: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("無法拷貝文件內容: %v", err)
	}

	err = writer.WriteField("purpose", purpose)
	if err != nil {
		return nil, fmt.Errorf("無法新增目的字段: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("無法關閉寫入器: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, Url_UploadFiles, body)
	if err != nil {
		return nil, fmt.Errorf("無法建立請求: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("請求失敗: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("無法讀取回應: %v", err)
	}

	var uploadResponse FileUploadResponse
	err = json.Unmarshal(respBody, &uploadResponse)
	if err != nil {
		return nil, fmt.Errorf("無法解析回應: %v", err)
	}

	return &uploadResponse, nil
}

// 檔案列表查詢
func ListFileRequest(apiKey string) (*ListFileResponse, error) {
	req, err := http.NewRequest(http.MethodGet, Url_ListFiles, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		errRes := ErrorResponse{}
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}

		return nil, errors.New(errRes.Error.Message)

	} else {
		res := ListFileResponse{}
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, err
		}

		return &res, nil
	}
}

// 檔案資訊查詢
func RetrieveFileRequest(apiKey, fileId string) (*RetrieveFileResponse, error) {

	url := strings.ReplaceAll(Url_RetrieveFile, "{file_id}", fileId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		errRes := ErrorResponse{}
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}

		return nil, errors.New(errRes.Error.Message)

	} else {
		res := RetrieveFileResponse{}
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, err
		}

		return &res, nil
	}
}

// 刪除檔案
func DeleteFileRequest(apiKey, fileId string) (*DeleteFileResponse, error) {
	url := strings.ReplaceAll(Url_DeleteFile, "{file_id}", fileId)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		errRes := ErrorResponse{}
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}

		return nil, errors.New(errRes.Error.Message)

	} else {
		res := DeleteFileResponse{}
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, err
		}

		return &res, nil
	}
}

// 檢索檔案內文
func RetrieveFileContentRequest(apiKey, fileId string, batchType string) (*RetrieveFileContentResponse, error) {
	url := strings.ReplaceAll(Url_RetrueveFileContent, "{file_id}", fileId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		errRes := ErrorResponse{}
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, err
		}

		return nil, errors.New(errRes.Error.Message)

	} else {

		res := RetrieveFileContentResponse{}
		switch batchType {
		case BatchType_Completions:
			for _, data := range bytes.Split(body, []byte{'\n'}) {
				if len(data) == 0 {
					continue
				}
				rowData := BatchOutput{
					Response: BatchOutputResData{
						Body: &CompletionsResponse{},
					}}
				if err := json.Unmarshal(data, &rowData); err != nil {
					panic(err)
				}

				res.Data = append(res.Data, rowData)
			}

		case BatchType_Embeddings:
			for _, data := range bytes.Split(body, []byte{'\n'}) {
				if len(data) == 0 {
					continue
				}
				rowData := BatchOutput{
					Response: BatchOutputResData{
						// Body: &EmbeddingsResponse{},
					}}
				if err := json.Unmarshal(data, &rowData); err != nil {
					panic(err)
				}

				res.Data = append(res.Data, rowData)
			}
		}

		return &res, nil
	}
}
