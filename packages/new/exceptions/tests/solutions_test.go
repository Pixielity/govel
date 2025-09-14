package tests

import (
	"fmt"
	"testing"

	"govel/exceptions"
	"govel/exceptions/solutions"
)

// Test BaseSolution implementation
func TestBaseSolution(t *testing.T) {
	title := "Test Solution"
	description := "This is a test solution"
	links := map[string]string{
		"Documentation": "https://example.com/docs",
		"Tutorial":      "https://example.com/tutorial",
	}

	solution := solutions.NewBaseSolution(title).
		SetSolutionDescription(description).
		SetDocumentationLinks(links)

	if solution.GetSolutionTitle() != title {
		t.Errorf("Expected title '%s', got '%s'", title, solution.GetSolutionTitle())
	}

	if solution.GetSolutionDescription() != description {
		t.Errorf("Expected description '%s', got '%s'", description, solution.GetSolutionDescription())
	}

	solutionLinks := solution.GetDocumentationLinks()
	if len(solutionLinks) != len(links) {
		t.Errorf("Expected %d links, got %d", len(links), len(solutionLinks))
	}

	for key, value := range links {
		if solutionLinks[key] != value {
			t.Errorf("Expected link '%s' = '%s', got '%s'", key, value, solutionLinks[key])
		}
	}
}

// Test method chaining for BaseSolution
func TestBaseSolutionChaining(t *testing.T) {
	solution := solutions.NewBaseSolution("Test").
		SetSolutionDescription("Description").
		AddDocumentationLink("Link1", "URL1").
		AddDocumentationLink("Link2", "URL2")

	if solution.GetSolutionTitle() != "Test" {
		t.Error("Method chaining failed for title")
	}

	if solution.GetSolutionDescription() != "Description" {
		t.Error("Method chaining failed for description")
	}

	links := solution.GetDocumentationLinks()
	if len(links) != 2 {
		t.Errorf("Expected 2 links, got %d", len(links))
	}

	if links["Link1"] != "URL1" {
		t.Error("Method chaining failed for Link1")
	}

	if links["Link2"] != "URL2" {
		t.Error("Method chaining failed for Link2")
	}
}

// Test SolutionProviderRepository
func TestSolutionProviderRepository(t *testing.T) {
	repo := solutions.NewSolutionProviderRepository()

	// Test with no providers
	err := fmt.Errorf("test error")
	solutionsFound := repo.GetSolutionsForError(err)
	if len(solutionsFound) != 0 {
		t.Errorf("Expected 0 solutions with no providers, got %d", len(solutionsFound))
	}

	// Register a provider
	provider := &TestSolutionProvider{}
	repo.RegisterSolutionProvider(provider)

	providers := repo.GetProviders()
	if len(providers) != 1 {
		t.Errorf("Expected 1 provider, got %d", len(providers))
	}

	// Test with provider that can solve
	testErr := fmt.Errorf("test error that can be solved")
	solutionsFound = repo.GetSolutionsForError(testErr)
	if len(solutionsFound) != 1 {
		t.Errorf("Expected 1 solution, got %d", len(solutionsFound))
	}

	if solutionsFound[0].GetSolutionTitle() != "Test Solution" {
		t.Errorf("Expected 'Test Solution', got '%s'", solutionsFound[0].GetSolutionTitle())
	}

	// Test with provider that cannot solve
	cannotSolveErr := fmt.Errorf("different error")
	solutionsFound = repo.GetSolutionsForError(cannotSolveErr)
	if len(solutionsFound) != 0 {
		t.Errorf("Expected 0 solutions for unsolvable error, got %d", len(solutionsFound))
	}
}

