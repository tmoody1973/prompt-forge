package config

import "fmt"

// Azure OpenAI Configuration
const (
	AZURE_BASE_URL = ""
	AZURE_API_KEY  = ""
	API_VERSION    = ""
)

// Model deployment mappings (only for available deployments)
var ModelDeployments = map[string]string{
	"gpt-4.1": "gpt-4.1",
	"o3":      "o3",
}

// GetEndpointURL builds the complete endpoint URL for a given model
func GetEndpointURL(model string) string {
	deployment, exists := ModelDeployments[model]
	if !exists {
		deployment = "gpt-4.1" // fallback to default
	}
	return fmt.Sprintf("%s/%s/chat/completions?api-version=%s", AZURE_BASE_URL, deployment, API_VERSION)
}
