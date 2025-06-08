package services

import (
	"fmt"
	"strings"
	"unicode"

	"promptforge/internal/models"
)

type PromptAnalyzer struct {
	aiService *UnifiedAIService
}

func NewPromptAnalyzer(aiService *UnifiedAIService) *PromptAnalyzer {
	return &PromptAnalyzer{
		aiService: aiService,
	}
}

// AnalyzePrompt performs comprehensive prompt analysis using the enhanced methodology
func (pa *PromptAnalyzer) AnalyzePrompt(prompt string, model string) (string, error) {
	// Create the enhanced critique system prompt based on the provided template
	critiqueSystemPrompt := `[Prompt] Act as a perfect prompt engineer. Your task is to analyze the given prompt and provide insights into its structure, content, and potential issues that may affect the model's response.

[Prompt] First, provide the length of the prompt (in tokens, characters, and words), any special characters used, and any potential issues with the phrasing or wording that may affect the model's response. Provide insights into the model's understanding of the prompt and any potential biases or limitations in its response.

[Prompt] Next, provide a breakdown of the prompt's Task Definition, including subtasks and objectives.

[Prompt] Then, Contextual Relevance, or how well the prompt is suited to its context. Identify relevant links to the context, and how strong links are.

[Prompt] Next, Structure Analysis, including its composition and organization.

[Prompt] Evaluate the prompt's effectiveness in achieving its intended purpose. Include the Evaluation Criteria used in the analysis, which may include factors such as clarity, specificity, relevance, and coherence or any other relevant factor.

[Prompt] Give an Audience Analysis assessing prompt's suitability for its audience, accounting for factors such as language complexity and technical jargon.

[Prompt] Last, give a Language Analysis evaluating prompt's language use and potential impact on the model's response. Include assessments of grammar, vocabulary, and style, as well as potential cultural or regional biases that may be present.

[Prompt] If you encounter any issues or errors while analyzing the prompt, document them and provide recommendations for addressing them.

CRITICAL: You MUST format your response in valid HTML only. Do not use markdown. Use these HTML tags:
- <h2>Major Section Title</h2> for main sections
- <h3>Subsection Title</h3> for subsections
- <p>Regular text</p> for paragraphs
- <ul><li>Item</li></ul> for lists
- <ol><li>Numbered item</li></ol> for numbered lists
- <strong>Important text</strong> for emphasis
- <em>Italicized text</em> for subtle emphasis
- <div class="analysis-section">Content</div> for grouping sections
- <div class="metrics">Metrics content</div> for statistical information
- <div class="recommendation">Recommendation content</div> for suggestions

[Prompt] Use **Bold** Markdown format within the HTML, including a narrated summary of your findings providing an overview of the prompt's strengths and weaknesses, and recommendations for improving its effectiveness.

Start your response with HTML tags immediately, no preamble or acknowledgment text.`

	// Add basic metrics to the analysis
	metrics := pa.calculateBasicMetrics(prompt)

	// Create the analysis request with metrics prepended
	analysisPrompt := fmt.Sprintf(`Please analyze this prompt with the following basic metrics:

PROMPT METRICS:
- Characters: %d
- Words: %d  
- Lines: %d
- Special Characters: %s

PROMPT TO ANALYZE:
%s`, metrics.Characters, metrics.Words, metrics.Lines, strings.Join(metrics.SpecialChars, ", "), prompt)

	messages := []models.Message{
		{Role: "system", Content: critiqueSystemPrompt},
		{Role: "user", Content: analysisPrompt},
	}

	if model == "" {
		model = "gpt-4.1" // Default model
	}

	response, err := pa.aiService.CallWithDefaultProvider(messages, 0.7, 2000, model)
	if err != nil {
		return "", fmt.Errorf("failed to get comprehensive analysis: %v", err)
	}

	return response, nil
}

