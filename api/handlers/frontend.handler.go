package handlers

import (
	"assistant/api/assets/templ/components"
	"assistant/api/assets/templ/pages"
	"assistant/api/dto"
	"assistant/api/services"
	"assistant/api/utils"
	"assistant/llm"
	"assistant/types"
	"fmt"
	"log"
	"net/http"
)

var conversationService *services.ConversationService

type FrontendHandler struct {
	contracts     []types.Contract
	errorMessages []int
	isSidebarOpen bool
	selectedModel string
	projectName   string
}

type ViewProps struct {
	projectName   string
	conversation  []types.Message
	contracts     []types.Contract
	errorMessages []int
	selectedModel string
	isSidebarOpen bool
}

func NewFrontendHandler(targetContracts string, contracts []types.Contract, projectName string) (*FrontendHandler, error) {
	// Initialize the conversation service
	var err error

	conversationService, err = services.NewConversationService(targetContracts)
	if err != nil {
		return nil, err
	}

	return &FrontendHandler{
		contracts:     contracts,
		errorMessages: []int{},
		isSidebarOpen: true,
		selectedModel: llm.DefaultModelIdentifier,
		projectName:   projectName,
	}, nil
}

func (h *FrontendHandler) Get(w http.ResponseWriter, r *http.Request) {
	var props ViewProps

	props.conversation = conversationService.GetConversation()
	props.contracts = h.contracts
	props.isSidebarOpen = h.isSidebarOpen
	props.errorMessages = h.errorMessages
	props.selectedModel = h.selectedModel
	props.projectName = h.projectName

	// Display the view
	h.View(w, r, props)
}

func (h *FrontendHandler) View(w http.ResponseWriter, r *http.Request, props ViewProps) {
	pages.Home(props.projectName, props.contracts, props.conversation, props.errorMessages, props.isSidebarOpen, h.selectedModel).Render(r.Context(), w)
}

func (h *FrontendHandler) ToggleSidebar(w http.ResponseWriter, r *http.Request) {
	h.isSidebarOpen = !h.isSidebarOpen

	// Display the view
	components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
}

func (h *FrontendHandler) ChangeModel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var data dto.ChangeModelDto

	if err := utils.DecodeAndValidateFormData(r, &data); err != nil {
		ErrorResponse(w, fmt.Sprintf("Error processing form data: %v", err), http.StatusBadRequest)
		components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
		return
	}

	h.selectedModel = data.Model

	// Display the view
	components.HeaderTemplate(h.projectName, h.isSidebarOpen, h.selectedModel).Render(r.Context(), w)
}

func (h *FrontendHandler) ResetConversation(w http.ResponseWriter, r *http.Request) {
	conversationService.ResetConversation()

	components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
}

func (h *FrontendHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var data dto.GenerateReportDto

	if err := utils.DecodeAndValidateFormData(r, &data); err != nil {
		ErrorResponse(w, fmt.Sprintf("Error processing form data: %v", err), http.StatusBadRequest)
		components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
		return
	}

	// Get the request context to handle request cancellations
	ctx := r.Context()

	responseChan := make(chan []types.Message)
	errorChan := make(chan error)

	go func() {
		response, err := conversationService.GenerateReport(ctx, data, h.selectedModel)
		if err != nil {
			errorChan <- err
			return
		}

		responseChan <- response
	}()

	select {
	case err := <-errorChan:
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)

	case <-ctx.Done():
		log.Println("Client canceled the request")
		w.WriteHeader(http.StatusRequestTimeout)
		components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
		return

	case response := <-responseChan:
		// Display the view
		components.MainContent(h.contracts, response, h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
	}
}

func (h *FrontendHandler) PromptLLM(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var data dto.PromptLLMDto

	if err := utils.DecodeAndValidateFormData(r, &data); err != nil {
		ErrorResponse(w, fmt.Sprintf("Error processing form data: %v", err), http.StatusBadRequest)
		components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
		return
	}

	// Get the request context to handle request cancellations
	ctx := r.Context()

	responseChan := make(chan []types.Message)
	errorChan := make(chan error)

	go func() {
		var model string
		if data.GenerateImage == "on" {
			model = llm.GetImageGenerationModel().Identifier
		} else {
			model = h.selectedModel
		}

		response, err := conversationService.PromptLLM(ctx, data.Message, model)
		if err != nil {
			errorChan <- err
			return
		}

		responseChan <- response
	}()

	select {
	case err := <-errorChan:
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.errorMessages = append(h.errorMessages, len(conversationService.GetConversation())-1)
		components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)

	case <-ctx.Done():
		log.Println("Client canceled the request")
		w.WriteHeader(http.StatusRequestTimeout)
		h.errorMessages = append(h.errorMessages, len(conversationService.GetConversation())-1)
		components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
		return

	case response := <-responseChan:
		// Display the view
		components.MainContent(h.contracts, response, h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
	}
}
