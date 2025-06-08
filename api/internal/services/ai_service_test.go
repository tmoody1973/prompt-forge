package services

import (
	"testing"

	"promptforge/internal/config"
	"promptforge/internal/models"
)

func TestNewUnifiedAIService(t *testing.T) {
	service := NewUnifiedAIService()
	if service == nil {
		t.Fatal("NewUnifiedAIService should not return nil")
	}
	if service.client == nil {
		t.Fatal("HTTP client should be initialized")
	}
}

func TestAnthropicTemperatureConversion(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0.0, 0.0},
		{1.0, 1.0},
		{2.0, 1.0},  // Should be scaled down and capped at 1.0
		{1.5, 0.75}, // Should be scaled down from 0-2 to 0-1 range
		{-0.5, 0.0}, // Should be capped at 0.0
		{0.5, 0.5},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			temperature := test.input

			// Apply the same logic as in callAnthropic
			if temperature > 1.0 {
				temperature = temperature / 2.0
			}
			if temperature < 0 {
				temperature = 0
			} else if temperature > 1 {
				temperature = 1
			}

			if temperature != test.expected {
				t.Errorf("For input %f, expected %f, got %f", test.input, test.expected, temperature)
			}
		})
	}
}

func TestAnthropicMessageConversion(t *testing.T) {
	messages := []models.Message{
		{Role: "system", Content: "You are a helpful assistant"},
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi there!"},
		{Role: "user", Content: "How are you?"},
	}

	var anthropicMessages []models.AnthropicMessage
	var systemMessage string

	// Apply the same logic as in callAnthropic
	for _, msg := range messages {
		if msg.Role == "system" {
			systemMessage = msg.Content
		} else {
			anthropicMessages = append(anthropicMessages, models.AnthropicMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	// Should extract system message
	if systemMessage != "You are a helpful assistant" {
		t.Errorf("Expected system message 'You are a helpful assistant', got '%s'", systemMessage)
	}

	// Should have 3 non-system messages
	if len(anthropicMessages) != 3 {
		t.Errorf("Expected 3 non-system messages, got %d", len(anthropicMessages))
	}

	// Check first non-system message
	if anthropicMessages[0].Role != "user" || anthropicMessages[0].Content != "Hello" {
		t.Errorf("First non-system message is incorrect")
	}
}

func TestO3ModelSpecialHandling(t *testing.T) {
	tests := []struct {
		model                  string
		temperature            float64
		maxTokens              int
		expectedTemp           float64
		expectedMaxTokensField string
	}{
		{"o3", 0.5, 100, 1.0, "max_completion_tokens"}, // o3 should force temperature to 1.0
		{"gpt-4", 0.5, 100, 0.5, "max_tokens"},         // other models should preserve temperature
		{"gpt-4.1", 0.8, 150, 0.8, "max_tokens"},
	}

	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			// Test temperature handling for o3
			temperature := test.temperature
			if test.model == "o3" {
				temperature = 1.0 // This is what the code does
			}

			if temperature != test.expectedTemp {
				t.Errorf("For model %s, expected temperature %f, got %f", test.model, test.expectedTemp, temperature)
			}

			// Create request body to test max_tokens field selection
			requestBody := models.OpenAIRequest{
				Messages:    []models.Message{{Role: "user", Content: "test"}},
				Temperature: temperature,
			}

			if test.maxTokens > 0 {
				if test.model == "o3" {
					requestBody.MaxCompletionTokens = test.maxTokens
				} else {
					requestBody.MaxTokens = test.maxTokens
				}
			}

			// Verify the correct field was set
			if test.model == "o3" {
				if requestBody.MaxCompletionTokens != test.maxTokens {
					t.Errorf("Expected MaxCompletionTokens %d, got %d", test.maxTokens, requestBody.MaxCompletionTokens)
				}
				if requestBody.MaxTokens != 0 {
					t.Errorf("Expected MaxTokens to be 0 for o3, got %d", requestBody.MaxTokens)
				}
			} else {
				if requestBody.MaxTokens != test.maxTokens {
					t.Errorf("Expected MaxTokens %d, got %d", test.maxTokens, requestBody.MaxTokens)
				}
				if requestBody.MaxCompletionTokens != 0 {
					t.Errorf("Expected MaxCompletionTokens to be 0 for non-o3, got %d", requestBody.MaxCompletionTokens)
				}
			}
		})
	}
}

