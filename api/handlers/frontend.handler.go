package handlers

import (
	"assistant/api/assets/templ/components"
	"assistant/api/assets/templ/pages"
	"assistant/api/dto"
	"assistant/api/services"
	"assistant/api/utils"
	"assistant/llm"
	"assistant/types"
	"log"
	"net/http"
)

var conversationService *services.ConversationService

type FrontendHandler struct {
	contracts     []types.Contract
	errorMessages []int
	isSidebarOpen bool
	selectedModel string
}

type ViewProps struct {
	conversation  []types.Message
	contracts     []types.Contract
	textModels    []llm.Model
	errorMessages []int
	selectedModel string
	isSidebarOpen bool
}

func NewFrontendHandler(targetContracts string, contracts []types.Contract) *FrontendHandler {
	// Initialize the conversation service
	conversationService = services.NewConversationService(targetContracts)

	return &FrontendHandler{
		contracts:     contracts,
		errorMessages: []int{},
		isSidebarOpen: true,
		selectedModel: llm.DefaultModel,
	}
}

func (h *FrontendHandler) Get(w http.ResponseWriter, r *http.Request) {
	var props ViewProps

	props.conversation = conversationService.GetConversation()
	props.contracts = h.contracts
	props.isSidebarOpen = h.isSidebarOpen
	props.errorMessages = h.errorMessages
	props.textModels = llm.GetTextGenerationModels()
	props.selectedModel = h.selectedModel

	// Display the view
	h.View(w, r, props)
}

func (h *FrontendHandler) View(w http.ResponseWriter, r *http.Request, props ViewProps) {
	pages.Home(props.contracts, props.conversation, props.errorMessages, props.isSidebarOpen, props.textModels, h.selectedModel).Render(r.Context(), w)
}

func (h *FrontendHandler) ToggleSidebar(w http.ResponseWriter, r *http.Request) {
	h.isSidebarOpen = !h.isSidebarOpen

	// Display the view
	components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
}

func (h *FrontendHandler) ChangeModel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	h.selectedModel = r.FormValue("model")

	// Display the view
	components.HeaderTemplate(h.isSidebarOpen, llm.GetTextGenerationModels(), h.selectedModel).Render(r.Context(), w)
}

func (h *FrontendHandler) ResetConversation(w http.ResponseWriter, r *http.Request) {
	conversationService.ResetConversation()

	components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
}

func (h *FrontendHandler) PromptLLM(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var data dto.PromptLLMDto

	data.Message = r.FormValue("message")
	data.GenerateImage = r.FormValue("generateImage") == "on"

	if err := utils.ValidateData(&data); err != nil {
		log.Println("Validation error:", err)
		w.WriteHeader(http.StatusBadRequest)
		components.MainContent(h.contracts, conversationService.GetConversation(), h.errorMessages, h.isSidebarOpen).Render(r.Context(), w)
		return
	}

	// Get the request context to handle request cancellations
	ctx := r.Context()

	responseChan := make(chan []types.Message)
	errorChan := make(chan error)

	go func() {
		var model string
		if data.GenerateImage {
			model = llm.GetImageGenerationModel().Model
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
		// The client canceled the request
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
