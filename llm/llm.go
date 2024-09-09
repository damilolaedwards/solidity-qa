package llm

import (
	"assistant/api/client"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"os"
)

type Model struct {
	Name        string
	Identifier  string
	Generates   string
	MaxTokenLen int
	url         string
	headers     map[string]string
}

const DefaultModelIdentifier = "chatgpt-4o-latest"

var Models = map[string]Model{
	"chatgpt-4o-latest": {
		Name:       "ChatGPT 4o Latest",
		Identifier: "chatgpt-4o-latest",
		Generates:  "text",
		url:        "https://api.openai.com/v1/chat/completions",
		headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")),
		},
		MaxTokenLen: 128000,
	},
	"dall-e-3": {
		Name:       "DALL·E 3",
		Identifier: "dall-e-3",
		Generates:  "image",
		url:        "https://api.openai.com/v1/images/generations",
		headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")),
		},
		MaxTokenLen: 76800,
	},
	"claude-3-5-sonnet-20240620": {
		Name:       "Claude 3.5 Sonnet",
		Identifier: "claude-3-5-sonnet-20240620",
		Generates:  "text",
		url:        "https://api.anthropic.com/v1/messages",
		headers: map[string]string{
			"Content-Type":      "application/json",
			"x-api-key":         os.Getenv("CLAUDE_API_KEY"),
			"anthropic-version": "2023-06-01",
		},
		MaxTokenLen: 1048576,
	},
}

const TokenLimitExceeded = "the number of tokens exceeds the maximum, please reduce the amount of data you are sending"

func AskModel(messages []ApiMessage, model string, ctx context.Context) (string, error) {
	for {
		select {
		case <-ctx.Done():
			return "", nil
		default:
			m, err := getModel(model)
			if err != nil {
				return "", err
			}

			numTokens, err := CalculateNumTokens(messages)
			if err != nil {
				return "", err
			}
			if numTokens > m.MaxTokenLen {
				return "", fmt.Errorf(TokenLimitExceeded)
			}

			var requestBody any

			if m.Generates == "image" {
				requestBody = ImageGenerationRequest{
					Model:  m.Identifier,
					Prompt: messages[len(messages)-1].Content,
					N:      1,
					Size:   "1024x1024",
				}
			} else {
				requestBody = TextGenerationRequest{
					Model:     m.Identifier,
					Messages:  messages,
					MaxTokens: 3000,
				}
			}

			req, err := client.CreateRequest(m.url, requestBody, "POST", m.headers)
			if err != nil {
				return "", err
			}

			resp, err := client.DoRequest(req)
			if err != nil {
				return "", err
			}

			return handleResponse(m, resp)
		}
	}
}

func getModel(model string) (Model, error) {
	if model == "" {
		model = DefaultModelIdentifier
	}
	m, ok := Models[model]
	if !ok {
		return Model{}, fmt.Errorf("unknown model: %s", model)
	}
	return m, nil
}

func handleResponse(m Model, resp *http.Response) (string, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if m.Generates == "image" {
		return handleImageResponse(body)
	}

	if m.Identifier == "claude-3-5-sonnet-20240620" {
		return handleClaudeTextResponse(body)
	} else {
		return handleOpenAITextResponse(body)
	}
}

func handleImageResponse(body []byte) (string, error) {
	var response ImageGenerationResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if len(response.Data) > 0 {
		return response.Data[0].URL, nil
	}
	var errorResponse OpenAIErrorResponse
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		return "", err
	}
	if errorResponse.Error.Message != "" {
		return "", fmt.Errorf("error: %s", errorResponse.Error.Message)
	}
	return "", fmt.Errorf("no response from model")
}

func handleOpenAITextResponse(body []byte) (string, error) {
	var response OpenAITextGenerationResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}
	var errorResponse OpenAIErrorResponse
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		return "", err
	}
	if errorResponse.Error.Message != "" {
		return "", fmt.Errorf("error: %s", errorResponse.Error.Message)
	}
	return "", fmt.Errorf("no response from model")
}

func handleClaudeTextResponse(body []byte) (string, error) {
	var response ClaudeTextGenerationResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if len(response.Content) > 0 {
		return response.Content[0].Text, nil
	}
	var errorResponse ClaudeErrorResponse
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		return "", err
	}
	if errorResponse.Error.Message != "" {
		return "", fmt.Errorf("error: %s", errorResponse.Error.Message)
	}
	return "", fmt.Errorf("no response from model")
}

func GetDefaultModel() Model {
	return Models[DefaultModelIdentifier]
}

func GetTextGenerationModels() []Model {
	var textModels []Model

	for _, m := range Models {
		if m.Generates == "text" {
			textModels = append(textModels, m)
		}
	}
	return textModels
}

func GetImageGenerationModel() Model {
	for _, m := range Models {
		if m.Generates == "image" {
			return m
		}
	}

	panic("no image model found")
}
