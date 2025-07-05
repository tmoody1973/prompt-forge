package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"promptforge/internal/config"
	"promptforge/internal/database"
	"promptforge/internal/models"
	"promptforge/internal/services"
)

type Handlers struct {
	db             *database.Database
	aiService      *services.UnifiedAIService
	promptAnalyzer *services.PromptAnalyzer
	evalGenerator  *services.EvalGenerator
}

func NewHandlers(db *database.Database, aiService *services.UnifiedAIService) *Handlers {
	promptAnalyzer := services.NewPromptAnalyzer(aiService)
	evalGenerator := services.NewEvalGenerator(aiService)

	return &Handlers{
		db:             db,
		aiService:      aiService,
		promptAnalyzer: promptAnalyzer,
		evalGenerator:  evalGenerator,
	}
}

func (h *Handlers) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "PromptForge API",
	})
}

func (h *Handlers) GetProviders(c echo.Context) error {
	providers := map[string]interface{}{
		"default": config.AppConfig.DefaultProvider,
		"available": []string{
			string(config.ProviderOpenAI),
			string(config.ProviderAzureOpenAI),
			string(config.ProviderAnthropic),
		},
		"configured": map[string]bool{
			string(config.ProviderOpenAI):      config.AppConfig.OpenAI.APIKey != "",
			string(config.ProviderAzureOpenAI): config.AppConfig.AzureOpenAI.APIKey != "",
			string(config.ProviderAnthropic):   config.AppConfig.Anthropic.APIKey != "",
		},
	}

	return c.JSON(http.StatusOK, providers)
}

func (h *Handlers) CritiquePrompt(c echo.Context) error {
	var req models.CritiqueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	model := req.Model
	if model == "" {
		model = models.DefaultGPTModel // Default model
	}

	// Use the enhanced prompt analyzer
	response, err := h.promptAnalyzer.AnalyzePrompt(req.Prompt, model)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to get comprehensive analysis: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

func (h *Handlers) ExecutePrompt(c echo.Context) error {
	var req models.ExecuteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	messages := []models.Message{
		{Role: "user", Content: req.Prompt},
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.7 // Default temperature
	}

	model := req.Model
	if model == "" {
		model = models.DefaultGPTModel // Default model
	}

	response, err := h.aiService.CallWithDefaultProvider(messages, temperature, req.MaxTokens, model)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to execute prompt: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

func (h *Handlers) MultiModelExecute(c echo.Context) error {
	var req models.MultiModelExecuteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.MultiModelExecuteResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	if len(req.Models) == 0 {
		return c.JSON(http.StatusBadRequest, models.MultiModelExecuteResponse{
			Success: false,
			Error:   "At least one model must be specified",
		})
	}

	messages := []models.Message{
		{Role: "user", Content: req.Prompt},
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.7 // Default temperature
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 1000 // Default max tokens
	}

	var results []models.ModelExecutionResult

	// Execute prompt against each model
	for _, model := range req.Models {
		startTime := time.Now()

		response, err := h.aiService.CallWithDefaultProvider(messages, temperature, maxTokens, model)

		executionTime := time.Since(startTime).Milliseconds()

		result := models.ModelExecutionResult{
			Model:         model,
			ExecutionTime: executionTime,
		}

		if err != nil {
			result.Success = false
			result.Error = err.Error()
		} else {
			result.Success = true
			result.Response = response
			// Note: Token usage would need to be extracted from the AI service response
			// This is a placeholder for future enhancement
		}

		results = append(results, result)
	}

	return c.JSON(http.StatusOK, models.MultiModelExecuteResponse{
		Success: true,
		Data:    results,
	})
}

func (h *Handlers) GetHistory(c echo.Context) error {
	history, err := h.db.GetHistory()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.HistoryResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to retrieve history: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.HistoryResponse{
		Success: true,
		Data:    history,
	})
}

func (h *Handlers) SaveHistory(c echo.Context) error {
	var req models.SaveHistoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	if err := h.db.SaveHistory(req); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to save history: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    "History saved successfully",
	})
}

func (h *Handlers) ClearHistory(c echo.Context) error {
	if err := h.db.ClearHistory(); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to clear history: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    "History cleared successfully",
	})
}

func (h *Handlers) PromptEngineer(c echo.Context) error {
	var req models.PromptEngineerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	model := req.Model
	if model == "" {
		model = models.DefaultO3Model // Default to o3 for prompt engineering
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.7 // Default temperature
	}

	response, err := h.aiService.CallWithDefaultProvider(req.Messages, temperature, 2000, model)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to get prompt engineering response: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

func (h *Handlers) DualCritiquePrompt(c echo.Context) error {
	var req models.CritiqueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.DualAnalysisResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	model := req.Model
	if model == "" {
		model = models.DefaultGPTModel // Default model
	}

	// Use the dual prompt analyzer
	response, err := h.promptAnalyzer.DualAnalyzePrompt(req.Prompt, model)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.DualAnalysisResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to get dual analysis: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.DualAnalysisResponse{
		Success: true,
		Data:    response,
	})
}

// Conversation management handlers
func (h *Handlers) GetConversations(c echo.Context) error {
	conversations, err := h.db.GetConversations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ConversationResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to retrieve conversations: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.ConversationResponse{
		Success: true,
		Data:    conversations,
	})
}

