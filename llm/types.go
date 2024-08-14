package llm

type TextGenerationRequest struct {
	Model     string       `json:"model"`
	Messages  []ApiMessage `json:"messages"`
	MaxTokens int          `json:"max_tokens"`
}

type ImageGenerationRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ApiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAITextGenerationResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message      ApiMessage `json:"message"`
		FinishReason string     `json:"finish_reason"`
		Index        int        `json:"index"`
	} `json:"choices"`
}

type ClaudeTextGenerationResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
	Id           string      `json:"id"`
	Model        string      `json:"model"`
	Role         string      `json:"role"`
	StopReason   string      `json:"stop_reason"`
	StopSequence interface{} `json:"stop_sequence"`
	Type         string      `json:"type"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type ImageGenerationResponse struct {
	Created int    `json:"created"`
	Data    []Data `json:"data"`
}

type Data struct {
	URL string `json:"url"`
}

type OpenAIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

type ClaudeErrorResponse struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}
