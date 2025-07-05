# PromptForge 🔨

**The AI prompt engineering workbench that brings systematic testing to prompt development**

Transform your prompt engineering from guesswork into a disciplined craft. Built with Go for speed and reliability.

**[🎮 Try Interactive Demo](https://demo.arcade.software/bFGTYb7AuRV33Kei7ZFQ)**

![PromptForge Screenshot](screenshot.png)

## Why PromptForge?

Most prompt tools are glorified text editors. PromptForge brings **engineering discipline** to prompt development:

- **🧪 Systematic Testing** - Generate comprehensive test suites automatically
- **📊 Built-in Analytics** - Track what works across different scenarios
- **🔍 AI-Powered Analysis** - Get feedback before you even test
- **🎯 Variable Testing** - Test edge cases and consistency
- **📚 Version Control** - Never lose a working prompt again

## Quick Start

**One-line Docker setup:**
```bash
docker run -d -p 8080:8080 -e ANTHROPIC_API_KEY="your-key" promptforge:latest
```

**Or clone and run:**
```bash
git clone https://github.com/insaanimanav/promptforge.git
cd promptforge && ./start.sh
```

Open `http://localhost:8080` and start crafting better prompts.

## Core Features

### 🔨 Workbench Design
Unified interface with multiple operation panels for efficient workflow

### 🔍 Dual Analysis System
- **Quick Review**: Instant optimization suggestions
- **Deep Analysis**: Comprehensive prompt evaluation

### ⚡ Advanced Testing
- Execute with full parameter control
- Dynamic variable detection and substitution
- Multi-model support (Claude, GPT-4, Azure OpenAI)

### 📊 Evaluation Generator
Creates comprehensive test suites including:
- **Robustness**: Edge cases, typos, input variations
- **Safety**: Bias detection, harmful content resistance
- **Accuracy**: Factual correctness verification
- **Creativity**: Novel thinking scenarios

### 📚 Prompt Management
- Organized library with search and tags
- Complete execution history with filtering
- Export/import capabilities

## Supported Models

- **Claude 3.5 Sonnet** (200K context) - Excellent reasoning
- **GPT-4.1** (200K context) - Detailed analysis
- **O3** (1M context) - Fast execution
- **Azure OpenAI** - Enterprise-ready

## Configuration

### Docker (Recommended)
```bash
# Anthropic (Default)
-e ANTHROPIC_API_KEY="sk-ant-api03-..."

# OpenAI  
-e OPENAI_API_KEY="sk-..."

# Azure OpenAI
-e AZURE_OPENAI_API_KEY="your-key"
-e AZURE_OPENAI_BASE_URL="https://your-resource.openai.azure.com"
```

### Local Development
```bash
# Clone and setup
git clone https://github.com/insaanimanav/promptforge.git
cd promptforge

# Configure environment
cp .env.example .env
# Add your API keys to .env

# Run
cd api && go run main.go
```

## API Reference

Core endpoints for integration:
- `POST /api/critique` - Analyze prompts
- `POST /api/execute` - Test with parameters
- `POST /api/generate-eval` - Create test suites
- `GET /api/prompts` - Manage prompt library

## Contributing

Built for the prompt engineering community. Contributions welcome:

- 🌐 Additional AI model integrations
- 📊 Enhanced evaluation metrics
- 🎨 UI/UX improvements
- 🔧 Performance optimizations

1. Fork the repository
2. Create feature branch: `git checkout -b feature/name`
3. Submit pull request

## License

GPLv3 - Built with ❤️ for better AI interactions

---

**Transform prompt engineering from art to science**