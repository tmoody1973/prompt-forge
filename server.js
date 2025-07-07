// Load environment variables from .env file
require('dotenv').config();

const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 8080;

// Middleware
app.use(cors());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

// Serve static files from frontend directory
app.use(express.static(path.join(__dirname, 'frontend')));

// API Routes (Mock endpoints for demo)
app.get('/api/health', (req, res) => {
    res.json({
        status: 'healthy',
        service: 'PromptForge API',
        timestamp: new Date().toISOString()
    });
});

app.get('/api/providers', (req, res) => {
    res.json({
        default: 'anthropic',
        available: ['openai', 'azure-openai', 'anthropic'],
        configured: {
            openai: !!process.env.OPENAI_API_KEY,
            'azure-openai': !!process.env.AZURE_OPENAI_API_KEY,
            anthropic: !!process.env.ANTHROPIC_API_KEY
        }
    });
});

app.post('/api/critique', (req, res) => {
    const { prompt } = req.body;
    
    // Mock response for demo
    res.json({
        success: true,
        data: `
            <div class="analysis-section">
                <h2>Prompt Analysis Results</h2>
                <div class="metrics">
                    <strong>Length:</strong> ${prompt?.length || 0} characters<br>
                    <strong>Estimated Tokens:</strong> ${Math.ceil((prompt?.length || 0) / 4)}<br>
                    <strong>Complexity:</strong> Medium
                </div>
                
                <h3>Strengths</h3>
                <ul>
                    <li>Clear objective statement</li>
                    <li>Appropriate context provided</li>
                    <li>Well-structured format</li>
                </ul>
                
                <h3>Recommendations</h3>
                <div class="recommendation">
                    <strong>Specificity:</strong> Consider adding more specific examples or constraints to guide the AI's response more precisely.
                </div>
                
                <div class="recommendation">
                    <strong>Output Format:</strong> Specify the desired output format (bullet points, paragraphs, structured data) for more consistent results.
                </div>
            </div>
        `
    });
});

app.post('/api/dual-critique', (req, res) => {
    const { prompt } = req.body;
    
    res.json({
        success: true,
        data: {
            quick_report: `
                <div class="quick-analysis">
                    <div class="score">Score: 7/10</div>
                    <div class="strengths">
                        <strong>Strengths:</strong>
                        <ul>
                            <li>Clear and concise</li>
                            <li>Good structure</li>
                        </ul>
                    </div>
                    <div class="issues">
                        <strong>Issues:</strong>
                        <ul>
                            <li>Could be more specific</li>
                            <li>Missing output format</li>
                        </ul>
                    </div>
                    <div class="fixes">
                        <strong>Essential Fixes:</strong>
                        <ul>
                            <li>Add specific examples</li>
                            <li>Define output structure</li>
                        </ul>
                    </div>
                </div>
            `,
            detailed_report: `
                <div class="analysis-section">
                    <h2>Comprehensive Prompt Analysis</h2>
                    
                    <h3>Task Definition</h3>
                    <p>The prompt establishes a clear objective with well-defined parameters. The task scope is appropriate for the intended AI model capabilities.</p>
                    
                    <h3>Contextual Relevance</h3>
                    <p>Context is provided but could be enhanced with domain-specific examples. The prompt maintains good relevance to the stated objective.</p>
                    
                    <h3>Structure Analysis</h3>
                    <div class="metrics">
                        <strong>Organization:</strong> Good<br>
                        <strong>Flow:</strong> Logical progression<br>
                        <strong>Clarity:</strong> High
                    </div>
                    
                    <h3>Language Analysis</h3>
                    <p>The language is clear and professional. Grammar and vocabulary are appropriate for the target audience. No significant cultural biases detected.</p>
                    
                    <div class="recommendation">
                        <strong>Overall Recommendation:</strong> This is a solid prompt that would benefit from minor refinements in specificity and output formatting. Consider adding 1-2 concrete examples to guide the AI's understanding.
                    </div>
                </div>
            `
        }
    });
});

app.post('/api/execute', (req, res) => {
    const { prompt, model, temperature, max_tokens } = req.body;
    
    // Mock execution response
    setTimeout(() => {
        res.json({
            success: true,
            data: `This is a mock response for the prompt: "${prompt?.substring(0, 50)}..."\n\nModel: ${model || 'gpt-4.1'}\nTemperature: ${temperature || 0.7}\nMax Tokens: ${max_tokens || 1000}\n\nIn a real deployment, this would connect to your configured AI provider (OpenAI, Anthropic, or Azure OpenAI) and return the actual AI response.\n\nTo enable real AI responses, configure your API keys in the environment variables:\n- OPENAI_API_KEY\n- ANTHROPIC_API_KEY\n- AZURE_OPENAI_API_KEY`
        });
    }, 1500);
});

