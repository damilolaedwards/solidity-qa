package api

import (
	"assistant/api/handlers"
	"github.com/gorilla/mux"
)

func attachConversationRoutes(router *mux.Router, targetContracts string) {
	ch := handlers.NewConversationHandler(targetContracts)
	conversationRoute := "/conversation"
	router.HandleFunc(conversationRoute, ch.GetConversation).Methods("GET")
	router.HandleFunc(conversationRoute, ch.PromptLLM).Methods("POST")
	router.HandleFunc(conversationRoute, ch.ResetConversation).Methods("DELETE")
}
