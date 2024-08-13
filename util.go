package gptapi

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
)

func ImageEncode(imagepath string) (string, error) {
	extName := filepath.Ext(imagepath)
	if extName == "" {
		return "", fmt.Errorf("[EncodeLocationImage] Error imagepath: %s", imagepath)
	}

	if !slices.Contains(SupImage, extName) {
		return "", fmt.Errorf("[EncodeLocationImage] Error not support file type: %s", extName)
	}

	file, err := os.Open(imagepath)
	if err != nil {
		return "", fmt.Errorf("[EncodeLocationImage] Error open: %s ,err: %v", imagepath, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("[EncodeLocationImage] Error ReadAll: %s ,err: %v", imagepath, err)
	}

	b64Str := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:image/%s;base64,%s", extName[1:], b64Str), nil
}

// 將批次指令寫入檔案內
func NewJsonlFile(dirPath, filename string, records []Record) error {

	jsonlPath := fmt.Sprintf("%s/%s.jsonl", dirPath, filename)
	// 创建文件
	file, err := os.Create(jsonlPath)
	if err != nil {
		return fmt.Errorf("[WriteToJSONL] Error err: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, record := range records {
		recordJSON, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("[WriteToJSONL] Error err: %v", err)
		}

		if _, err := writer.Write(recordJSON); err != nil {
			return fmt.Errorf("[WriteToJSONL] Error err: %v", err)
		}
		_ = writer.WriteByte('\n')
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("[WriteToJSONL] Error err: %v", err)
	}

	return nil
}

// //// User Message 區塊
func NewUserTextMessage(text string) IMessage {
	return &UserMessage{
		Role:    MessageContentRole_User,
		Content: text,
	}
}

func NewUserImageMessage(text, imageUrl string) IMessage {
	return &UserMessage{
		Role: MessageContentRole_User,
		Content: []ContentImage{
			{
				Type: MessageContentType_Image,
				ImageURL: &ContentImageData{
					URL:    imageUrl,
					Detail: ImageDetailMode_Low,
				},
			}, {
				Type: MessageContentType_Text,
				Text: text,
			},
		},
	}
}

// //// System Message 區塊
func NewSystemTextMessage(text string) IMessage {
	return &SystemMessage{
		Role:    MessageContentRole_System,
		Content: text,
	}
}

// //// Assistant Message 區塊
func NewAssistantTextMessage(text string) IMessage {
	return &AssistantMessage{
		Role:    MessageContentRole_Assistant,
		Content: text,
	}
}

// //// Tool 區塊

func NewToolMessage(id string, content IContent) IMessage {
	js, _ := json.Marshal(content)
	return &ToolMessage{
		Role:       MessageContentRole_Tool,
		Content:    string(js),
		ToolCallId: id,
	}
}

// 生成新 Tool
//
// @name:外部呼叫名稱
// @description:api功能說明
// @paramet:api參數
func NewTool(name, description string, paramet FunctionParameters) Tool {
	return Tool{
		Type: "function",
		ToolFunction: ToolFunction{
			Name:        name,
			Description: description,
			Parameters:  paramet,
		},
	}
}

// 生成新 FunctionParameters
//
// [][3]string 0: key:參數名稱, type:參數資料型態, description:參數說明提供給openAI辨識
func NewToolFunctionParameters(parameterDatas [][3]string) FunctionParameters {

	params := FunctionParameters{
		Type:                 "object",
		Properties:           make(map[string]Parameter),
		AdditionalProperties: false,
	}
	for _, dataRow := range parameterDatas {
		params.Properties[dataRow[0]] = Parameter{
			Type:        dataRow[1],
			Description: dataRow[2],
		}
		params.Required = append(params.Required, dataRow[0])
	}

	return params
}
