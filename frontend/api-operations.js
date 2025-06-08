// Simple markdown to HTML converter
function markdownToHtml(markdown) {
    return markdown
        .replace(/^### (.*$)/gim, '<h3>$1</h3>')
        .replace(/^## (.*$)/gim, '<h2>$1</h2>')
        .replace(/^# (.*$)/gim, '<h1>$1</h1>')
        .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
        .replace(/\*(.*?)\*/g, '<em>$1</em>')
        .replace(/^\s*[\-\*\+]\s+(.*)$/gim, '<li>$1</li>')
        .replace(/(<li>.*<\/li>)/s, '<ul>$1</ul>')
        .replace(/\n\n/g, '</p><p>')
        .replace(/^(.*)$/gim, '<p>$1</p>')
        .replace(/<p><\/p>/g, '')
        .replace(/<p>(<h[1-6]>)/g, '$1')
        .replace(/(<\/h[1-6]>)<\/p>/g, '$1')
        .replace(/<p>(<ul>)/g, '$1')
        .replace(/(<\/ul>)<\/p>/g, '$1');
}

// Review prompt with dual analysis
async function reviewPrompt() {
    const prompt = document.getElementById('main-prompt').value.trim();
    if (!prompt) {
        alert('Please enter a prompt first!');
        return;
    }

    const btn = document.getElementById('review-btn-text');
    const originalText = btn.textContent;
    
    btn.innerHTML = '<span class="spinner"></span>Analyzing...';
    
    try {
        const response = await fetch(`${AppState.API_BASE}/dual-critique`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                prompt,
                model: getCurrentModel()
            })
        });
        
        const data = await response.json();
        
        // Switch to review tab and show results
        switchTabProgrammatically('review');
        const content = document.getElementById('result-content');
        
        if (data.success) {
            content.className = 'results-content success';
            
            // Create dual report display
            let quickReport = data.data.quick_report;
            let detailedReport = data.data.detailed_report;
            
            // Check if reports contain HTML tags
            if (!quickReport.includes('<') || !quickReport.includes('>')) {
                quickReport = markdownToHtml(quickReport);
            }
            if (!detailedReport.includes('<') || !detailedReport.includes('>')) {
                detailedReport = markdownToHtml(detailedReport);
            }
            
            // Create the dual report interface
            content.innerHTML = `
                <div class="dual-analysis-container">
                    <div class="analysis-tabs">
                        <button class="analysis-tab active" onclick="switchAnalysisTab('quick')">🚀 Quick Analysis</button>
                        <button class="analysis-tab" onclick="switchAnalysisTab('detailed')">📋 Detailed Analysis</button>
                    </div>
                    <div class="analysis-content">
                        <div id="quick-analysis" class="analysis-panel active">
                            ${quickReport}
                        </div>
                        <div id="detailed-analysis" class="analysis-panel" style="display: none;">
                            ${detailedReport}
                        </div>
                    </div>
                </div>
            `;
        } else {
            content.className = 'results-content error';
            content.innerHTML = `<p><strong>Error:</strong> ${data.error}</p>`;
        }
    } catch (error) {
        switchTabProgrammatically('review');
        const content = document.getElementById('result-content');
        content.className = 'results-content error';
        content.innerHTML = `<p><strong>Network Error:</strong> ${error.message}</p>`;
    } finally {
        btn.textContent = originalText;
    }
}

// Switch between quick and detailed analysis tabs
function switchAnalysisTab(tabName) {
    const tabs = document.querySelectorAll('.analysis-tab');
    const panels = document.querySelectorAll('.analysis-panel');
    
    tabs.forEach(tab => tab.classList.remove('active'));
    panels.forEach(panel => {
        panel.classList.remove('active');
        panel.style.display = 'none';
    });
    
    document.querySelector(`.analysis-tab[onclick*="${tabName}"]`).classList.add('active');
    const targetPanel = document.getElementById(`${tabName}-analysis`);
    targetPanel.classList.add('active');
    targetPanel.style.display = 'block';
}

// Variables helper functions
function extractVariables(prompt) {
    const variableRegex = /\{\{([^}]+)\}\}/g;
    const variables = [];
    let match;
    
    while ((match = variableRegex.exec(prompt)) !== null) {
        const varName = match[1].trim();
        if (!variables.includes(varName)) {
            variables.push(varName);
        }
    }
    
    return variables;
}

function substituteVariables(prompt) {
    const variableInputs = document.querySelectorAll('.variable-input');
    let processedPrompt = prompt;
    
    variableInputs.forEach(input => {
        const varName = input.dataset.variable;
        const value = input.value.trim();
        const placeholder = new RegExp(`\\{\\{\\s*${varName}\\s*\\}\\}`, 'g');
        
        if (value) {
            processedPrompt = processedPrompt.replace(placeholder, value);
        }
    });
    
    return processedPrompt;
}

function validateVariables() {
    const variableInputs = document.querySelectorAll('.variable-input');
    const emptyVariables = [];
    
    variableInputs.forEach(input => {
        if (!input.value.trim()) {
            emptyVariables.push(input.dataset.variable);
        }
    });
    
    return {
        isValid: emptyVariables.length === 0,
        emptyVariables: emptyVariables
    };
}

