package handlers

import (
	"bufio"
	"os"
	"strconv"
)

// SourceCodeExtractor handles extracting source code around error lines
type SourceCodeExtractor struct {
	contextLines int
}

// NewSourceCodeExtractor creates a new source code extractor
func NewSourceCodeExtractor(contextLines int) *SourceCodeExtractor {
	return &SourceCodeExtractor{
		contextLines: contextLines,
	}
}

// ExtractSourceCode extracts source code around the specified line
func (s *SourceCodeExtractor) ExtractSourceCode(filename string, targetLine int) map[string]string {
	codeSnippet := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return codeSnippet
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	startLine := max(1, targetLine-s.contextLines)
	endLine := targetLine + s.contextLines

	for scanner.Scan() {
		currentLine++
		
		if currentLine >= startLine && currentLine <= endLine {
			lineKey := strconv.Itoa(currentLine)
			codeSnippet[lineKey] = scanner.Text()
		}
		
		if currentLine > endLine {
			break
		}
	}

	return codeSnippet
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
