package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"textarium/internal/database"
	"textarium/internal/handlers"
	"textarium/internal/services"
)

func main() {
	// Initialize database
	db, err := database.NewDatabase()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize services
	openaiService := services.NewOpenAIService()

	// Initialize handlers with dependencies
	h := handlers.NewHandlers(db, openaiService)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Serve static files
	e.Static("/", "../frontend")

	// API Routes
	api := e.Group("/api")
	api.GET("/health", h.HealthCheck)
	api.POST("/critique", h.CritiquePrompt)
	api.POST("/dual-critique", h.DualCritiquePrompt)
	api.POST("/execute", h.ExecutePrompt)
	api.POST("/prompt-engineer", h.PromptEngineer)
	api.GET("/history", h.GetHistory)
	api.POST("/history", h.SaveHistory)
	api.DELETE("/history", h.ClearHistory)

	// Conversation management routes
	api.GET("/conversations", h.GetConversations)
	api.GET("/conversations/:id", h.GetConversation)
	api.POST("/conversations", h.SaveConversation)
	api.DELETE("/conversations/:id", h.DeleteConversation)

	// Prompt Library routes
	api.GET("/prompts", h.GetSavedPrompts)
	api.GET("/prompts/:id", h.GetSavedPrompt)
	api.POST("/prompts", h.SavePrompt)
	api.PUT("/prompts/:id", h.UpdatePrompt)
	api.DELETE("/prompts/:id", h.DeletePrompt)
	api.POST("/prompts/:id/use", h.UsePrompt)

	// Eval Generator routes
	api.POST("/generate-eval", h.GenerateEval)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("📚 Textarium server starting on port %s\n", port)
	fmt.Printf("📦 Database initialized successfully\n")
	fmt.Printf("🧠 Enhanced prompt analyzer ready\n")
	fmt.Printf("🏗️  Properly structured codebase loaded\n")

	e.Logger.Fatal(e.Start(":" + port))
}