// Test Prompt execution with variables and parameters
async function executeTest() {
    const prompt = document.getElementById('main-prompt').value.trim();
    if (!prompt) {
        alert('Please enter a prompt first!');
        return;
    }
    
    // Check if there are variables and validate them
    const variables = extractVariables(prompt);
    let processedPrompt = prompt;
    
    if (variables.length > 0) {
        const validation = validateVariables();
        if (!validation.isValid) {
            alert(`Please fill in values for: ${validation.emptyVariables.join(', ')}`);
            return;
        }
        processedPrompt = substituteVariables(prompt);
    }
    
    // Get parameters
    const temperature = parseFloat(document.getElementById('test-temperature').value);
    const maxTokens = parseInt(document.getElementById('test-max-tokens').value) || 1000;
    const model = getCurrentModel();
    
    const btn = document.getElementById('test-btn-text');
    const originalText = btn.textContent;
    
    btn.innerHTML = '<span class="spinner"></span>Executing...';
    
    try {
        const response = await fetch(`${AppState.API_BASE}/execute`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                prompt: processedPrompt, 
                model: model,
                temperature, 
                max_tokens: maxTokens 
            })
        });
        
        const data = await response.json();
        
        // Save history to database
        const historyPrompt = variables.length > 0 ? `[Variables] ${prompt}` : prompt;
        await saveToHistory({
            prompt: historyPrompt,
            model: model,
            temperature,
            max_tokens: maxTokens,
            success: data.success,
            response: data.success ? data.data : "",
            error_msg: data.success ? "" : data.error
        });
        
        // Switch to execution tab and show results
        switchTabProgrammatically('execution');
        const content = document.getElementById('result-content');
        
        if (data.success) {
            content.className = 'results-content success';
            let resultHtml = `
                <h3>Test Execution Results</h3>
                <p><strong>Model:</strong> ${model} | <strong>Temperature:</strong> ${temperature} | <strong>Max Tokens:</strong> ${maxTokens}</p>
            `;
            
            if (variables.length > 0) {
                resultHtml += `
                    <div style="margin: 12px 0; padding: 8px; background: #2d2d30; border-radius: 4px;">
                        <strong>Variables Used:</strong>
                        <ul style="margin: 8px 0; padding-left: 20px;">
                `;
                variables.forEach(variable => {
                    const input = document.getElementById(`var-${variable}`);
                    const value = input ? input.value : '';
                    resultHtml += `<li style="margin: 4px 0; font-size: 11px;">{{${variable}}} → "${value}"</li>`;
                });
                resultHtml += `
                        </ul>
                        <strong>Processed Prompt:</strong>
                        <pre style="margin: 8px 0; white-space: pre-wrap; font-size: 11px;">${processedPrompt}</pre>
                    </div>
                `;
            }
            
            resultHtml += `
                <hr style="margin: 12px 0; border: 1px solid #3e3e42;">
                <div style="white-space: pre-wrap;">${data.data}</div>
            `;
            
            content.innerHTML = resultHtml;
        } else {
            content.className = 'results-content error';
            content.innerHTML = `<p><strong>Error:</strong> ${data.error}</p>`;
        }
    } catch (error) {
        switchTabProgrammatically('execution');
        const content = document.getElementById('result-content');
        content.className = 'results-content error';
        content.innerHTML = `<p><strong>Network Error:</strong> ${error.message}</p>`;
    } finally {
        btn.textContent = originalText;
    }
}

// Save history to database
async function saveToHistory(historyData) {
    try {
        await fetch(`${AppState.API_BASE}/history`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(historyData)
        });
    } catch (error) {
        console.error('Failed to save history:', error);
        // Fallback to local storage if database save fails
        const historyItem = {
            timestamp: new Date().toLocaleString(),
            prompt: historyData.prompt.substring(0, 100) + (historyData.prompt.length > 100 ? '...' : ''),
            temperature: historyData.temperature,
            maxTokens: historyData.max_tokens,
            success: historyData.success,
            response: historyData.success ? historyData.response : historyData.error_msg
        };
        AppState.executionHistory.unshift(historyItem);
    }
}

// Load history from database
async function loadHistoryFromDB() {
    try {
        const response = await fetch(`${AppState.API_BASE}/history`);
        const data = await response.json();
        
        if (data.success) {
            return data.data;
        } else {
            console.error('Failed to load history:', data.error);
            return [];
        }
    } catch (error) {
        console.error('Failed to load history:', error);
        return [];
    }
}

// Clear history in database
async function clearHistoryDB() {
    try {
        const response = await fetch(`${AppState.API_BASE}/history`, {
            method: 'DELETE'
        });
        const data = await response.json();
        
        if (data.success) {
            console.log('History cleared successfully');
        } else {
            console.error('Failed to clear history:', data.error);
        }
    } catch (error) {
        console.error('Failed to clear history:', error);
    }
}

// Make functions globally accessible
window.executeTest = executeTest;
window.reviewPrompt = reviewPrompt;
window.saveToHistory = saveToHistory; 