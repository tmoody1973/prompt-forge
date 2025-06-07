# PromptForge API - Enhanced & Restructured

A powerful prompt engineering workbench with comprehensive analysis capabilities and a properly structured Go codebase.

## 🏗️ Architecture

The codebase is now properly organized using Go best practices with clean separation of concerns:

```
api/
├── main.go                     # Application entry point
├── internal/                   # Internal packages
│   ├── config/                 # Configuration management
│   │   └── config.go          # Azure OpenAI settings & model mappings
│   ├── models/                 # Data structures
│   │   └── models.go          # Request/Response types & database models
│   ├── services/              # Business logic
│   │   ├── openai.go          # Azure OpenAI API client
│   │   └── prompt_analyzer.go # Enhanced prompt analysis engine
│   ├── database/              # Data persistence
│   │   └── database.go        # SQLite operations & history management
│   └── handlers/              # HTTP handlers
│       └── handlers.go        # REST API endpoints
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## 🧠 Enhanced Prompt Analyzer

The prompt reviewer has been completely upgraded with comprehensive analysis capabilities inspired by advanced prompt engineering methodologies:

### Analysis Dimensions

1. **Prompt Metrics**: Characters, words, lines, special characters
2. **Task Definition**: Subtasks and objectives breakdown
3. **Contextual Relevance**: Context suitability and link strength
4. **Structure Analysis**: Composition and organization evaluation
5. **Evaluation Criteria**: Clarity, specificity, relevance, coherence
6. **Audience Analysis**: Language complexity and technical jargon assessment
7. **Language Analysis**: Grammar, vocabulary, style, and cultural bias detection

### Key Features

- **Multi-part prompt structure** using `[Prompt]` tokens for better AI comprehension
- **Automatic metrics calculation** for quantitative analysis
- **HTML-formatted responses** for rich presentation
- **Comprehensive error documentation** with improvement recommendations
- **Advanced prompt engineering methodology** based on industry best practices

## 🚀 Quick Start

```bash
# Build the application
go build -buildvcs=false -o promptforge-server

# Run the server
./promptforge-server

# Or use the provided script
./start.sh
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| POST | `/api/critique` | Enhanced prompt analysis |
| POST | `/api/execute` | Execute prompt |
| GET | `/api/history` | Get execution history |
| POST | `/api/history` | Save execution to history |
| DELETE | `/api/history` | Clear history |

## 🔧 Configuration

Configuration is centralized in `internal/config/config.go`:

- Azure OpenAI endpoint and API key
- Model deployment mappings
- API version settings

## 🎯 Enhanced Prompt Analysis

The new analysis engine provides detailed insights:

```json
{
  "success": true,
  "data": "<h2>Prompt Analysis</h2><div class='metrics'>...</div>..."
}
```

### Analysis Process

1. **Metrics Calculation**: Automatic computation of basic prompt statistics
2. **Multi-dimensional Analysis**: Comprehensive evaluation across 7 key dimensions
3. **Structured Output**: HTML-formatted response with clear sections
4. **Actionable Recommendations**: Specific suggestions for improvement

## 🏛️ Architecture Benefits

### Separation of Concerns
- **Config**: Centralized configuration management
- **Models**: Clean data structure definitions
- **Services**: Business logic isolation
- **Database**: Data persistence abstraction
- **Handlers**: HTTP layer separation

### Maintainability
- Single responsibility principle
- Dependency injection
- Clear interfaces
- Modular structure

### Testability
- Isolated components
- Mock-friendly interfaces
- Clear dependencies

### Scalability
- Pluggable services
- Configurable components
- Easy feature additions

## 🛠️ Development

### Adding New Features

1. Define models in `internal/models/`
2. Implement business logic in `internal/services/`
3. Add HTTP handlers in `internal/handlers/`
4. Wire up in `main.go`

### Database Operations

All database operations are centralized in `internal/database/database.go`:
- Connection management
- Table initialization
- CRUD operations
- Error handling

## 📊 Enhanced Analysis Example

The new prompt analyzer provides comprehensive feedback like:

- **Quantitative metrics** (length, complexity)
- **Structural assessment** (organization, clarity)
- **Contextual evaluation** (relevance, appropriateness)
- **Language analysis** (grammar, style, bias)
- **Improvement recommendations** (specific, actionable)

## 🔄 Migration from Previous Version

The restructured codebase maintains full API compatibility while providing:
- Better code organization
- Enhanced analysis capabilities
- Improved maintainability
- Cleaner separation of concerns

## 🤝 Contributing

1. Follow the established package structure
2. Maintain separation of concerns
3. Add appropriate error handling
4. Update documentation as needed

---

**Note**: This enhanced version includes comprehensive prompt analysis capabilities using advanced prompt engineering methodologies, providing detailed insights across multiple dimensions for optimal prompt optimization. 