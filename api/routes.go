package api

import (
	"assistant/api/handlers"
	"assistant/types"
	"fmt"
	"github.com/gorilla/mux"
)

func attachFrontendRoutes(router *mux.Router, contracts []types.Contract, targetContracts string, projectName string) error {
	ch, err := handlers.NewFrontendHandler(targetContracts, contracts, projectName)
	if err != nil {
		return err
	}

	frontendRoute := "/"

	router.HandleFunc(frontendRoute, ch.Get).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%stoggle-sidebar", frontendRoute), ch.ToggleSidebar).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%schange-model", frontendRoute), ch.ChangeModel).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%sreset", frontendRoute), ch.ResetConversation).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%sprompt", frontendRoute), ch.PromptLLM).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%sgenerate-report", frontendRoute), ch.GenerateReport).Methods("POST")

	return nil
}
