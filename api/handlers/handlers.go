package handlers

import (
	"assistant/api/dto"
	"assistant/api/utils"
	"assistant/llm"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type ConversationHandler struct {
	conversation []llm.Message
	mu           sync.Mutex
}

func NewConversationHandler(targetContracts string) *ConversationHandler {
	return &ConversationHandler{
		conversation: []llm.Message{llm.InitialPrompt(targetContracts)},
	}
}

func (ch *ConversationHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	response := map[string][]llm.Message{"conversation": ch.conversation[1:]}
	writeJSONResponse(w, http.StatusOK, response)
}

func (ch *ConversationHandler) PromptLLM(w http.ResponseWriter, r *http.Request) {
	var data dto.PromptLLMDto

	if err := utils.DecodeAndValidateRequestBody(r, &data); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	// Get the request context to handle request cancellations
	ctx := r.Context()

	responseChan := make(chan string)
	go func() {
		ch.mu.Lock()
		ch.conversation = append(ch.conversation, llm.Message{
			Role:    "user",
			Content: data.Message,
		})
		ch.mu.Unlock()

		response, err := llm.AskModel(ch.conversation, ctx)
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, errorResponse(err.Error()))
			return
		}

		responseChan <- response
	}()

	select {
	case <-ctx.Done():
		// The client canceled the request
		log.Println("Client canceled the request")
		w.WriteHeader(http.StatusRequestTimeout)
		return
	case response := <-responseChan:
		ch.mu.Lock()
		ch.conversation = append(ch.conversation, llm.Message{
			Role:    "system",
			Content: response,
		})
		ch.mu.Unlock()

		writeJSONResponse(w, http.StatusOK, map[string]string{"response": response})
	}
}

func (ch *ConversationHandler) ResetConversation(w http.ResponseWriter, r *http.Request) {
	ch.mu.Lock()
	ch.conversation = ch.conversation[0:1] // Keep the first prompt
	ch.mu.Unlock()

	writeJSONResponse(w, http.StatusOK, messageResponse("Conversation reset successfully"))
}

func messageResponse(message string) map[string]string {
	return map[string]string{"message": message}
}

func errorResponse(errMessage string) map[string]string {
	return map[string]string{"error": errMessage}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