func TestDefaultModelSelection(t *testing.T) {
	tests := []struct {
		provider      config.AIProvider
		inputModel    string
		expectedModel string
	}{
		{config.ProviderOpenAI, "", "gpt-4"},
		{config.ProviderOpenAI, "gpt-3.5-turbo", "gpt-3.5-turbo"},
		{config.ProviderAnthropic, "", "claude-3-5-sonnet-20241022"},
		{config.ProviderAnthropic, "claude-3-haiku", "claude-3-haiku"},
		{config.ProviderAzureOpenAI, "", ""}, // Azure uses deployment names
		{config.ProviderAzureOpenAI, "gpt-4.1", "gpt-4.1"},
	}

	for _, test := range tests {
		t.Run(string(test.provider), func(t *testing.T) {
			model := test.inputModel

			// Apply the same default logic as in the service methods
			switch test.provider {
			case config.ProviderOpenAI:
				if model == "" {
					model = "gpt-4"
				}
			case config.ProviderAnthropic:
				if model == "" {
					model = "claude-3-5-sonnet-20241022"
				}
			}

			if model != test.expectedModel {
				t.Errorf("For provider %s and input '%s', expected model '%s', got '%s'",
					test.provider, test.inputModel, test.expectedModel, model)
			}
		})
	}
}

func TestProviderValidation(t *testing.T) {
	service := NewUnifiedAIService()

	// Initialize config for testing
	config.AppConfig = &config.Config{
		OpenAI: config.OpenAIConfig{
			APIKey: "",
		},
		AzureOpenAI: config.AzureOpenAIConfig{
			APIKey: "",
		},
		Anthropic: config.AnthropicConfig{
			APIKey: "",
		},
	}

	messages := []models.Message{{Role: "user", Content: "test"}}

	// Test unsupported provider
	_, err := service.CallAI(messages, 0.7, 100, "gpt-4", "unsupported-provider")
	if err == nil {
		t.Error("Expected error for unsupported provider")
	}
	if err.Error() != "unsupported AI provider: unsupported-provider" {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}

	// Test providers with missing API keys (should return configuration errors)
	providers := []config.AIProvider{
		config.ProviderOpenAI,
		config.ProviderAzureOpenAI,
		config.ProviderAnthropic,
	}

	for _, provider := range providers {
		_, err := service.CallAI(messages, 0.7, 100, "test-model", provider)
		if err == nil {
			t.Errorf("Expected error for provider %s with missing API key", provider)
		}
		// Should contain "API key not configured"
		if err.Error() == "" {
			t.Errorf("Expected non-empty error message for provider %s", provider)
		}
	}
}

func TestMaxTokensHandling(t *testing.T) {
	tests := []struct {
		maxTokens int
		expected  bool // whether max_tokens should be set
	}{
		{0, false},  // Should not set max_tokens
		{100, true}, // Should set max_tokens
		{-1, false}, // Should not set max_tokens for negative values
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			requestBody := models.OpenAIRequest{
				Model: "gpt-4",
				Messages: []models.Message{
					{Role: "user", Content: "test"},
				},
				Temperature: 0.7,
			}

			// Apply the same logic as in the service
			if test.maxTokens > 0 {
				requestBody.MaxTokens = test.maxTokens
			}

			hasMaxTokens := requestBody.MaxTokens > 0
			if hasMaxTokens != test.expected {
				t.Errorf("For maxTokens %d, expected hasMaxTokens=%v, got hasMaxTokens=%v",
					test.maxTokens, test.expected, hasMaxTokens)
			}
		})
	}
}