// Test multiple solution providers
func TestMultipleSolutionProviders(t *testing.T) {
	repo := solutions.NewSolutionProviderRepository()

	provider1 := &TestSolutionProvider{}
	provider2 := &AnotherTestSolutionProvider{}

	repo.RegisterSolutionProviders([]solutions.HasSolutionsForThrowable{provider1, provider2})

	providers := repo.GetProviders()
	if len(providers) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(providers))
	}

	// Test error that both can solve
	err := fmt.Errorf("test error that can be solved")
	solutionsFound := repo.GetSolutionsForError(err)
	if len(solutionsFound) != 2 {
		t.Errorf("Expected 2 solutions, got %d", len(solutionsFound))
	}
}

// Test exception with solution
func TestExceptionWithSolution(t *testing.T) {
	exc := exceptions.NewException("Test error", 400)

	// Initially no solution
	if exc.HasSolution() {
		t.Error("Exception should not have solution initially")
	}

	if exc.GetSolution() != nil {
		t.Error("GetSolution should return nil initially")
	}

	// Set solution
	solution := solutions.NewBaseSolution("Test Solution")
	exc.SetSolution(solution)

	if !exc.HasSolution() {
		t.Error("Exception should have solution after setting")
	}

	if exc.GetSolution() == nil {
		t.Error("GetSolution should not return nil after setting")
	}

	if exc.GetSolution().GetSolutionTitle() != "Test Solution" {
		t.Errorf("Expected 'Test Solution', got '%s'", exc.GetSolution().GetSolutionTitle())
	}
}

// Test exception rendering with solutions
func TestExceptionRenderingWithSolutions(t *testing.T) {
	exc := exceptions.NewException("Test error", 400)
	solution := solutions.NewBaseSolution("Test Solution").
		SetSolutionDescription("Test description").
		AddDocumentationLink("Test Link", "https://example.com")

	exc.SetSolution(solution)

	rendered := exc.Render()

	// Check that solution is included in rendered output
	if rendered["solution"] == nil {
		t.Error("Rendered exception should include solution")
	}

	solutionData, ok := rendered["solution"].(map[string]interface{})
	if !ok {
		t.Error("Solution should be a map")
	}

	if solutionData["title"] != "Test Solution" {
		t.Errorf("Expected solution title 'Test Solution', got '%v'", solutionData["title"])
	}

	if solutionData["description"] != "Test description" {
		t.Errorf("Expected solution description 'Test description', got '%v'", solutionData["description"])
	}

	if solutionData["runnable"] != false {
		t.Error("Basic solution should not be runnable")
	}
}

// Test exception with runnable solution
func TestExceptionWithRunnableSolution(t *testing.T) {
	exc := exceptions.NewException("Test error", 500)
	runnableSolution := &TestRunnableSolution{
		BaseSolution: solutions.NewBaseSolution("Runnable Test Solution"),
	}

	exc.SetSolution(runnableSolution)

	rendered := exc.Render()
	solutionData := rendered["solution"].(map[string]interface{})

	if solutionData["runnable"] != true {
		t.Error("Runnable solution should be marked as runnable")
	}

	if solutionData["action_description"] != "Test action description" {
		t.Errorf("Expected action description, got '%v'", solutionData["action_description"])
	}

	if solutionData["run_button_text"] != "Test Button" {
		t.Errorf("Expected button text, got '%v'", solutionData["run_button_text"])
	}
}

