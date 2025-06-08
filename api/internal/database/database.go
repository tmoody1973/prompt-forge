package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"promptforge/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("sqlite3", "./promptforge.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	database := &Database{db: db}

	// Initialize tables
	if err := database.initTables(); err != nil {
		return nil, err
	}

	return database, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) initTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		prompt TEXT NOT NULL,
		model TEXT NOT NULL,
		temperature REAL NOT NULL,
		max_tokens INTEGER,
		success BOOLEAN NOT NULL,
		response TEXT,
		error_msg TEXT
	);
	`

	_, err := d.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create history table: %v", err)
	}

	// Create conversations table
	createConversationsSQL := `
	CREATE TABLE IF NOT EXISTS conversations (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = d.db.Exec(createConversationsSQL)
	if err != nil {
		return fmt.Errorf("failed to create conversations table: %v", err)
	}

	// Create conversation_messages table
	createMessagesSQL := `
	CREATE TABLE IF NOT EXISTS conversation_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		conversation_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
	);
	`

	_, err = d.db.Exec(createMessagesSQL)
	if err != nil {
		return fmt.Errorf("failed to create conversation_messages table: %v", err)
	}

	// Create saved_prompts table
	createPromptsSQL := `
	CREATE TABLE IF NOT EXISTS saved_prompts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		description TEXT DEFAULT '',
		category TEXT DEFAULT 'General',
		tags TEXT DEFAULT '[]',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		usage_count INTEGER DEFAULT 0
	);
	`

	_, err = d.db.Exec(createPromptsSQL)
	if err != nil {
		return fmt.Errorf("failed to create saved_prompts table: %v", err)
	}

	return nil
}

func (d *Database) GetHistory() ([]models.HistoryItem, error) {
	query := `
		SELECT id, timestamp, prompt, model, temperature, max_tokens, success, response, COALESCE(error_msg, '') as error_msg
		FROM history 
		ORDER BY timestamp DESC 
		LIMIT 50
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %v", err)
	}
	defer rows.Close()

	var history []models.HistoryItem
	for rows.Next() {
		var item models.HistoryItem
		err := rows.Scan(
			&item.ID, &item.Timestamp, &item.Prompt, &item.Model,
			&item.Temperature, &item.MaxTokens, &item.Success,
			&item.Response, &item.ErrorMsg,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history row: %v", err)
		}
		history = append(history, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating history rows: %v", err)
	}

	return history, nil
}

func (d *Database) SaveHistory(req models.SaveHistoryRequest) error {
	query := `
		INSERT INTO history (prompt, model, temperature, max_tokens, success, response, error_msg)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query, req.Prompt, req.Model, req.Temperature, req.MaxTokens, req.Success, req.Response, req.ErrorMsg)
	if err != nil {
		return fmt.Errorf("failed to save history: %v", err)
	}

	return nil
}

func (d *Database) ClearHistory() error {
	query := `DELETE FROM history`

	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to clear history: %v", err)
	}

	return nil
}

