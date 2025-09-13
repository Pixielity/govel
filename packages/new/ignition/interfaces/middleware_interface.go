package interfaces

// MiddlewareInterface interface for processing error reports
type MiddlewareInterface interface {
	// Process takes an error report and returns a modified error report
	// This allows middleware to add, modify, or filter error report data
	Process(ErrorReportInterface) ErrorReportInterface
}
