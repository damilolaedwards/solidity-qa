package dto

type ChangeModelDto struct {
	Model string `json:"model" form:"model" validate:"required"`
}

type GenerateReportDto struct {
	ReportType        string `json:"reportType" form:"reportType" validate:"required"`
	AdditionalMessage string `json:"additionalMessage" form:"additionalMessage"`
}

type PromptLLMDto struct {
	Message       string `json:"message" form:"message" validate:"required"`
	GenerateImage string `json:"generateImage" form:"message"`
}
