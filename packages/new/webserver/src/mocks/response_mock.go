// Package mocks - Response mock
package mocks

type ResponseMock struct{}

func (ResponseMock) Status(int) interface{} { return ResponseMock{} }
func (ResponseMock) Header(string, string) interface{} { return ResponseMock{} }
func (ResponseMock) Headers(map[string]string) interface{} { return ResponseMock{} }
func (ResponseMock) RemoveHeader(string) interface{} { return ResponseMock{} }
func (ResponseMock) Json(interface{}) interface{} { return ResponseMock{} }
func (ResponseMock) Text(string) interface{} { return ResponseMock{} }
func (ResponseMock) HTML(string) interface{} { return ResponseMock{} }
func (ResponseMock) Send([]byte, string) interface{} { return ResponseMock{} }
func (ResponseMock) Stream(interface{}, string) interface{} { return ResponseMock{} }
func (ResponseMock) File(string) interface{} { return ResponseMock{} }
func (ResponseMock) Download(string, ...string) interface{} { return ResponseMock{} }
func (ResponseMock) Redirect(string, ...int) interface{} { return ResponseMock{} }
func (ResponseMock) Cookie(*interface{}) interface{} { return ResponseMock{} }
func (ResponseMock) ClearCookie(string) interface{} { return ResponseMock{} }
func (ResponseMock) NoContent(...int) interface{} { return ResponseMock{} }
func (ResponseMock) RawWriter(func(w interface{})) interface{} { return ResponseMock{} }
