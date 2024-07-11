package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"io"
	"net/http"
	"os"
)

const OpenAIModel = "gpt-4-turbo"

func AskGPT4Turbo(messages []Message) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
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
		return "", fmt.Errorf("the number of tokens exceeds the maximum, please reduce the amount of data you are sending")
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