func (h *Handlers) GetConversation(c echo.Context) error {
	conversationID := c.Param("id")
	if conversationID == "" {
		return c.JSON(http.StatusBadRequest, models.ConversationDetailResponse{
			Success: false,
			Error:   "Conversation ID is required",
		})
	}

	conversation, err := h.db.GetConversation(conversationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ConversationDetailResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to retrieve conversation: %v", err),
		})
	}

	if conversation == nil {
		return c.JSON(http.StatusNotFound, models.ConversationDetailResponse{
			Success: false,
			Error:   "Conversation not found",
		})
	}

	return c.JSON(http.StatusOK, models.ConversationDetailResponse{
		Success: true,
		Data:    conversation,
	})
}

func (h *Handlers) SaveConversation(c echo.Context) error {
	var req models.SaveConversationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	if err := h.db.SaveConversation(req); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to save conversation: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    "Conversation saved successfully",
	})
}

func (h *Handlers) DeleteConversation(c echo.Context) error {
	conversationID := c.Param("id")
	if conversationID == "" {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Conversation ID is required",
		})
	}

	if err := h.db.DeleteConversation(conversationID); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to delete conversation: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    "Conversation deleted successfully",
	})
}

// Prompt Library handlers
func (h *Handlers) GetSavedPrompts(c echo.Context) error {
	prompts, err := h.db.GetSavedPrompts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.PromptLibraryResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to retrieve saved prompts: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.PromptLibraryResponse{
		Success: true,
		Data:    prompts,
	})
}

func (h *Handlers) GetSavedPrompt(c echo.Context) error {
	promptID := c.Param("id")
	if promptID == "" {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Prompt ID is required",
		})
	}

	// Convert string to int64
	var id int64
	if _, err := fmt.Sscanf(promptID, "%d", &id); err != nil {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Invalid prompt ID format",
		})
	}

	prompt, err := h.db.GetSavedPrompt(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.PromptResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to retrieve prompt: %v", err),
		})
	}

	if prompt == nil {
		return c.JSON(http.StatusNotFound, models.PromptResponse{
			Success: false,
			Error:   "Prompt not found",
		})
	}

	return c.JSON(http.StatusOK, models.PromptResponse{
		Success: true,
		Data:    prompt,
	})
}

func (h *Handlers) SavePrompt(c echo.Context) error {
	var req models.SavePromptRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	// Validate required fields
	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Title is required",
		})
	}

	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Content is required",
		})
	}

	prompt, err := h.db.SavePrompt(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.PromptResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to save prompt: %v", err),
		})
	}

	return c.JSON(http.StatusCreated, models.PromptResponse{
		Success: true,
		Data:    prompt,
	})
}

func (h *Handlers) UpdatePrompt(c echo.Context) error {
	promptID := c.Param("id")
	if promptID == "" {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Prompt ID is required",
		})
	}

	// Convert string to int64
	var id int64
	if _, err := fmt.Sscanf(promptID, "%d", &id); err != nil {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Invalid prompt ID format",
		})
	}

	var req models.UpdatePromptRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	// Set the ID from the URL parameter
	req.ID = id

	// Validate required fields
	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Title is required",
		})
	}

	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Content is required",
		})
	}

	prompt, err := h.db.UpdatePrompt(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.PromptResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to update prompt: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.PromptResponse{
		Success: true,
		Data:    prompt,
	})
}

func (h *Handlers) DeletePrompt(c echo.Context) error {
	promptID := c.Param("id")
	if promptID == "" {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Prompt ID is required",
		})
	}

	// Convert string to int64
	var id int64
	if _, err := fmt.Sscanf(promptID, "%d", &id); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid prompt ID format",
		})
	}

	if err := h.db.DeletePrompt(id); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to delete prompt: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    "Prompt deleted successfully",
	})
}

func (h *Handlers) UsePrompt(c echo.Context) error {
	promptID := c.Param("id")
	if promptID == "" {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Prompt ID is required",
		})
	}

	// Convert string to int64
	var id int64
	if _, err := fmt.Sscanf(promptID, "%d", &id); err != nil {
		return c.JSON(http.StatusBadRequest, models.PromptResponse{
			Success: false,
			Error:   "Invalid prompt ID format",
		})
	}

	// Increment usage count
	if err := h.db.IncrementPromptUsage(id); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to increment prompt usage: %v", err),
		})
	}

	// Get the prompt to return
	prompt, err := h.db.GetSavedPrompt(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.PromptResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to retrieve prompt: %v", err),
		})
	}

	if prompt == nil {
		return c.JSON(http.StatusNotFound, models.PromptResponse{
			Success: false,
			Error:   "Prompt not found",
		})
	}

	return c.JSON(http.StatusOK, models.PromptResponse{
		Success: true,
		Data:    prompt,
	})
}

func (h *Handlers) GenerateEval(c echo.Context) error {
	var req models.EvalGenerateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.EvalResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	// Validate required fields
	if req.Prompt == "" {
		return c.JSON(http.StatusBadRequest, models.EvalResponse{
			Success: false,
			Error:   "Prompt is required",
		})
	}

	if len(req.EvalTypes) == 0 {
		return c.JSON(http.StatusBadRequest, models.EvalResponse{
			Success: false,
			Error:   "At least one evaluation type is required",
		})
	}

	if req.SampleSize <= 0 {
		req.SampleSize = 10 // Default sample size
	}

	if req.Model == "" {
		req.Model = "gpt-4.1" // Default model
	}

	if req.Difficulty == "" {
		req.Difficulty = "mixed" // Default difficulty
	}

	// Generate evaluation suite
	evalData, err := h.evalGenerator.GenerateEvaluationSuite(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.EvalResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to generate evaluation suite: %v", err),
		})
	}

	return c.JSON(http.StatusOK, models.EvalResponse{
		Success: true,
		Data:    evalData,
	})
}
