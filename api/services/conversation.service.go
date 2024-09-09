package services

import (
	"assistant/api/client"
	"assistant/api/dto"
	"assistant/llm"
	"assistant/types"
	"context"
	"fmt"
	"sync"
)

var ReportTypes = map[string]string{
	"Access Control":       "https://raw.githubusercontent.com/brainycodelab/report-datasets/master/access_controls.md",
	"Auditing and Logging": "https://raw.githubusercontent.com/brainycodelab/report-datasets/master/auditing_and_logging.md",
	"Authentication":       "https://raw.githubusercontent.com/brainycodelab/report-datasets/master/authentication.md",
	"Configuration":        "https://raw.githubusercontent.com/brainycodelab/report-datasets/master/configuration.md",
	"Data Validation":      "https://raw.githubusercontent.com/brainycodelab/report-datasets/master/data_validation.md",
	"Denial of Service":    "https://raw.githubusercontent.com/brainycodelab/report-datasets/master/denial_of_service.md",
	"Patching":             "https://raw.githubusercontent.com/brainycodelab/report-datasets/master/patching.md",
	"Testing":              "https://raw.githubusercontent.com/brainycodelab/report-datasets/master/testing.md",
	"Timing":               "https://raw.githubusercontent.com/brainycodelab/report-datasets/master/timing.md",
	"Undefined Behaviour":  "https://raw.githubusercontent.com/brainycodelab/report-datasets/master/undefined_behaviour.md",
}

type ConversationService struct {
	conversation []types.Message
	mu           sync.Mutex
}

func NewConversationService(targetContracts string) (*ConversationService, error) {
	// Ensure that initial prompts don't surpass the maximum number of tokens
	initPrompts := llm.InitialPrompts(targetContracts)

	numTokens, err := llm.CalculateNumTokens(generateApiMessages(initPrompts))
	if err != nil {
		return nil, err
	}

	if numTokens > llm.GetDefaultModel().MaxTokenLen {
		return nil, fmt.Errorf("target contracts exceed maximum token length")
	}

	return &ConversationService{
		conversation: initPrompts,
	}, nil
}

func conversationResponse(conversation []types.Message) []types.Message {
	return conversation[2:]
}

func (ch *ConversationService) GetConversation() []types.Message {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	return conversationResponse(ch.conversation)
}

func (ch *ConversationService) PromptLLM(ctx context.Context, prompt string, model string) ([]types.Message, error) {
	ch.mu.Lock()
	ch.conversation = append(ch.conversation, types.Message{
		Role:    "user",
		Content: prompt,
		Type:    "text",
		Model:   llm.Models[model].Name,
	})
	ch.mu.Unlock()

	messages := ch.filterMessages(ch.conversation)

	response, err := llm.AskModel(generateApiMessages(messages), model, ctx)
	if err != nil {
		return nil, err
	}

	ch.mu.Lock()
	if model == llm.GetImageGenerationModel().Identifier {
		ch.conversation = append(ch.conversation, types.Message{
			Role:    "assistant",
			Content: response,
			Type:    "image",
			Model:   llm.Models[model].Name,
		})
	} else {
		ch.conversation = append(ch.conversation, types.Message{
			Role:    "assistant",
			Content: response,
			Type:    "text",
			Model:   llm.Models[model].Name,
		})
	}
	ch.mu.Unlock()

	return conversationResponse(ch.conversation), nil
}

func (ch *ConversationService) GenerateReport(ctx context.Context, data dto.GenerateReportDto, model string) ([]types.Message, error) {
	sampleFileUrl, ok := ReportTypes[data.ReportType]
	if !ok {
		return nil, fmt.Errorf("invalid report type")
	}

	reportSample, err := client.FetchFileContent(sampleFileUrl, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch file content: %v", err)
	}

	prompt := llm.GenerateReportPrompt(data.ReportType, reportSample, data.AdditionalMessage)

	ch.mu.Lock()
	ch.conversation = append(ch.conversation, types.Message{
		Role:    "user",
		Content: prompt,
		Type:    "text",
		Model:   llm.Models[model].Name,
		Hidden:  true,
	})
	ch.mu.Unlock()

	messages := ch.filterMessages(ch.conversation)

	response, err := llm.AskModel(generateApiMessages(messages), model, ctx)
	if err != nil {
		return nil, err
	}

	ch.mu.Lock()
	ch.conversation = append(ch.conversation, types.Message{
		Role:    "assistant",
		Content: response,
		Type:    "text",
		Model:   llm.Models[model].Name,
	})
	ch.mu.Unlock()

	return conversationResponse(ch.conversation), nil
}

func (ch *ConversationService) ResetConversation() {
	ch.mu.Lock()
	ch.conversation = ch.conversation[0:2]
	ch.mu.Unlock()
}

func GetReportTypes() []string {
	reportTypes := make([]string, 0, len(ReportTypes))

	for key := range ReportTypes {
		reportTypes = append(reportTypes, key)
	}

	return reportTypes
}

// filterMessages ensures that the user and assistant messages are in pairs,
// and non-user messages are always included.
func (ch *ConversationService) filterMessages(messages []types.Message) []types.Message {
	var filteredMessages []types.Message

	for i := 0; i < len(messages); i++ {
		if messages[i].Role == "user" {
			if i+1 < len(messages) && messages[i+1].Role == "assistant" {
				// User message followed by assistant message
				filteredMessages = append(filteredMessages, messages[i], messages[i+1])
				i++ // Skip the next message as we've already added it
			} else if i == len(messages)-1 {
				// User message is the last message in the array
				filteredMessages = append(filteredMessages, messages[i])
			}
			// If user message is not last and not followed by assistant, it's filtered out
		} else {
			// Non-user messages are always included
			filteredMessages = append(filteredMessages, messages[i])
		}
	}

	return filteredMessages
}

func generateApiMessages(messages []types.Message) []llm.ApiMessage {
	var apiMessages []llm.ApiMessage
	for _, message := range messages {
		apiMessages = append(apiMessages, llm.ApiMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return apiMessages
}
