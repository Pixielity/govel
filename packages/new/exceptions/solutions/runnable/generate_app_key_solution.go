package runnable

import (
	"fmt"
	"os"
	"os/exec"

	"govel/packages/exceptions/core/solution"
	solutionInterface "govel/packages/exceptions/interfaces/solution"
)

// GenerateAppKeySolution provides a runnable solution for missing application keys
type GenerateAppKeySolution struct {
	*solution.BaseSolution
}

// NewGenerateAppKeySolution creates a new runnable solution for generating app keys
func NewGenerateAppKeySolution() *GenerateAppKeySolution {
	base := solution.NewBaseSolution("Generate Application Key").
		SetSolutionDescription("Your application needs a unique encryption key for security. This solution will generate a new APP_KEY for your .env file.").
		AddDocumentationLink("GoVel Security Docs", "https://govel.dev/docs/security").
		AddDocumentationLink("Environment Configuration", "https://govel.dev/docs/configuration")

	return &GenerateAppKeySolution{
		BaseSolution: base,
	}
}

// GetSolutionActionDescription returns what the runnable action will do
func (s *GenerateAppKeySolution) GetSolutionActionDescription() string {
	return "Generate a new APP_KEY and add it to your .env file"
}

// GetRunButtonText returns the text for the run button
func (s *GenerateAppKeySolution) GetRunButtonText() string {
	return "Generate App Key"
}

// Run executes the solution to generate an app key
func (s *GenerateAppKeySolution) Run(parameters map[string]interface{}) error {
	// Generate a random 32-byte key (base64 encoded)
	cmd := exec.Command("openssl", "rand", "-base64", "32")
	key, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to generate app key: %w", err)
	}

	keyString := string(key)
	keyString = keyString[:len(keyString)-1] // Remove trailing newline

	// Find .env file
	envPath := ".env"
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Try .env.example
		if _, err := os.Stat(".env.example"); err == nil {
			// Copy .env.example to .env
			input, err := os.ReadFile(".env.example")
			if err != nil {
				return fmt.Errorf("failed to read .env.example: %w", err)
			}
			err = os.WriteFile(envPath, input, 0644)
			if err != nil {
				return fmt.Errorf("failed to create .env file: %w", err)
			}
		} else {
			// Create basic .env file
			content := fmt.Sprintf("APP_KEY=%s\n", keyString)
			err = os.WriteFile(envPath, []byte(content), 0644)
			if err != nil {
				return fmt.Errorf("failed to create .env file: %w", err)
			}
			return nil
		}
	}

	// Read existing .env file
	content, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	// Update or add APP_KEY
	contentStr := string(content)
	if len(contentStr) > 0 && contentStr[len(contentStr)-1] != '\n' {
		contentStr += "\n"
	}
	contentStr += fmt.Sprintf("APP_KEY=%s\n", keyString)

	// Write back to .env file
	err = os.WriteFile(envPath, []byte(contentStr), 0644)
	if err != nil {
		return fmt.Errorf("failed to update .env file: %w", err)
	}

	return nil
}

// GetRunParameters returns the parameters this solution expects
func (s *GenerateAppKeySolution) GetRunParameters() map[string]interface{} {
	return map[string]interface{}{} // No parameters needed
}

// Ensure GenerateAppKeySolution implements the RunnableSolution interface
var _ solutionInterface.RunnableSolution = (*GenerateAppKeySolution)(nil)
