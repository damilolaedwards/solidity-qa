package services

import (
	"assistant/llm"
	"assistant/types"
	"context"
	"sync"
)

type ConversationService struct {
	conversation []types.Message
	mu           sync.Mutex
}

func NewConversationService(targetContracts string) *ConversationService {
	return &ConversationService{
		conversation: llm.InitialPrompts(targetContracts),
	}
}

func conversationResponse(conversation []types.Message) []types.Message {
	return conversation[2:]
}

func (ch *ConversationService) GetConversation() []types.Message {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	return conversationResponse(ch.conversation)
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

func (ch *ConversationService) PromptLLM(ctx context.Context, prompt string, model string) ([]types.Message, error) {
	ch.mu.Lock()
	ch.conversation = append(ch.conversation, types.Message{
		Role:    "user",
		Content: prompt,
		Type:    "text",
		Model:   model,
	})
	ch.mu.Unlock()

	messages := ch.conversation

	if model == "claude-3-5-sonnet-20240620" {
		// Filter out messages from other models
		var filteredMessages []types.Message
		for index, message := range messages {
			if message.Model == "claude-3-5-sonnet-20240620" || index == 0 || index == 1 {
				filteredMessages = append(filteredMessages, message)
			}
		}

		messages = filteredMessages
	}

	response, err := llm.AskModel(generateApiMessages(messages), model, ctx)
	if err != nil {
		return nil, err
	}

	ch.mu.Lock()
	if model == llm.GetImageGenerationModel().Model {
		ch.conversation = append(ch.conversation, types.Message{
			Role:    "assistant",
			Content: response,
			Type:    "image",
			Model:   model,
		})
	} else {
		ch.conversation = append(ch.conversation, types.Message{
			Role:    "assistant",
			Content: response,
			Type:    "text",
			Model:   model,
		})
	}
	ch.mu.Unlock()

	return conversationResponse(ch.conversation), nil
}

func (ch *ConversationService) ResetConversation() {
	ch.mu.Lock()
	ch.conversation = ch.conversation[0:2]
	ch.mu.Unlock()
}
