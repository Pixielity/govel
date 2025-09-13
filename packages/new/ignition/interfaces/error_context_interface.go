package interfaces

// ErrorContextInterface interface for contextual information about the error
type ErrorContextInterface interface {
	GetRequest() RequestContextInterface
	SetRequest(RequestContextInterface)
	GetEnvironment() EnvContextInterface
	SetEnvironment(EnvContextInterface)
	GetUser() interface{}
	SetUser(interface{})
	HasRequest() bool
	HasEnvironment() bool
	HasUser() bool
	IsEmpty() bool
}