// DualAnalyzePrompt performs both quick and detailed analysis in one go
func (pa *PromptAnalyzer) DualAnalyzePrompt(prompt string, model string) (*models.DualAnalysisData, error) {
	if model == "" {
		model = "gpt-4.1" // Default model
	}

	// Calculate basic metrics once for both analyses
	metrics := pa.calculateBasicMetrics(prompt)

	// Generate both reports concurrently
	quickReportChan := make(chan string, 1)
	detailedReportChan := make(chan string, 1)
	errorChan := make(chan error, 2)

	// Quick analysis
	go func() {
		report, err := pa.generateQuickAnalysis(prompt, metrics, model)
		if err != nil {
			errorChan <- fmt.Errorf("quick analysis failed: %v", err)
			return
		}
		quickReportChan <- report
	}()

	// Detailed analysis (reuse existing method)
	go func() {
		report, err := pa.AnalyzePrompt(prompt, model)
		if err != nil {
			errorChan <- fmt.Errorf("detailed analysis failed: %v", err)
			return
		}
		detailedReportChan <- report
	}()

	// Collect results
	var quickReport, detailedReport string
	var receivedReports int

	for receivedReports < 2 {
		select {
		case report := <-quickReportChan:
			quickReport = report
			receivedReports++
		case report := <-detailedReportChan:
			detailedReport = report
			receivedReports++
		case err := <-errorChan:
			return nil, err
		}
	}

	return &models.DualAnalysisData{
		QuickReport:    quickReport,
		DetailedReport: detailedReport,
	}, nil
}

// generateQuickAnalysis creates a succinct analysis report
func (pa *PromptAnalyzer) generateQuickAnalysis(prompt string, metrics PromptMetrics, model string) (string, error) {
	quickSystemPrompt := `You are a prompt analysis expert. Provide a QUICK, SUCCINCT analysis of the given prompt.

Keep your response focused and brief. Analyze these key aspects:
1. Overall Quality Score (1-10) 
2. Key Strengths (2-3 points max)
3. Critical Issues (2-3 points max)
4. Essential Fixes (2-3 points max)

CRITICAL: Format your response in valid HTML only. Use these tags:
- <div class="quick-analysis">Main container</div>
- <div class="score">Score: X/10</div>
- <div class="strengths"><strong>Strengths:</strong> bullet points</div>
- <div class="issues"><strong>Issues:</strong> bullet points</div>
- <div class="fixes"><strong>Essential Fixes:</strong> bullet points</div>
- <ul><li>Bullet point</li></ul> for lists
- <strong>Bold text</strong> for emphasis

Keep it concise - maximum 200 words total.`

	analysisPrompt := fmt.Sprintf(`Analyze this prompt quickly:

METRICS: %d chars, %d words, %d lines
PROMPT: %s`, metrics.Characters, metrics.Words, metrics.Lines, prompt)

	messages := []models.Message{
		{Role: "system", Content: quickSystemPrompt},
		{Role: "user", Content: analysisPrompt},
	}

	response, err := pa.aiService.CallWithDefaultProvider(messages, 0.5, 500, model)
	if err != nil {
		return "", fmt.Errorf("failed to get quick analysis: %v", err)
	}

	return response, nil
}

// PromptMetrics holds basic metrics about the prompt
type PromptMetrics struct {
	Characters   int
	Words        int
	Lines        int
	SpecialChars []string
}

// calculateBasicMetrics computes basic metrics for the prompt
func (pa *PromptAnalyzer) calculateBasicMetrics(prompt string) PromptMetrics {
	metrics := PromptMetrics{}

	// Character count
	metrics.Characters = len(prompt)

	// Word count
	words := strings.Fields(prompt)
	metrics.Words = len(words)

	// Line count
	lines := strings.Split(prompt, "\n")
	metrics.Lines = len(lines)

	// Special characters detection
	specialChars := make(map[string]bool)
	for _, char := range prompt {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			if char != ' ' && char != '\n' && char != '\t' {
				charStr := string(char)
				if !specialChars[charStr] {
					specialChars[charStr] = true
				}
			}
		}
	}

	// Convert map to slice
	for char := range specialChars {
		metrics.SpecialChars = append(metrics.SpecialChars, char)
	}

	if len(metrics.SpecialChars) == 0 {
		metrics.SpecialChars = []string{"None detected"}
	}

	return metrics
}
