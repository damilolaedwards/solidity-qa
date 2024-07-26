package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const OpenAIModel = "gpt-4-turbo"
const TokenLimitExceeded = "the number of tokens exceeds the maximum, please reduce the amount of data you are sending"

var apiKey = os.Getenv("OPENAI_API_KEY")

// TODO: Allow the cancelling of a request
func AskGPT4Turbo(messages []Message) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	url := "https://api.openai.com/v1/chat/completions"

	// Calculate the number of tokens to make sure we don't go over the limit
	numTokens, err := calculateNumTokens(messages)
	if err != nil {
		return "", err
	}
	if numTokens > 128000 {
		return "", fmt.Errorf(TokenLimitExceeded)
	}

	requestBody := ChatRequest{
		Model:    OpenAIModel,
		Messages: messages,
	}

	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response ChatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	var errorResponse ChatErrorResponse
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		return "", err
	}

	if errorResponse.Error.Message != "" {
		return "", fmt.Errorf("OpenAI API error: %s", errorResponse.Error.Message)
	}

	return "", fmt.Errorf("no response from OpenAI API")
}

// UploadFilesToOpenAI uploads multiple files to OpenAI and returns their IDs.
func UploadFilesToOpenAI(filePaths []string) ([]string, error) {
	// Ensure all provided files exist, are not directories and do not exceed the maximum file size limit
	var fileSizes int64 = 0
	for _, filePath := range filePaths {
		info, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("file %s does not exist", filePath)
			}
			return nil, fmt.Errorf("failed to stat file %s: %w", filePath, err)
		}

		if info.IsDir() {
			return nil, fmt.Errorf("file %s is a directory", filePath)
		}

		if info.Size() > 512*1024*1024 {
			return nil, fmt.Errorf("file %s size exceeds the maximum limit of 512 MB", filePath)
		}

		fileSizes += info.Size()
		if fileSizes > 100*1000*1024*1024 {
			return nil, fmt.Errorf("file size exceeds the maximum limit of 100 GB")
		}
	}

	var fileIDs []string

	for _, filePath := range filePaths {
		fileID, err := UploadFileToOpenAI(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file %s: %w", filePath, err)
		}
		fileIDs = append(fileIDs, fileID)
	}

	return fileIDs, nil
}

// UploadFileToOpenAI uploads a file at a specified path to OpenAI and returns the file ID.
func UploadFileToOpenAI(filePath string) (string, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file %s does not exist", filePath)
		}
		return "", fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("file %s is a directory", filePath)
	}

	if info.Size() > 512*1024*1024 {
		return "", fmt.Errorf("file %s size exceeds the maximum limit of 512 MB", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file: ", err)
		}
	}(file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	// Add file purpose
	err = writer.WriteField("purpose", "assistants")
	if err != nil {
		return "", fmt.Errorf("failed to write field: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/files", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body: ", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.ID, nil
}

// DeleteFileFromOpenAI deletes a file from OpenAI using the given file ID.
func DeleteFileFromOpenAI(fileID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.openai.com/v1/files/%s", fileID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body: ", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func calculateNumTokens(messages []Message) (int, error) {
	var numTokens = 0

	tkm, err := tiktoken.EncodingForModel(OpenAIModel)
	if err != nil {
		return numTokens, fmt.Errorf("unable to get model encoding: %v", err)
	}

	for _, message := range messages {
		token := tkm.Encode(message.Content, nil, nil)
		numTokens += len(token)
	}

	return numTokens, nil
}