// Test HTTP exceptions with built-in solutions
func TestHTTPExceptionsWithSolutions(t *testing.T) {
	testCases := []struct {
		exception     exceptions.ExceptionInterface
		expectedTitle string
	}{
		{exceptions.NewNotFoundException("Not found"), "Resource Not Found"},
		{exceptions.NewUnauthorizedException("Unauthorized"), "Authentication Required"},
		{exceptions.NewForbiddenException("Forbidden"), "Access Forbidden"},
		{exceptions.NewBadRequestException("Bad request"), "Bad Request"},
		{exceptions.NewInternalServerErrorException("Server error"), "Internal Server Error"},
		{exceptions.NewServiceUnavailableException("Unavailable", 60), "Service Temporarily Unavailable"},
		{exceptions.NewTooManyRequestsException("Rate limited", 30), "Rate Limit Exceeded"},
	}

	for _, tc := range testCases {
		if !tc.exception.HasSolution() {
			t.Errorf("Exception '%s' should have a solution", tc.exception.GetMessage())
			continue
		}

		solution := tc.exception.GetSolution()
		if solution.GetSolutionTitle() != tc.expectedTitle {
			t.Errorf("Expected solution title '%s', got '%s'", tc.expectedTitle, solution.GetSolutionTitle())
		}

		// Check that description is not empty
		if solution.GetSolutionDescription() == "" {
			t.Errorf("Solution for '%s' should have description", tc.exception.GetMessage())
		}

		// Check that there are documentation links
		if len(solution.GetDocumentationLinks()) == 0 {
			t.Errorf("Solution for '%s' should have documentation links", tc.exception.GetMessage())
		}
	}
}

// Test HTTP exception provider
func TestHTTPExceptionProvider(t *testing.T) {
	provider := solutions.NewHTTPExceptionProvider()

	testCases := []struct {
		error        error
		canSolve     bool
		numSolutions int
	}{
		{fmt.Errorf("404 not found"), true, 1},
		{fmt.Errorf("Not Found error occurred"), true, 1},
		{fmt.Errorf("401 unauthorized access"), true, 1},
		{fmt.Errorf("Unauthorized request"), true, 1},
		{fmt.Errorf("500 internal server error"), true, 1},
		{fmt.Errorf("Internal server problem"), true, 1},
		{fmt.Errorf("Some random error"), false, 0},
		{fmt.Errorf("Database connection error"), false, 0},
	}

	for _, tc := range testCases {
		canSolve := provider.CanSolve(tc.error)
		if canSolve != tc.canSolve {
			t.Errorf("Provider CanSolve for '%s': expected %t, got %t", tc.error.Error(), tc.canSolve, canSolve)
		}

		if canSolve {
			solutionsFound := provider.GetSolutions(tc.error)
			if len(solutionsFound) != tc.numSolutions {
				t.Errorf("Provider solutions for '%s': expected %d, got %d", tc.error.Error(), tc.numSolutions, len(solutionsFound))
			}
		}
	}
}

// Test runnable solutions
func TestRunnableSolutions(t *testing.T) {
	testCases := []struct {
		solution      solutions.RunnableSolution
		expectedTitle string
	}{
		{solutions.NewGenerateAppKeySolution(), "Generate Application Key"},
		{solutions.NewCreateDirectorySolution("test/dir"), "Create Missing Directory"},
		{solutions.NewInstallDependencySolution("test-package", "npm install test-package"), "Install Missing Dependency"},
		{solutions.NewFixPermissionsSolution("test/file", 0644), "Fix File Permissions"},
	}

	for _, tc := range testCases {
		if tc.solution.GetSolutionTitle() != tc.expectedTitle {
			t.Errorf("Expected title '%s', got '%s'", tc.expectedTitle, tc.solution.GetSolutionTitle())
		}

		// Test that all runnable solutions have required methods
		if tc.solution.GetSolutionActionDescription() == "" {
			t.Errorf("Solution '%s' should have action description", tc.expectedTitle)
		}

		if tc.solution.GetRunButtonText() == "" {
			t.Errorf("Solution '%s' should have button text", tc.expectedTitle)
		}

		// Test run parameters
		params := tc.solution.GetRunParameters()
		if params == nil {
			t.Errorf("Solution '%s' should return non-nil parameters map", tc.expectedTitle)
		}
	}
}

