# PromptForge 🔨

> **A comprehensive AI prompt engineering workbench with systematic evaluation capabilities**

PromptForge is a powerful, modern tool for crafting, analyzing, and systematically evaluating AI prompts. Built with a Go backend and clean frontend, it provides everything you need for professional prompt engineering workflows.

![PromptForge Screenshot](screenshot.png)

[![License: GPLv3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/promptforge)](https://goreportcard.com/report/github.com/yourusername/promptforge)

## ✨ Features

### Core Functionality
- **🔨 Workbench Design**: Unified interface with multiple operation panels for efficient workflow
- **🔍 AI-Powered Analysis**: Dual analysis reports with both quick insights and comprehensive feedback
- **⚡ Advanced Testing**: Execute prompts with full parameter control and variable support
- **📊 Systematic Evaluations**: Generate comprehensive test suites for prompt validation
- **📚 Prompt Library**: Save, organize, and manage your prompt collections
- **🎛️ Multi-Model Support**: Works with multiple AI models (Azure OpenAI, O3, GPT-4.1)

### Advanced Capabilities
- **📈 Execution History**: Track and analyze all prompt testing sessions
- **🧪 Test Case Generation**: Automatically create diverse evaluation scenarios
- **🎯 Evaluation Criteria**: Define custom scoring metrics for systematic assessment
- **🔄 Variable Management**: Dynamic variable detection and substitution
- **💾 Data Persistence**: Local SQLite database for storing prompts and results
- **🎨 Modern UI**: Clean, responsive interface with professional design

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Azure OpenAI API access (or compatible AI service)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/promptforge.git
   cd promptforge
   ```

2. **Install dependencies**
   ```bash
   cd api
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your API credentials
   ```

4. **Run the application**
   ```bash
   # Using the provided start script
   ./start.sh
   
   # Or manually
   cd api && go run main.go
   ```

5. **Access the application**
   Navigate to `http://localhost:8080`

## 🎯 Usage Guide

### Getting Started
1. **Enter your prompt** in the main editor with syntax highlighting and line numbers
2. **Choose your operation** from the sidebar:
   - 🔍 **Get Review**: AI-powered prompt analysis and optimization suggestions
   - 🧪 **Test Prompt**: Execute with variables and advanced parameters
   - 📚 **Prompt Library**: Save and manage your prompt collections
   - 📊 **Generate Evals**: Create systematic evaluation test suites

3. **View results** in the tabbed interface:
   - **Review**: Analysis reports with improvement recommendations
   - **Execution**: AI responses with detailed parameters and settings
   - **History**: Complete execution history with filtering
   - **Library**: Organized prompt collections with search and tags
   - **Evaluations**: Generated test cases and evaluation criteria

### Evaluation System
The evaluation generator creates comprehensive test suites including:
- **Robustness Testing**: Edge cases, typos, and input variations
- **Creativity Assessment**: Novel thinking and originality scenarios
- **Safety & Alignment**: Bias detection and harmful content resistance
- **Factual Accuracy**: Correctness and reliability verification

## 🛠️ API Reference

### Core Endpoints
- `GET /api/health` - Health check and service status
- `POST /api/critique` - Single prompt analysis
- `POST /api/dual-critique` - Comprehensive dual analysis
- `POST /api/execute` - Execute prompt with parameters
- `POST /api/generate-eval` - Generate evaluation test suite

### Prompt Management
- `GET /api/prompts` - List saved prompts
- `POST /api/prompts` - Save new prompt
- `PUT /api/prompts/:id` - Update existing prompt
- `DELETE /api/prompts/:id` - Delete prompt

### Data & History
- `GET /api/history` - Execution history
- `POST /api/conversations` - Save conversation sessions
- `GET /api/conversations` - List conversation history

## ⚙️ Configuration

### Environment Variables
Create a `.env` file in the `api/` directory:

```bash
# AI Service Configuration
AZURE_OPENAI_ENDPOINT=your-azure-openai-endpoint
AZURE_API_KEY=your-api-key
AZURE_API_VERSION=2024-02-15-preview

# Server Configuration
PORT=8080
DATABASE_PATH=./promptforge.db

# Optional: Custom Model Configurations
DEFAULT_MODEL=gpt-4.1
MAX_TOKENS_LIMIT=4000
```

### Supported Models
- **GPT-4.1**: 200K context window, optimal for detailed analysis
- **O3**: 1M context window, faster execution
- **Custom**: Configure additional models via environment variables

## 🚦 Development

### Development Setup
```bash
# Clone and setup
git clone https://github.com/yourusername/promptforge.git
cd promptforge

# Install Go dependencies
cd api && go mod tidy

# Set up environment
cp .env.example .env
# Configure your .env file

# Run in development mode
go run main.go
```

### Building for Production
```bash
# Build binary
cd api
go build -o promptforge main.go

# Run production server
./promptforge
```

### Testing
```bash
# Run tests
go test ./...

# Test API endpoints
./test_api.sh
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Workflow
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Commit your changes: `git commit -m 'Add amazing feature'`
5. Push to the branch: `git push origin feature/amazing-feature`
6. Submit a Pull Request

### Areas for Contribution
- 🌐 Additional AI model integrations
- 📊 Enhanced evaluation metrics
- 🎨 UI/UX improvements
- 🔧 Performance optimizations
- 📚 Documentation improvements

## 📄 License

This project is licensed under the GPLv3 License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Built with ❤️ for the prompt engineering community**

[🌟 Star this repo](https://github.com/yourusername/promptforge) | [🐛 Report Bug](https://github.com/yourusername/promptforge/issues) | [💡 Request Feature](https://github.com/yourusername/promptforge/issues)

</div> 