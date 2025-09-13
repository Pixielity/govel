package interfaces

// SolutionInterface interface for error solutions compatible with Laravel Ignition
type SolutionInterface interface {
	GetClass() string
	SetClass(string)
	GetTitle() string
	SetTitle(string)
	GetDescription() string
	SetDescription(string)
	GetLinks() map[string]string
	SetLinks(map[string]string)
	AddLink(string, string)
	GetIsRunnable() bool
	SetIsRunnable(bool)
	GetAiGenerated() bool
	SetAiGenerated(bool)
	GetActionDescription() string
	SetActionDescription(string)
	GetRunButtonText() string
	SetRunButtonText(string)
	GetExecuteEndpoint() string
	SetExecuteEndpoint(string)
	GetRunParameters() []interface{}
	SetRunParameters([]interface{})
	HasLinks() bool
	HasRunCode() bool
	IsEmpty() bool
}
