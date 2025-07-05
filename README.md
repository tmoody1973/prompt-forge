# PromptForge 🔨

> **A comprehensive AI prompt engineering workbench with systematic evaluation capabilities**

PromptForge is a powerful, modern tool for crafting, analyzing, and systematically evaluating AI prompts. Built with a Go backend and clean frontend, it provides everything you need for professional prompt engineering workflows.

<!--ARCADE EMBED START--><div style="position: relative; padding-bottom: calc(50.681341719077565% + 41px); height: 0; width: 100%;"><iframe src="https://demo.arcade.software/bFGTYb7AuRV33Kei7ZFQ?embed&embed_mobile=inline&embed_desktop=inline&show_copy_link=true" title="Run Prompt Evaluations in PromptForge" frameborder="0" loading="lazy" webkitallowfullscreen mozallowfullscreen allowfullscreen allow="clipboard-write" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; color-scheme: light;" ></iframe></div><!--ARCADE EMBED END-->

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
- **🐳 Docker Ready**: One-command deployment with Docker/Podman support
- **🔧 Multi-Provider**: Support for Anthropic, OpenAI, and Azure OpenAI

![PromptForge Interface](image.png)

## 🚀 Quick Start

### Prerequisites
- **Option 1 (Docker)**: Docker or Podman installed
- **Option 2 (Local)**: Go 1.21 or higher
- **AI Service**: Anthropic, OpenAI, or Azure OpenAI API access

## 🐳 Docker Deployment (Recommended)

### Using Docker
```bash
# 1. Clone the repository
git clone https://github.com/insaanimanav/promptforge.git
cd promptforge

# 2. Build the Docker image
docker build -t promptforge:latest -f Dockerfile .

# 3. Run with your API key
docker run -d \
  --name promptforge \
  -p 8080:8080 \
  -e ANTHROPIC_API_KEY="your-api-key-here" \
  promptforge:latest

# 4. Access the application
# Frontend: http://localhost:8080
# API: http://localhost:8080/api
```

### Using Podman
```bash
# 1. Clone the repository
git clone https://github.com/insaanimanav/promptforge.git
cd promptforge

# 2. Build the Podman image
podman build -t promptforge:latest -f Dockerfile .

# 3. Run with your API key
podman run -d \
  --name promptforge \
  -p 8080:8080 \
  -e ANTHROPIC_API_KEY="your-api-key-here" \
  promptforge:latest

# 4. Access the application
# Frontend: http://localhost:8080
# API: http://localhost:8080/api
```

### Docker Environment Variables
Configure AI providers using environment variables:

```bash
# Anthropic (Default)
-e ANTHROPIC_API_KEY="sk-ant-api03-..."

# OpenAI
-e OPENAI_API_KEY="sk-..."

# Azure OpenAI
-e AZURE_OPENAI_API_KEY="your-key"
-e AZURE_OPENAI_BASE_URL="https://your-resource.openai.azure.com"
-e AZURE_OPENAI_API_VERSION="2024-02-15-preview"

# Set default provider
-e DEFAULT_AI_PROVIDER="anthropic"  # or "openai", "azure-openai"
```

### Docker Compose (Alternative)
Create a `docker-compose.yml` file for easier management:

```yaml
version: '3.8'
services:
  promptforge:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ANTHROPIC_API_KEY=your-api-key-here
      - DEFAULT_AI_PROVIDER=anthropic
    volumes:
      - promptforge_data:/data
    restart: unless-stopped

volumes:
  promptforge_data:
```

Run with Docker Compose:
```bash
# Start the application
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the application
docker-compose down
```

### Container Management
```bash
# View running containers
docker ps

# View logs
docker logs promptforge

# Stop container
docker stop promptforge

# Remove container
docker rm promptforge

# Health check
curl http://localhost:8080/api/health
```

## 🏠 Local Development

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

### Docker Environment Variables
When running with Docker/Podman, configure using environment variables:

```bash
# AI Providers (choose one or more)
ANTHROPIC_API_KEY=sk-ant-api03-...           # Anthropic Claude
OPENAI_API_KEY=sk-...                        # OpenAI GPT models
AZURE_OPENAI_API_KEY=your-key                # Azure OpenAI
AZURE_OPENAI_BASE_URL=https://your-resource.openai.azure.com
AZURE_OPENAI_API_VERSION=2024-02-15-preview

# Default Provider Selection
DEFAULT_AI_PROVIDER=anthropic                # anthropic, openai, or azure-openai

# Server Configuration (optional)
PORT=8080                                    # Server port
DATABASE_PATH=/data/promptforge.db           # Database location
```

### Local Development (.env file)
For local development, create a `.env` file in the `api/` directory:

```bash
# AI Service Configuration
ANTHROPIC_API_KEY=your-anthropic-api-key
OPENAI_API_KEY=your-openai-api-key
AZURE_OPENAI_ENDPOINT=your-azure-openai-endpoint
AZURE_API_KEY=your-api-key
AZURE_API_VERSION=2024-02-15-preview

# Server Configuration
PORT=8080
DATABASE_PATH=./promptforge.db
DEFAULT_AI_PROVIDER=anthropic

# Optional: Custom Model Configurations
DEFAULT_MODEL=claude-3-sonnet-20240229
MAX_TOKENS_LIMIT=4000
```

### Supported Models
- **Claude 3.5 Sonnet**: 200K context, excellent reasoning and coding
- **Claude 3 Haiku**: 200K context, fast and cost-effective
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