// Test common runnable solutions provider
func TestCommonRunnableSolutionsProvider(t *testing.T) {
	provider := solutions.NewCommonRunnableSolutionsProvider()

	testCases := []struct {
		error          error
		canSolve       bool
		expectRunnable bool
	}{
		{fmt.Errorf("app key missing"), true, true},
		{fmt.Errorf("encryption key not found"), true, true},
		{fmt.Errorf("directory not found"), true, true},
		{fmt.Errorf("permission denied"), true, false}, // Can solve but doesn't provide runnable solutions for generic permission errors
		{fmt.Errorf("random error"), false, false},
	}

	for _, tc := range testCases {
		canSolve := provider.CanSolve(tc.error)
		if canSolve != tc.canSolve {
			t.Errorf("Provider CanSolve for '%s': expected %t, got %t", tc.error.Error(), tc.canSolve, canSolve)
		}

		if canSolve {
			solutionsFound := provider.GetSolutions(tc.error)
			// Some errors can be solved but don't have specific solutions implemented yet
			if tc.expectRunnable && len(solutionsFound) == 0 {
				t.Errorf("Provider should return solutions for '%s'", tc.error.Error())
			}

			// Check if solutions are runnable when expected
			for _, solution := range solutionsFound {
				_, isRunnable := solution.(solutions.RunnableSolution)
				if tc.expectRunnable && !isRunnable {
					t.Errorf("Expected runnable solution for '%s'", tc.error.Error())
				}
			}
		}
	}
}

// Test providers with exceptions that provide their own solutions
func TestExceptionProvidesSolution(t *testing.T) {
	repo := solutions.NewSolutionProviderRepository()
	provider := &TestSolutionProvider{}
	repo.RegisterSolutionProvider(provider)

	// Create an exception that provides its own solution
	exc := &TestExceptionWithOwnSolution{
		Exception: exceptions.NewException("Test error with own solution", 400),
	}

	solutionsFound := repo.GetSolutionsForError(exc)

	// Should get both the exception's own solution and provider solutions
	if len(solutionsFound) < 1 {
		t.Error("Should find at least the exception's own solution")
	}

	// The first solution should be the exception's own solution
	if solutionsFound[0].GetSolutionTitle() != "Custom Exception Solution" {
		t.Errorf("Expected 'Custom Exception Solution', got '%s'", solutionsFound[0].GetSolutionTitle())
	}
}

// Helper test types

type TestSolutionProvider struct{}

func (p *TestSolutionProvider) CanSolve(err error) bool {
	return err.Error() == "test error that can be solved"
}

func (p *TestSolutionProvider) GetSolutions(err error) []solutions.Solution {
	if p.CanSolve(err) {
		return []solutions.Solution{
			solutions.NewBaseSolution("Test Solution").
				SetSolutionDescription("This is a test solution"),
		}
	}
	return []solutions.Solution{}
}

type AnotherTestSolutionProvider struct{}

func (p *AnotherTestSolutionProvider) CanSolve(err error) bool {
	return err.Error() == "test error that can be solved"
}

func (p *AnotherTestSolutionProvider) GetSolutions(err error) []solutions.Solution {
	if p.CanSolve(err) {
		return []solutions.Solution{
			solutions.NewBaseSolution("Another Test Solution").
				SetSolutionDescription("This is another test solution"),
		}
	}
	return []solutions.Solution{}
}

type TestRunnableSolution struct {
	*solutions.BaseSolution
}

func (s *TestRunnableSolution) GetSolutionActionDescription() string {
	return "Test action description"
}

func (s *TestRunnableSolution) GetRunButtonText() string {
	return "Test Button"
}

func (s *TestRunnableSolution) Run(parameters map[string]interface{}) error {
	return nil // Mock implementation
}

func (s *TestRunnableSolution) GetRunParameters() map[string]interface{} {
	return map[string]interface{}{
		"test_param": "test_value",
	}
}

type TestExceptionWithOwnSolution struct {
	*exceptions.Exception
}

func (e *TestExceptionWithOwnSolution) GetSolution() solutions.Solution {
	return solutions.NewBaseSolution("Custom Exception Solution").
		SetSolutionDescription("This exception provides its own solution")
}