// Conversation methods
func (d *Database) GetConversations() ([]models.Conversation, error) {
	query := `
		SELECT id, title, created_at, updated_at
		FROM conversations 
		ORDER BY updated_at DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations: %v", err)
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		err := rows.Scan(&conv.ID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation row: %v", err)
		}
		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating conversation rows: %v", err)
	}

	return conversations, nil
}

func (d *Database) GetConversation(conversationID string) (*models.Conversation, error) {
	// Get conversation details
	query := `SELECT id, title, created_at, updated_at FROM conversations WHERE id = ?`

	var conv models.Conversation
	err := d.db.QueryRow(query, conversationID).Scan(&conv.ID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Conversation not found
		}
		return nil, fmt.Errorf("failed to get conversation: %v", err)
	}

	// Get conversation messages
	messagesQuery := `
		SELECT id, conversation_id, role, content, timestamp
		FROM conversation_messages 
		WHERE conversation_id = ? 
		ORDER BY timestamp ASC
	`

	rows, err := d.db.Query(messagesQuery, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversation messages: %v", err)
	}
	defer rows.Close()

	var messages []models.ConversationMessage
	for rows.Next() {
		var msg models.ConversationMessage
		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &msg.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message row: %v", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating message rows: %v", err)
	}

	conv.Messages = messages
	return &conv, nil
}

func (d *Database) SaveConversation(req models.SaveConversationRequest) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Check if conversation exists
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM conversations WHERE id = ?)", req.ConversationID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check conversation existence: %v", err)
	}

	if exists {
		// Update existing conversation
		_, err = tx.Exec("UPDATE conversations SET title = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
			req.Title, req.ConversationID)
		if err != nil {
			return fmt.Errorf("failed to update conversation: %v", err)
		}

		// Clear existing messages
		_, err = tx.Exec("DELETE FROM conversation_messages WHERE conversation_id = ?", req.ConversationID)
		if err != nil {
			return fmt.Errorf("failed to clear existing messages: %v", err)
		}
	} else {
		// Create new conversation
		title := req.Title
		if title == "" {
			title = "New Conversation"
		}

		_, err = tx.Exec("INSERT INTO conversations (id, title) VALUES (?, ?)", req.ConversationID, title)
		if err != nil {
			return fmt.Errorf("failed to create conversation: %v", err)
		}
	}

	// Insert messages
	for _, msg := range req.Messages {
		_, err = tx.Exec(
			"INSERT INTO conversation_messages (conversation_id, role, content, timestamp) VALUES (?, ?, ?, ?)",
			req.ConversationID, msg.Role, msg.Content, msg.Timestamp,
		)
		if err != nil {
			return fmt.Errorf("failed to save message: %v", err)
		}
	}

	return tx.Commit()
}

func (d *Database) DeleteConversation(conversationID string) error {
	query := `DELETE FROM conversations WHERE id = ?`

	_, err := d.db.Exec(query, conversationID)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %v", err)
	}

	return nil
}

// Prompt Library methods
func (d *Database) GetSavedPrompts() ([]models.SavedPrompt, error) {
	query := `
		SELECT id, title, content, description, category, tags, created_at, updated_at, usage_count
		FROM saved_prompts 
		ORDER BY updated_at DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query saved prompts: %v", err)
	}
	defer rows.Close()

	var prompts []models.SavedPrompt
	for rows.Next() {
		var prompt models.SavedPrompt
		err := rows.Scan(
			&prompt.ID, &prompt.Title, &prompt.Content, &prompt.Description,
			&prompt.Category, &prompt.Tags, &prompt.CreatedAt, &prompt.UpdatedAt,
			&prompt.UsageCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan saved prompt row: %v", err)
		}
		prompts = append(prompts, prompt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating saved prompt rows: %v", err)
	}

	return prompts, nil
}

func (d *Database) GetSavedPrompt(promptID int64) (*models.SavedPrompt, error) {
	query := `
		SELECT id, title, content, description, category, tags, created_at, updated_at, usage_count
		FROM saved_prompts 
		WHERE id = ?
	`

	var prompt models.SavedPrompt
	err := d.db.QueryRow(query, promptID).Scan(
		&prompt.ID, &prompt.Title, &prompt.Content, &prompt.Description,
		&prompt.Category, &prompt.Tags, &prompt.CreatedAt, &prompt.UpdatedAt,
		&prompt.UsageCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Prompt not found
		}
		return nil, fmt.Errorf("failed to get saved prompt: %v", err)
	}

	return &prompt, nil
}

func (d *Database) SavePrompt(req models.SavePromptRequest) (*models.SavedPrompt, error) {
	// Convert tags to JSON string
	tagsJSON := "[]"
	if len(req.Tags) > 0 {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %v", err)
		}
		tagsJSON = string(tagsBytes)
	}

	query := `
		INSERT INTO saved_prompts (title, content, description, category, tags)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(query, req.Title, req.Content, req.Description, req.Category, tagsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to save prompt: %v", err)
	}

	promptID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	// Return the saved prompt
	return d.GetSavedPrompt(promptID)
}

func (d *Database) UpdatePrompt(req models.UpdatePromptRequest) (*models.SavedPrompt, error) {
	// Convert tags to JSON string
	tagsJSON := "[]"
	if len(req.Tags) > 0 {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %v", err)
		}
		tagsJSON = string(tagsBytes)
	}

	query := `
		UPDATE saved_prompts 
		SET title = ?, content = ?, description = ?, category = ?, tags = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := d.db.Exec(query, req.Title, req.Content, req.Description, req.Category, tagsJSON, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update prompt: %v", err)
	}

	// Return the updated prompt
	return d.GetSavedPrompt(req.ID)
}

func (d *Database) DeletePrompt(promptID int64) error {
	query := `DELETE FROM saved_prompts WHERE id = ?`

	_, err := d.db.Exec(query, promptID)
	if err != nil {
		return fmt.Errorf("failed to delete prompt: %v", err)
	}

	return nil
}

func (d *Database) IncrementPromptUsage(promptID int64) error {
	query := `UPDATE saved_prompts SET usage_count = usage_count + 1 WHERE id = ?`

	_, err := d.db.Exec(query, promptID)
	if err != nil {
		return fmt.Errorf("failed to increment prompt usage: %v", err)
	}

	return nil
}
