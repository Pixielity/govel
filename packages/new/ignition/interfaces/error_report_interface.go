package interfaces

import "time"

// ErrorReportInterface interface for structured error reports
type ErrorReportInterface interface {
	GetMessage() string
	SetMessage(string)
	GetType() string
	SetType(string)
	GetFile() string
	SetFile(string)
	GetLine() int
	SetLine(int)
	GetStack() []StackFrameInterface
	SetStack([]StackFrameInterface)
	AddStackFrame(StackFrameInterface)
	GetContext() ErrorContextInterface
	SetContext(ErrorContextInterface)
	GetSolutions() []SolutionInterface
	SetSolutions([]SolutionInterface)
	AddSolution(SolutionInterface)
	GetTimestamp() time.Time
	SetTimestamp(time.Time)
	IsEmpty() bool
	HasSolutions() bool
	GetStackFrameCount() int
}
