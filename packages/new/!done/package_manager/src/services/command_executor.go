package services

import (
	"context"
	"fmt"
	"govel/package_manager/interfaces"
	"govel/package_manager/models"
	"os"
	"os/exec"
	"strings"
	"time"
)

// CommandExecutor implements ExecutorInterface for running commands and scripts
type CommandExecutor struct{}

// NewCommandExecutor creates a new command executor instance
func NewCommandExecutor() interfaces.ExecutorInterface {
	return &CommandExecutor{}
}

// ExecuteCommand executes a single command in the specified working directory
func (ce *CommandExecutor) ExecuteCommand(ctx context.Context, command string, workDir string) (*models.CommandResult, error) {
	startTime := time.Now()

	result := &models.CommandResult{
		Command:  command,
		Success:  false,
		Output:   "",
		Error:    "",
		ExitCode: 0,
		Duration: 0,
	}

	// Parse the command
	parts := ce.parseCommand(command)
	if len(parts) == 0 {
		result.Error = "empty command"
		result.Duration = time.Since(startTime)
		return result, fmt.Errorf("empty command")
	}

	// Create the command
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Dir = workDir

	// Set environment variables
	cmd.Env = os.Environ()

	// Capture output
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Error = err.Error()
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1
		}
		return result, nil // Don't return error, just capture it in result
	}

	result.Success = true
	result.ExitCode = 0

	return result, nil
}

// ExecuteScript executes a script (potentially multi-line) in the specified working directory
func (ce *CommandExecutor) ExecuteScript(ctx context.Context, script string, workDir string) (*models.CommandResult, error) {
	startTime := time.Now()

	result := &models.CommandResult{
		Command:  script,
		Success:  false,
		Output:   "",
		Error:    "",
		ExitCode: 0,
		Duration: 0,
	}

	// Handle multi-line scripts or scripts with &&
	commands := ce.parseScript(script)
	var outputs []string
	var errors []string

	for _, command := range commands {
		command = strings.TrimSpace(command)
		if command == "" {
			continue
		}

		cmdResult, err := ce.ExecuteCommand(ctx, command, workDir)
		if err != nil {
			result.Error = err.Error()
			result.Duration = time.Since(startTime)
			return result, err
		}

		outputs = append(outputs, cmdResult.Output)
		if cmdResult.Error != "" {
			errors = append(errors, cmdResult.Error)
		}

		// If a command fails, stop execution
		if !cmdResult.Success {
			result.ExitCode = cmdResult.ExitCode
			result.Output = strings.Join(outputs, "\n")
			result.Error = strings.Join(errors, "\n")
			result.Duration = time.Since(startTime)
			return result, nil
		}
	}

	result.Success = true
	result.ExitCode = 0
	result.Output = strings.Join(outputs, "\n")
	if len(errors) > 0 {
		result.Error = strings.Join(errors, "\n")
	}
	result.Duration = time.Since(startTime)

	return result, nil
}

// ExecuteHooks executes a series of hook commands
func (ce *CommandExecutor) ExecuteHooks(ctx context.Context, hooks []string, workDir string) ([]*models.CommandResult, error) {
	var results []*models.CommandResult

	for _, hook := range hooks {
		if strings.TrimSpace(hook) == "" {
			continue
		}

		result, err := ce.ExecuteScript(ctx, hook, workDir)
		if err != nil {
			return results, fmt.Errorf("hook failed: %w", err)
		}

		results = append(results, result)

		// If a hook fails, stop execution
		if !result.Success {
			return results, fmt.Errorf("hook failed with exit code %d: %s", result.ExitCode, result.Error)
		}
	}

	return results, nil
}

// parseCommand parses a command string into parts, handling quotes
func (ce *CommandExecutor) parseCommand(command string) []string {
	var parts []string
	var current strings.Builder
	var inQuotes bool
	var quoteChar rune

	for i, char := range command {
		switch char {
		case '"', '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else {
				current.WriteRune(char)
			}
		case ' ', '\t':
			if inQuotes {
				current.WriteRune(char)
			} else if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}

		// Handle last part
		if i == len(command)-1 && current.Len() > 0 {
			parts = append(parts, current.String())
		}
	}

	return parts
}

// parseScript parses a script into individual commands
func (ce *CommandExecutor) parseScript(script string) []string {
	// Handle scripts with && operators
	if strings.Contains(script, "&&") {
		return strings.Split(script, "&&")
	}

	// Handle scripts with newlines
	if strings.Contains(script, "\n") {
		lines := strings.Split(script, "\n")
		var commands []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				commands = append(commands, line)
			}
		}
		return commands
	}

	// Single command
	return []string{script}
}

// ExecuteCommandWithStreaming executes a command with real-time output streaming
func (ce *CommandExecutor) ExecuteCommandWithStreaming(ctx context.Context, command string, workDir string) (*models.CommandResult, error) {
	startTime := time.Now()

	result := &models.CommandResult{
		Command:  command,
		Success:  false,
		Output:   "",
		Error:    "",
		ExitCode: 0,
		Duration: 0,
	}

	// Parse the command
	parts := ce.parseCommand(command)
	if len(parts) == 0 {
		result.Error = "empty command"
		result.Duration = time.Since(startTime)
		return result, fmt.Errorf("empty command")
	}

	// Create the command
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Dir = workDir
	cmd.Env = os.Environ()

	// Set up streaming output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err := cmd.Run()
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Error = err.Error()
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1
		}
		return result, nil
	}

	result.Success = true
	result.ExitCode = 0
	result.Output = "Command executed successfully (streamed output)"

	return result, nil
}

// ValidateCommand checks if a command is safe to execute
func (ce *CommandExecutor) ValidateCommand(command string) error {
	// Basic validation - can be extended with more security checks
	command = strings.TrimSpace(command)

	if command == "" {
		return fmt.Errorf("empty command")
	}

	// Check for potentially dangerous commands
	dangerousCommands := []string{
		"rm -rf /",
		"sudo rm",
		"format",
		"del /s",
		"rmdir /s",
	}

	lowerCommand := strings.ToLower(command)
	for _, dangerous := range dangerousCommands {
		if strings.Contains(lowerCommand, strings.ToLower(dangerous)) {
			return fmt.Errorf("potentially dangerous command detected: %s", command)
		}
	}

	return nil
}

// GetWorkingDirectory returns the current working directory or creates it
func (ce *CommandExecutor) GetWorkingDirectory(path string) (string, error) {
	if path == "" {
		return os.Getwd()
	}

	// Ensure the directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("failed to create working directory: %w", err)
	}

	return path, nil
}
