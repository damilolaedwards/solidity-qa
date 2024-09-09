package dto

type ChangeModelDto struct {
	Model string `json:"model" validate:"required"`
}

type GenerateReportDto struct {
	ReportType        string `json:"reportType" validate:"required"`
	AdditionalMessage string `json:"additionalMessage"`
}

type PromptLLMDto struct {
	Message       string `json:"message" validate:"required"`
	GenerateImage string `json:"generateImage"`
}
