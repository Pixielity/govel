package interfaces

// StackTraceBuilderInterface interface for building stack traces
type StackTraceBuilderInterface interface {
	// BuildStackTrace builds a stack trace from the given error
	BuildStackTrace(error) []StackFrameInterface
	
	// SetContextLines sets the number of context lines to include around each frame
	SetContextLines(int)
	
	// GetContextLines returns the number of context lines
	GetContextLines() int
	
	// ShouldSkipFrame determines if a frame should be skipped
	ShouldSkipFrame(string) bool
}
