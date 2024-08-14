package dto

type PromptLLMDto struct {
	Message       string `json:"message" validate:"required"`
	GenerateImage bool   `json:"generateImage" validate:"boolean"`
}
