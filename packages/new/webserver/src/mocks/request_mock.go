// Package mocks - Request mock
package mocks

import (
	"net/http"
	"net/url"
)

type RequestMock struct{}

func (RequestMock) Param(string) string { return "" }
func (RequestMock) ParamInt(string) int { return 0 }
func (RequestMock) Params() map[string]string { return map[string]string{} }
func (RequestMock) Query(string, ...string) string { return "" }
func (RequestMock) QueryInt(string, ...int) int { return 0 }
func (RequestMock) QueryBool(string, ...bool) bool { return false }
func (RequestMock) Queries() url.Values { return url.Values{} }
func (RequestMock) Input(string, ...string) string { return "" }
func (RequestMock) InputInt(string, ...int) int { return 0 }
func (RequestMock) InputBool(string, ...bool) bool { return false }
func (RequestMock) All() map[string]interface{} { return map[string]interface{}{} }
func (RequestMock) Only(...string) map[string]interface{} { return map[string]interface{}{} }
func (RequestMock) Except(...string) map[string]interface{} { return map[string]interface{}{} }
func (RequestMock) Has(string) bool { return false }
func (RequestMock) Filled(string) bool { return false }
func (RequestMock) Json(interface{}) error { return nil }
func (RequestMock) Body() ([]byte, error) { return nil, nil }
func (RequestMock) BodyString() (string, error) { return "", nil }
func (RequestMock) BodyReader() interface{} { return nil }
func (RequestMock) File(string) (*interface{}, error) { return nil, nil }
func (RequestMock) Files(string) ([]*interface{}, error) { return nil, nil }
func (RequestMock) AllFiles() (map[string][]*interface{}, error) { return map[string][]*interface{}{}, nil }
func (RequestMock) HasFile(string) bool { return false }
func (RequestMock) Header(string, ...string) string { return "" }
func (RequestMock) Headers() http.Header { return http.Header{} }
func (RequestMock) HasHeader(string) bool { return false }
func (RequestMock) Bearer() string { return "" }
func (RequestMock) Method() string { return "GET" }
func (RequestMock) URL() *url.URL { return &url.URL{} }
func (RequestMock) Path() string { return "/" }
func (RequestMock) FullURL() string { return "/" }
func (RequestMock) Scheme() string { return "http" }
func (RequestMock) Host() string { return "localhost" }
func (RequestMock) Hostname() string { return "localhost" }
func (RequestMock) Port() int { return 80 }
func (RequestMock) IP() string { return "127.0.0.1" }
func (RequestMock) UserAgent() string { return "mock" }
func (RequestMock) Referer() string { return "" }
func (RequestMock) ContentType() string { return "" }
func (RequestMock) ContentLength() int64 { return 0 }
func (RequestMock) Cookie(string, ...string) string { return "" }
func (RequestMock) Cookies() []*http.Cookie { return []*http.Cookie{} }
func (RequestMock) HasCookie(string) bool { return false }
func (RequestMock) IsJson() bool { return false }
func (RequestMock) IsXml() bool { return false }
func (RequestMock) IsForm() bool { return false }
func (RequestMock) IsMultipart() bool { return false }
func (RequestMock) IsSecure() bool { return false }
func (RequestMock) IsAjax() bool { return false }
func (RequestMock) Accepts(string) bool { return true }
func (RequestMock) WantsJson() bool { return false }
func (RequestMock) Context() interface{} { return nil }
func (RequestMock) SetContext(string, interface{}) {}
func (RequestMock) GetContext(string) interface{} { return nil }
func (RequestMock) Fresh() bool { return false }
func (RequestMock) Stale() bool { return true }
func (RequestMock) IfModifiedSince() interface{} { return nil }
func (RequestMock) IfNoneMatch() string { return "" }
