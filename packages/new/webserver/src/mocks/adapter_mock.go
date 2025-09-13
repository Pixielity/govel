// Package mocks - Adapter mock
package mocks

import (
	"context"
	"govel/new/webserver/interfaces"
)

type AdapterMock struct {
	Handled []struct{ Method, Path string }
}

func (a *AdapterMock) Init(_ map[string]interface{}, _ []interfaces.MiddlewareInterface) error {
	return nil
}
func (a *AdapterMock) SetConfig(_ string, _ interface{})       {}
func (a *AdapterMock) GetConfig(_ string) interface{}          { return nil }
func (a *AdapterMock) Use(_ ...interfaces.MiddlewareInterface) {}
func (a *AdapterMock) Handle(method, path string, _ interfaces.HandlerInterface, _ ...interfaces.MiddlewareInterface) {
	a.Handled = append(a.Handled, struct{ Method, Path string }{method, path})
}
func (a *AdapterMock) Group(_ string, register func(), _ ...interfaces.MiddlewareInterface) {
	register()
}
func (a *AdapterMock) Listen(_ string) error            { return nil }
func (a *AdapterMock) ListenTLS(_, _, _ string) error   { return nil }
func (a *AdapterMock) Shutdown(_ context.Context) error { return nil }
