package database

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"promptforge/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *Database {
	// Create temporary database file
	tmpfile := "test_" + time.Now().Format("20060102150405") + ".db"

	db, err := sql.Open("sqlite3", tmpfile)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	database := &Database{db: db}

	// Initialize tables
	if err := database.initTables(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	// Clean up function
	t.Cleanup(func() {
		database.Close()
		os.Remove(tmpfile)
	})

	return database
}

func TestNewDatabase(t *testing.T) {
	// Test with a temporary database
	tmpfile := "test_new_db.db"
	defer os.Remove(tmpfile)

	// Mock the database path temporarily
	originalPath := "./promptforge.db"

	// Create a new database instance
	db, err := sql.Open("sqlite3", tmpfile)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	database := &Database{db: db}
	err = database.initTables()
	if err != nil {
		t.Fatalf("Failed to initialize tables: %v", err)
	}

	defer database.Close()

	// Test that tables were created by trying to insert data
	_, err = database.db.Exec("INSERT INTO history (prompt, model, temperature, max_tokens, success, response) VALUES (?, ?, ?, ?, ?, ?)",
		"test", "gpt-4", 0.7, 100, true, "test response")
	if err != nil {
		t.Errorf("Failed to insert into history table: %v", err)
	}

	_ = originalPath // Avoid unused variable warning
}

func TestHistoryOperations(t *testing.T) {
	db := setupTestDB(t)

	// Test SaveHistory
	saveReq := models.SaveHistoryRequest{
		Prompt:      "Test prompt",
		Model:       "gpt-4",
		Temperature: 0.7,
		MaxTokens:   100,
		Success:     true,
		Response:    "Test response",
		ErrorMsg:    "",
	}

	err := db.SaveHistory(saveReq)
	if err != nil {
		t.Fatalf("Failed to save history: %v", err)
	}

	// Test GetHistory
	history, err := db.GetHistory()
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}

	if len(history) != 1 {
		t.Errorf("Expected 1 history item, got %d", len(history))
	}

	item := history[0]
	if item.Prompt != saveReq.Prompt {
		t.Errorf("Expected prompt '%s', got '%s'", saveReq.Prompt, item.Prompt)
	}
	if item.Model != saveReq.Model {
		t.Errorf("Expected model '%s', got '%s'", saveReq.Model, item.Model)
	}
	if item.Temperature != saveReq.Temperature {
		t.Errorf("Expected temperature %f, got %f", saveReq.Temperature, item.Temperature)
	}
	if item.Success != saveReq.Success {
		t.Errorf("Expected success %v, got %v", saveReq.Success, item.Success)
	}

	// Test ClearHistory
	err = db.ClearHistory()
	if err != nil {
		t.Fatalf("Failed to clear history: %v", err)
	}

	history, err = db.GetHistory()
	if err != nil {
		t.Fatalf("Failed to get history after clear: %v", err)
	}

	if len(history) != 0 {
		t.Errorf("Expected 0 history items after clear, got %d", len(history))
	}
}

func TestSavedPromptOperations(t *testing.T) {
	db := setupTestDB(t)

	// Test SavePrompt
	saveReq := models.SavePromptRequest{
		Title:       "Test Prompt",
		Content:     "This is a test prompt for unit testing",
		Description: "A test description",
		Category:    "Testing",
		Tags:        []string{"test", "unit", "sample"},
	}

	savedPrompt, err := db.SavePrompt(saveReq)
	if err != nil {
		t.Fatalf("Failed to save prompt: %v", err)
	}

	if savedPrompt.ID == 0 {
		t.Error("Expected saved prompt to have non-zero ID")
	}
	if savedPrompt.Title != saveReq.Title {
		t.Errorf("Expected title '%s', got '%s'", saveReq.Title, savedPrompt.Title)
	}

	promptID := savedPrompt.ID

	// Test GetSavedPrompts
	prompts, err := db.GetSavedPrompts()
	if err != nil {
		t.Fatalf("Failed to get saved prompts: %v", err)
	}

	if len(prompts) != 1 {
		t.Errorf("Expected 1 saved prompt, got %d", len(prompts))
	}

	// Test GetSavedPrompt
	prompt, err := db.GetSavedPrompt(promptID)
	if err != nil {
		t.Fatalf("Failed to get saved prompt by ID: %v", err)
	}

	if prompt.ID != promptID {
		t.Errorf("Expected prompt ID %d, got %d", promptID, prompt.ID)
	}

	// Test UpdatePrompt
	updateReq := models.UpdatePromptRequest{
		ID:          promptID,
		Title:       "Updated Test Prompt",
		Content:     "This is an updated test prompt",
		Description: "Updated description",
		Category:    "Updated Testing",
		Tags:        []string{"updated", "test"},
	}

	updatedPrompt, err := db.UpdatePrompt(updateReq)
	if err != nil {
		t.Fatalf("Failed to update prompt: %v", err)
	}

	if updatedPrompt.Title != updateReq.Title {
		t.Errorf("Expected updated title '%s', got '%s'", updateReq.Title, updatedPrompt.Title)
	}

	// Test IncrementPromptUsage
	originalUsage := updatedPrompt.UsageCount
	err = db.IncrementPromptUsage(promptID)
	if err != nil {
		t.Fatalf("Failed to increment prompt usage: %v", err)
	}

	// Verify usage was incremented
	prompt, err = db.GetSavedPrompt(promptID)
	if err != nil {
		t.Fatalf("Failed to get prompt after usage increment: %v", err)
	}

	if prompt.UsageCount != originalUsage+1 {
		t.Errorf("Expected usage count %d, got %d", originalUsage+1, prompt.UsageCount)
	}

	// Test DeletePrompt
	err = db.DeletePrompt(promptID)
	if err != nil {
		t.Fatalf("Failed to delete prompt: %v", err)
	}

	// Verify prompt was deleted
	prompt, err = db.GetSavedPrompt(promptID)
	if err != nil {
		t.Fatalf("Unexpected error when getting deleted prompt: %v", err)
	}
	if prompt != nil {
		t.Error("Expected nil when getting deleted prompt")
	}
}

