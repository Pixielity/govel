package runnable

import (
	"fmt"
	"os"
	"os/exec"

	"govel/packages/exceptions/core/solution"
	solutionInterface "govel/packages/exceptions/interfaces/solution"
)

// CreateDirectorySolution provides a runnable solution for creating missing directories
type CreateDirectorySolution struct {
	*solution.BaseSolution
	directoryPath string
}

// NewCreateDirectorySolution creates a new runnable solution for creating directories
func NewCreateDirectorySolution(directoryPath string) *CreateDirectorySolution {
	base := solution.NewBaseSolution("Create Missing Directory").
		SetSolutionDescription(fmt.Sprintf("The directory '%s' is required but doesn't exist. This solution will create it with proper permissions.", directoryPath)).
		AddDocumentationLink("GoVel File System Docs", "https://govel.dev/docs/filesystem")

	return &CreateDirectorySolution{
		BaseSolution:  base,
		directoryPath: directoryPath,
	}
}

// GetSolutionActionDescription returns what the runnable action will do
func (s *CreateDirectorySolution) GetSolutionActionDescription() string {
	return fmt.Sprintf("Create the directory '%s' with proper permissions", s.directoryPath)
}

// GetRunButtonText returns the text for the run button
func (s *CreateDirectorySolution) GetRunButtonText() string {
	return "Create Directory"
}

// Run executes the solution to create the directory
func (s *CreateDirectorySolution) Run(parameters map[string]interface{}) error {
	err := os.MkdirAll(s.directoryPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", s.directoryPath, err)
	}
	return nil
}

// GetRunParameters returns the parameters this solution expects
func (s *CreateDirectorySolution) GetRunParameters() map[string]interface{} {
	return map[string]interface{}{
		"directory_path": s.directoryPath,
	}
}

// InstallDependencySolution provides a runnable solution for installing missing dependencies
type InstallDependencySolution struct {
	*solution.BaseSolution
	packageName string
	command     string
}

// NewInstallDependencySolution creates a new runnable solution for installing dependencies
func NewInstallDependencySolution(packageName, command string) *InstallDependencySolution {
	base := solution.NewBaseSolution("Install Missing Dependency").
		SetSolutionDescription(fmt.Sprintf("The package '%s' is required but not installed. This solution will install it using '%s'.", packageName, command)).
		AddDocumentationLink("GoVel Dependencies Docs", "https://govel.dev/docs/dependencies")

	return &InstallDependencySolution{
		BaseSolution: base,
		packageName:  packageName,
		command:      command,
	}
}

// GetSolutionActionDescription returns what the runnable action will do
func (s *InstallDependencySolution) GetSolutionActionDescription() string {
	return fmt.Sprintf("Run '%s' to install the missing dependency", s.command)
}

// GetRunButtonText returns the text for the run button
func (s *InstallDependencySolution) GetRunButtonText() string {
	return "Install Dependency"
}

// Run executes the solution to install the dependency
func (s *InstallDependencySolution) Run(parameters map[string]interface{}) error {
	// Parse command into parts
	parts := []string{"sh", "-c", s.command}
	
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = "." // Run in current directory
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install dependency '%s': %s\nOutput: %s", s.packageName, err.Error(), string(output))
	}
	
	return nil
}

// GetRunParameters returns the parameters this solution expects
func (s *InstallDependencySolution) GetRunParameters() map[string]interface{} {
	return map[string]interface{}{
		"package_name": s.packageName,
		"command":      s.command,
	}
}

// FixPermissionsSolution provides a runnable solution for fixing file permissions
type FixPermissionsSolution struct {
	*solution.BaseSolution
	path        string
	permissions os.FileMode
}

// NewFixPermissionsSolution creates a new runnable solution for fixing permissions
func NewFixPermissionsSolution(path string, permissions os.FileMode) *FixPermissionsSolution {
	base := solution.NewBaseSolution("Fix File Permissions").
		SetSolutionDescription(fmt.Sprintf("The file or directory '%s' has incorrect permissions. This solution will set the permissions to %o.", path, permissions)).
		AddDocumentationLink("GoVel File System Docs", "https://govel.dev/docs/filesystem")

	return &FixPermissionsSolution{
		BaseSolution: base,
		path:         path,
		permissions:  permissions,
	}
}

// GetSolutionActionDescription returns what the runnable action will do
func (s *FixPermissionsSolution) GetSolutionActionDescription() string {
	return fmt.Sprintf("Set permissions of '%s' to %o", s.path, s.permissions)
}

// GetRunButtonText returns the text for the run button
func (s *FixPermissionsSolution) GetRunButtonText() string {
	return "Fix Permissions"
}

// Run executes the solution to fix permissions
func (s *FixPermissionsSolution) Run(parameters map[string]interface{}) error {
	err := os.Chmod(s.path, s.permissions)
	if err != nil {
		return fmt.Errorf("failed to change permissions of '%s': %w", s.path, err)
	}
	return nil
}

// GetRunParameters returns the parameters this solution expects
func (s *FixPermissionsSolution) GetRunParameters() map[string]interface{} {
	return map[string]interface{}{
		"path":        s.path,
		"permissions": fmt.Sprintf("%o", s.permissions),
	}
}

// Ensure all runnable solutions implement the RunnableSolution interface
var _ solutionInterface.RunnableSolution = (*CreateDirectorySolution)(nil)
var _ solutionInterface.RunnableSolution = (*InstallDependencySolution)(nil)
var _ solutionInterface.RunnableSolution = (*FixPermissionsSolution)(nil)
