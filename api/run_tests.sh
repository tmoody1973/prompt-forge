#!/bin/bash

echo "🧪 Running PromptForge Backend Unit Tests"
echo "========================================"

# Change to the API directory
cd "$(dirname "$0")"

# Initialize Go modules if needed
echo "📦 Ensuring dependencies are up to date..."
go mod tidy

echo ""
echo "🔧 Running tests for config package..."
go test -v ./internal/config/

echo ""
echo "📊 Running tests for models package..."
go test -v ./internal/models/

echo ""
echo "🤖 Running tests for services package..."
go test -v ./internal/services/

echo ""
echo "💾 Running tests for database package..."
go test -v ./internal/database/

echo ""
echo "🌐 Running tests for handlers package..."
go test -v ./internal/handlers/

echo ""
echo "📈 Running all tests with coverage..."
go test -v -coverprofile=coverage.out ./...

echo ""
echo "📊 Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo ""
echo "✅ Test Summary:"
echo "=================="
echo "📁 Config tests: Configuration management and environment variables"
echo "📊 Models tests: Data structures and JSON serialization"
echo "🤖 Services tests: AI service logic and provider handling"
echo "💾 Database tests: CRUD operations and data persistence"
echo "🌐 Handlers tests: HTTP request/response handling"
echo ""
echo "📈 Coverage report generated: coverage.html"
echo "🎉 All tests completed!" 