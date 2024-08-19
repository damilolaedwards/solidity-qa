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
		// Ensure that the user and assistant messages are in pairs
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