func TestConversationOperations(t *testing.T) {
	db := setupTestDB(t)

	// Test SaveConversation
	saveReq := models.SaveConversationRequest{
		ConversationID: "conv-123",
		Title:          "Test Conversation",
		Messages: []models.ConversationMessage{
			{
				ConversationID: "conv-123",
				Role:           "user",
				Content:        "Hello",
			},
			{
				ConversationID: "conv-123",
				Role:           "assistant",
				Content:        "Hi there!",
			},
		},
	}

	err := db.SaveConversation(saveReq)
	if err != nil {
		t.Fatalf("Failed to save conversation: %v", err)
	}

	// Test GetConversations
	conversations, err := db.GetConversations()
	if err != nil {
		t.Fatalf("Failed to get conversations: %v", err)
	}

	if len(conversations) != 1 {
		t.Errorf("Expected 1 conversation, got %d", len(conversations))
	}

	if conversations[0].ID != saveReq.ConversationID {
		t.Errorf("Expected conversation ID '%s', got '%s'", saveReq.ConversationID, conversations[0].ID)
	}

	// Test GetConversation
	conversation, err := db.GetConversation(saveReq.ConversationID)
	if err != nil {
		t.Fatalf("Failed to get conversation by ID: %v", err)
	}

	if conversation.ID != saveReq.ConversationID {
		t.Errorf("Expected conversation ID '%s', got '%s'", saveReq.ConversationID, conversation.ID)
	}

	if len(conversation.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(conversation.Messages))
	}

	// Test DeleteConversation
	err = db.DeleteConversation(saveReq.ConversationID)
	if err != nil {
		t.Fatalf("Failed to delete conversation: %v", err)
	}

	// Verify conversation was deleted
	conversation, err = db.GetConversation(saveReq.ConversationID)
	if err != nil {
		t.Fatalf("Unexpected error when getting deleted conversation: %v", err)
	}
	if conversation != nil {
		t.Error("Expected nil when getting deleted conversation")
	}
}

func TestHistoryLimit(t *testing.T) {
	db := setupTestDB(t)

	// Insert more than 50 history items to test the limit
	for i := 0; i < 60; i++ {
		saveReq := models.SaveHistoryRequest{
			Prompt:      "Test prompt " + string(rune(i)),
			Model:       "gpt-4",
			Temperature: 0.7,
			MaxTokens:   100,
			Success:     true,
			Response:    "Test response " + string(rune(i)),
		}

		err := db.SaveHistory(saveReq)
		if err != nil {
			t.Fatalf("Failed to save history item %d: %v", i, err)
		}
	}

	// Get history and verify limit
	history, err := db.GetHistory()
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}

	if len(history) != 50 {
		t.Errorf("Expected history to be limited to 50 items, got %d", len(history))
	}
}

func TestDatabaseErrorHandling(t *testing.T) {
	db := setupTestDB(t)

	// Test GetSavedPrompt with non-existent ID
	prompt, err := db.GetSavedPrompt(99999)
	if err != nil {
		t.Errorf("Unexpected error when getting non-existent prompt: %v", err)
	}
	if prompt != nil {
		t.Error("Expected nil when getting non-existent prompt")
	}

	// Test GetConversation with non-existent ID
	conversation, err := db.GetConversation("non-existent")
	if err != nil {
		t.Errorf("Unexpected error when getting non-existent conversation: %v", err)
	}
	if conversation != nil {
		t.Error("Expected nil when getting non-existent conversation")
	}

	// Test DeletePrompt with non-existent ID (should not error)
	err = db.DeletePrompt(99999)
	if err != nil {
		t.Errorf("DeletePrompt should not error for non-existent ID, got: %v", err)
	}

	// Test DeleteConversation with non-existent ID (should not error)
	err = db.DeleteConversation("non-existent")
	if err != nil {
		t.Errorf("DeleteConversation should not error for non-existent ID, got: %v", err)
	}
}

func TestJSONTagHandling(t *testing.T) {
	db := setupTestDB(t)

	// Test saving and retrieving a prompt with tags
	saveReq := models.SavePromptRequest{
		Title:   "Test Prompt with Tags",
		Content: "Test content",
		Tags:    []string{"tag1", "tag2", "tag3"},
	}

	savedPrompt, err := db.SavePrompt(saveReq)
	if err != nil {
		t.Fatalf("Failed to save prompt with tags: %v", err)
	}

	// Retrieve and verify tags were stored as JSON
	prompt, err := db.GetSavedPrompt(savedPrompt.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve prompt: %v", err)
	}

	// The tags should be stored as JSON string
	if prompt.Tags == "" {
		t.Error("Expected tags to be stored as JSON string")
	}

	// Should contain the JSON representation of the tags
	expectedSubstrings := []string{"tag1", "tag2", "tag3"}
	for _, substr := range expectedSubstrings {
		if !contains(prompt.Tags, substr) {
			t.Errorf("Expected tags JSON to contain '%s', got: %s", substr, prompt.Tags)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