app.post('/api/multi-model-execute', (req, res) => {
    const { prompt, models, temperature, max_tokens } = req.body;
    
    // Mock multi-model response
    setTimeout(() => {
        const results = models.map((model, index) => ({
            model: model,
            success: true,
            response: `Mock response from ${model} for: "${prompt?.substring(0, 30)}..."\n\nThis demonstrates how different models might respond to the same prompt with varying styles and approaches.`,
            execution_time_ms: 800 + (index * 200),
            token_usage: {
                prompt_tokens: Math.ceil((prompt?.length || 0) / 4),
                completion_tokens: 150 + (index * 50),
                total_tokens: Math.ceil((prompt?.length || 0) / 4) + 150 + (index * 50)
            }
        }));
        
        res.json({
            success: true,
            data: results
        });
    }, 2000);
});

// Mock endpoints for other features
app.get('/api/history', (req, res) => {
    res.json({
        success: true,
        data: [
            {
                id: 1,
                timestamp: new Date().toISOString(),
                prompt: "Write a professional email...",
                model: "gpt-4.1",
                temperature: 0.7,
                max_tokens: 1000,
                success: true,
                response: "Mock historical response"
            }
        ]
    });
});

app.post('/api/history', (req, res) => {
    res.json({ success: true, data: "History saved successfully" });
});

app.delete('/api/history', (req, res) => {
    res.json({ success: true, data: "History cleared successfully" });
});

app.get('/api/prompts', (req, res) => {
    res.json({
        success: true,
        data: [
            {
                id: 1,
                title: "Email Writer",
                content: "Write a professional email that...",
                description: "Template for professional emails",
                category: "Communication",
                tags: '["email", "professional", "template"]',
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
                usage_count: 5
            }
        ]
    });
});

app.post('/api/prompts', (req, res) => {
    res.json({
        success: true,
        data: {
            id: Date.now(),
            ...req.body,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            usage_count: 0
        }
    });
});

app.get('/api/prompts/:id', (req, res) => {
    const { id } = req.params;
    
    // Mock response for a specific prompt
    res.json({
        success: true,
        data: {
            id: parseInt(id),
            title: "Sample Prompt",
            content: "This is a sample prompt content for ID " + id,
            description: "A sample prompt loaded from the server",
            category: "General",
            tags: '["sample", "demo", "test"]',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            usage_count: 1
        }
    });
});

app.get('/api/conversations', (req, res) => {
    res.json({
        success: true,
        data: []
    });
});

app.post('/api/conversations', (req, res) => {
    res.json({ success: true, data: "Conversation saved successfully" });
});

app.post('/api/generate-eval', (req, res) => {
    const { prompt, eval_types, sample_size } = req.body;
    
    setTimeout(() => {
        res.json({
            success: true,
            data: {
                test_cases: Array.from({ length: sample_size || 5 }, (_, i) => ({
                    input: `Test case ${i + 1} for: ${prompt?.substring(0, 30)}...`,
                    category: eval_types?.[i % eval_types.length] || 'robustness',
                    difficulty: ['easy', 'medium', 'hard'][i % 3]
                })),
                criteria: eval_types?.map(type => ({
                    name: type.charAt(0).toUpperCase() + type.slice(1),
                    description: `Evaluation criteria for ${type}`,
                    weight: Math.floor(100 / eval_types.length)
                })) || [],
                base_prompt: prompt,
                metadata: {
                    generated_at: new Date().toISOString(),
                    model: 'gpt-4.1',
                    sample_size: sample_size || 5,
                    eval_types: eval_types || ['robustness'],
                    difficulty: 'mixed'
                }
            }
        });
    }, 2000);
});

// Catch all handler: send back React's index.html file for client-side routing
app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, 'frontend', 'index.html'));
});

app.listen(PORT, () => {
    console.log('🔨 PromptForge Server Started');
    console.log('================================');
    console.log(`📍 Server running on port ${PORT}`);
    console.log(`🌐 Open: http://localhost:${PORT}`);
    console.log('🔍 API Health: /api/health');
    console.log('⚡ Ready for prompt engineering!');
    console.log('================================');
    
    // Log environment status
    const providers = {
        'OpenAI': !!process.env.OPENAI_API_KEY,
        'Anthropic': !!process.env.ANTHROPIC_API_KEY,
        'Azure OpenAI': !!process.env.AZURE_OPENAI_API_KEY
    };
    
    console.log('\n🤖 AI Provider Status:');
    Object.entries(providers).forEach(([name, configured]) => {
        console.log(`   ${configured ? '✅' : '❌'} ${name}: ${configured ? 'Configured' : 'Not configured'}`);
    });
    
    if (!Object.values(providers).some(Boolean)) {
        console.log('\n⚠️  No AI providers configured. Add API keys to environment variables:');
        console.log('   - OPENAI_API_KEY="sk-..."');
        console.log('   - ANTHROPIC_API_KEY="sk-ant-..."');
        console.log('   - AZURE_OPENAI_API_KEY="your-key"');
        console.log('\n📝 Currently running in demo mode with mock responses.');
    }
    
    console.log('\n🚀 PromptForge is ready!');
});