package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMessageJSONSerialization(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "Test message",
	}

	// Test marshaling
	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal Message: %v", err)
	}

	expected := `{"role":"user","content":"Test message"}`
	if string(jsonData) != expected {
		t.Errorf("Expected JSON: %s, got: %s", expected, string(jsonData))
	}

	// Test unmarshaling
	var unmarshaled Message
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Message: %v", err)
	}

	if unmarshaled.Role != msg.Role || unmarshaled.Content != msg.Content {
		t.Errorf("Unmarshaled message doesn't match original")
	}
}

func TestOpenAIRequestSerialization(t *testing.T) {
	req := OpenAIRequest{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		Temperature:         0.7,
		MaxTokens:           100,
		MaxCompletionTokens: 150,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal OpenAIRequest: %v", err)
	}

	var unmarshaled OpenAIRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal OpenAIRequest: %v", err)
	}

	if unmarshaled.Model != req.Model {
		t.Errorf("Expected model %s, got %s", req.Model, unmarshaled.Model)
	}
	if unmarshaled.Temperature != req.Temperature {
		t.Errorf("Expected temperature %f, got %f", req.Temperature, unmarshaled.Temperature)
	}
	if len(unmarshaled.Messages) != len(req.Messages) {
		t.Errorf("Expected %d messages, got %d", len(req.Messages), len(unmarshaled.Messages))
	}
}

func TestCritiqueRequestValidation(t *testing.T) {
	tests := []struct {
		name     string
		request  CritiqueRequest
		hasError bool
	}{
		{
			name: "Valid request",
			request: CritiqueRequest{
				Prompt: "Test prompt",
				Model:  "gpt-4",
			},
			hasError: false,
		},
		{
			name: "Empty prompt should still be valid (let API handle validation)",
			request: CritiqueRequest{
				Prompt: "",
				Model:  "gpt-4",
			},
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, err := json.Marshal(test.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			var unmarshaled CritiqueRequest
			err = json.Unmarshal(jsonData, &unmarshaled)
			if (err != nil) != test.hasError {
				t.Errorf("Expected error: %v, got error: %v", test.hasError, err != nil)
			}
		})
	}
}

func TestHistoryItemTimestamp(t *testing.T) {
	now := time.Now()
	item := HistoryItem{
		ID:          1,
		Timestamp:   now,
		Prompt:      "Test prompt",
		Model:       "gpt-4",
		Temperature: 0.7,
		MaxTokens:   100,
		Success:     true,
		Response:    "Test response",
	}

	jsonData, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal HistoryItem: %v", err)
	}

	var unmarshaled HistoryItem
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal HistoryItem: %v", err)
	}

	// Compare timestamps with some tolerance for JSON serialization precision
	if unmarshaled.Timestamp.Unix() != now.Unix() {
		t.Errorf("Timestamp mismatch: expected %v, got %v", now.Unix(), unmarshaled.Timestamp.Unix())
	}
}

func TestAPIResponseStructure(t *testing.T) {
	// Test success response
	successResp := APIResponse{
		Success: true,
		Data:    "Operation successful",
	}

	jsonData, err := json.Marshal(successResp)
	if err != nil {
		t.Fatalf("Failed to marshal success response: %v", err)
	}

	var unmarshaled APIResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal success response: %v", err)
	}

	if !unmarshaled.Success {
		t.Error("Expected success to be true")
	}
	if unmarshaled.Data != "Operation successful" {
		t.Errorf("Expected data 'Operation successful', got '%s'", unmarshaled.Data)
	}

	// Test error response
	errorResp := APIResponse{
		Success: false,
		Error:   "Something went wrong",
	}

	jsonData, err = json.Marshal(errorResp)
	if err != nil {
		t.Fatalf("Failed to marshal error response: %v", err)
	}

	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if unmarshaled.Success {
		t.Error("Expected success to be false")
	}
	if unmarshaled.Error != "Something went wrong" {
		t.Errorf("Expected error 'Something went wrong', got '%s'", unmarshaled.Error)
	}
}

func TestConversationStructure(t *testing.T) {
	now := time.Now()
	conv := Conversation{
		ID:        "conv-123",
		Title:     "Test Conversation",
		CreatedAt: now,
		UpdatedAt: now,
		Messages: []ConversationMessage{
			{
				ID:             1,
				ConversationID: "conv-123",
				Role:           "user",
				Content:        "Hello",
				Timestamp:      now,
			},
		},
	}

	jsonData, err := json.Marshal(conv)
	if err != nil {
		t.Fatalf("Failed to marshal Conversation: %v", err)
	}

	var unmarshaled Conversation
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Conversation: %v", err)
	}

	if unmarshaled.ID != conv.ID {
		t.Errorf("Expected ID %s, got %s", conv.ID, unmarshaled.ID)
	}
	if len(unmarshaled.Messages) != len(conv.Messages) {
		t.Errorf("Expected %d messages, got %d", len(conv.Messages), len(unmarshaled.Messages))
	}
}

func TestSavedPromptStructure(t *testing.T) {
	now := time.Now()
	prompt := SavedPrompt{
		ID:          1,
		Title:       "Test Prompt",
		Content:     "This is a test prompt",
		Description: "A prompt for testing",
		Category:    "Testing",
		Tags:        `["test", "sample"]`,
		CreatedAt:   now,
		UpdatedAt:   now,
		UsageCount:  5,
	}

	jsonData, err := json.Marshal(prompt)
	if err != nil {
		t.Fatalf("Failed to marshal SavedPrompt: %v", err)
	}

	var unmarshaled SavedPrompt
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal SavedPrompt: %v", err)
	}

	if unmarshaled.Title != prompt.Title {
		t.Errorf("Expected title %s, got %s", prompt.Title, unmarshaled.Title)
	}
	if unmarshaled.UsageCount != prompt.UsageCount {
		t.Errorf("Expected usage count %d, got %d", prompt.UsageCount, unmarshaled.UsageCount)
	}
}

func TestEvalStructures(t *testing.T) {
	testCase := TestCase{
		Input:      "Test input",
		Category:   "Basic",
		Difficulty: "Easy",
		Expected:   "Expected output",
	}

	criterion := EvalCriterion{
		Name:        "Accuracy",
		Description: "How accurate is the response",
		Weight:      10,
	}

	evalData := EvalData{
		TestCases:  []TestCase{testCase},
		Criteria:   []EvalCriterion{criterion},
		BasePrompt: "Base prompt for evaluation",
		Metadata: EvalMetadata{
			GeneratedAt: time.Now(),
			Model:       "gpt-4",
			SampleSize:  10,
			Difficulty:  "Medium",
		},
	}

	jsonData, err := json.Marshal(evalData)
	if err != nil {
		t.Fatalf("Failed to marshal EvalData: %v", err)
	}

	var unmarshaled EvalData
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal EvalData: %v", err)
	}

	if len(unmarshaled.TestCases) != 1 {
		t.Errorf("Expected 1 test case, got %d", len(unmarshaled.TestCases))
	}
	if len(unmarshaled.Criteria) != 1 {
		t.Errorf("Expected 1 criterion, got %d", len(unmarshaled.Criteria))
	}
}
