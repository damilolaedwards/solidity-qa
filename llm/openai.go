package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func AskGPT4Turbo(messages []Message) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	url := "https://api.openai.com/v1/chat/completions"

	requestBody := ChatRequest{
		Model:    "gpt-4-turbo",
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
