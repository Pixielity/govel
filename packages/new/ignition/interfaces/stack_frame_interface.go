package interfaces

// StackFrameInterface interface for single frames in the stack trace
type StackFrameInterface interface {
	GetFunction() string
	SetFunction(string)
	GetFile() string
	SetFile(string)
	GetLine() int
	SetLine(int)
	GetCode() map[string]string
	SetCode(map[string]string)
	AddCodeLine(string, string)
	IsApplicationFrame() bool
	SetApplicationFrame(bool)
	GetRelativeFile() string
	GetFileName() string
	GetPackageName() string
	GetMethodName() string
	IsGoStdLib() bool
	IsThirdParty() bool
	HasCode() bool
	GetCodeLineCount() int
	ToString() string
}
