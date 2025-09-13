// Package mocks - Webserver mock
package mocks

import (
	"context"
	"govel/packages/new/webserver/src/interfaces"
)

type WebserverMock struct {
	Routes []struct{ Method, Path string }
}

func (m *WebserverMock) Get(path string, h interfaces.HandlerInterface) interfaces.WebserverInterface {
	m.Routes = append(m.Routes, struct{ Method, Path string }{"GET", path})
	return m
}
func (m *WebserverMock) Post(path string, h interfaces.HandlerInterface) interfaces.WebserverInterface {
	m.Routes = append(m.Routes, struct{ Method, Path string }{"POST", path})
	return m
}
func (m *WebserverMock) Put(path string, h interfaces.HandlerInterface) interfaces.WebserverInterface {
	m.Routes = append(m.Routes, struct{ Method, Path string }{"PUT", path})
	return m
}
func (m *WebserverMock) Patch(path string, h interfaces.HandlerInterface) interfaces.WebserverInterface {
	m.Routes = append(m.Routes, struct{ Method, Path string }{"PATCH", path})
	return m
}
func (m *WebserverMock) Delete(path string, h interfaces.HandlerInterface) interfaces.WebserverInterface {
	m.Routes = append(m.Routes, struct{ Method, Path string }{"DELETE", path})
	return m
}
func (m *WebserverMock) Options(path string, h interfaces.HandlerInterface) interfaces.WebserverInterface {
	m.Routes = append(m.Routes, struct{ Method, Path string }{"OPTIONS", path})
	return m
}
func (m *WebserverMock) Head(path string, h interfaces.HandlerInterface) interfaces.WebserverInterface {
	m.Routes = append(m.Routes, struct{ Method, Path string }{"HEAD", path})
	return m
}
func (m *WebserverMock) Group(_ string, _ func(interfaces.WebserverInterface)) interfaces.WebserverInterface {
	return m
}
func (m *WebserverMock) Use(_ ...interfaces.MiddlewareInterface) interfaces.WebserverInterface {
	return m
}
func (m *WebserverMock) Middleware(_ ...interfaces.MiddlewareInterface) interfaces.WebserverInterface {
	return m
}
func (m *WebserverMock) Listen(_ ...string) error                                        { return nil }
func (m *WebserverMock) ListenTLS(_, _ string, _ ...string) error                        { return nil }
func (m *WebserverMock) Shutdown(_ context.Context) error                                { return nil }
func (m *WebserverMock) SetConfig(_ string, _ interface{}) interfaces.WebserverInterface { return m }
func (m *WebserverMock) GetConfig(_ string) interface{}                                  { return nil }
func (m *WebserverMock) Static(_, _ string) interfaces.WebserverInterface                { return m }
func (m *WebserverMock) Port(_ int) interfaces.WebserverInterface                        { return m }
func (m *WebserverMock) Host(_ string) interfaces.WebserverInterface                     { return m }
