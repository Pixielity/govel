package interfaces

// RequestContextInterface interface for HTTP request information
type RequestContextInterface interface {
	GetMethod() string
	SetMethod(string)
	GetURL() string
	SetURL(string)
	GetHeaders() map[string]string
	SetHeaders(map[string]string)
	GetHeader(string) string
	SetHeader(string, string)
	GetBody() string
	SetBody(string)
	HasHeader(string) bool
	GetHeaderCount() int
	HasBody() bool
	IsGET() bool
	IsPOST() bool
	IsPUT() bool
	IsDELETE() bool
	IsAjax() bool
	IsSecure() bool
}
