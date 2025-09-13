// Package webserver - Unified Request implementation
// This file implements a framework-agnostic Request that satisfies interfaces.RequestInterface.
// Adapters should wrap their native request types into this structure or populate its fields
// to expose a consistent API to handlers and middleware.
package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"govel/packages/new/webserver/src/interfaces"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Request is a concrete implementation of interfaces.RequestInterface.
// It is intentionally generic and does not depend on a specific framework.
// Adapters can construct this from their native request types.
type Request struct {
	// Route parameters like ":id" captured from the path
	params map[string]string
	// Query string parameters
	query url.Values
	// Headers
	headers http.Header
	// Cookies
	cookies []*http.Cookie
	// Method (GET, POST, ...)
	method string
	// URL object
	url *url.URL
	// Path without query
	path string
	// Raw body bytes (lazy read support optional)
	body []byte
	// Cached string body
	bodyString string
	// Content length
	contentLength int64
	// Client IP
	ip string
	// TLS/HTTPS flag
	secure bool
	// Arbitrary context values
	context map[string]interface{}
	// Creation time
	startedAt time.Time
}

// NewRequest creates an empty Request with sane defaults.
func NewRequest() *Request {
	return &Request{
		params:    map[string]string{},
		query:     url.Values{},
		headers:   http.Header{},
		cookies:   []*http.Cookie{},
		method:    "GET",
		url:       &url.URL{},
		path:      "/",
		body:      nil,
		context:   map[string]interface{}{},
		startedAt: time.Now(),
	}
}

// WithNativeHTTP allows constructing a Request from a standard http.Request.
// Useful for adapters built on top of net/http.
func WithNativeHTTP(r *http.Request, params map[string]string) *Request {
	if r == nil {
		return NewRequest()
	}
	u := *r.URL
	cookies := []*http.Cookie{}
	for _, c := range r.Cookies() {
		cookies = append(cookies, c)
	}
	bodyBytes := []byte{}
	if r.Body != nil {
		b, _ := io.ReadAll(r.Body)
		bodyBytes = b
		_ = r.Body.Close()
	}
	return &Request{
		params:        cloneStringMap(params),
		query:         r.URL.Query(),
		headers:       r.Header.Clone(),
		cookies:       cookies,
		method:        r.Method,
		url:           &u,
		path:          r.URL.Path,
		body:          bodyBytes,
		bodyString:    string(bodyBytes),
		contentLength: r.ContentLength,
		ip:            clientIPFromHeaders(r),
		secure:        r.TLS != nil,
		context:       map[string]interface{}{},
		startedAt:     time.Now(),
	}
}

// Helpers
func cloneStringMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = v
	}
	return out
}

func clientIPFromHeaders(r *http.Request) string {
	if r == nil {
		return ""
	}
	// X-Forwarded-For may contain multiple IPs, the first is the client
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return xr
	}
	ip := r.RemoteAddr
	if i := strings.LastIndex(ip, ":"); i > 0 {
		ip = ip[:i]
	}
	return ip
}

// Parameters
func (rq *Request) Param(key string) string { return rq.params[key] }
func (rq *Request) ParamInt(key string) int {
	v := strings.TrimSpace(rq.Param(key))
	if v == "" {
		return 0
	}
	var n int
	_, _ = fmt.Sscanf(v, "%d", &n)
	return n
}
func (rq *Request) Params() map[string]string { return cloneStringMap(rq.params) }

// Query
func (rq *Request) Query(key string, defaultValue ...string) string {
	vals := rq.query[key]
	if len(vals) > 0 {
		return vals[0]
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}
func (rq *Request) QueryInt(key string, defaultValue ...int) int {
	v := rq.Query(key)
	if v == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	var n int
	_, _ = fmt.Sscanf(v, "%d", &n)
	if n == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return n
}
func (rq *Request) QueryBool(key string, defaultValue ...bool) bool {
	v := strings.ToLower(strings.TrimSpace(rq.Query(key)))
	switch v {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}
func (rq *Request) Queries() url.Values { return rq.query }

// Input (body)
func (rq *Request) Input(key string, defaultValue ...string) string {
	// Simple JSON map parser fallback
	var m map[string]interface{}
	_ = json.Unmarshal(rq.body, &m)
	if m != nil {
		if v, ok := m[key]; ok {
			switch vv := v.(type) {
			case string:
				return vv
			case float64:
				return strings.TrimRight(strings.TrimRight(fmtFloat(vv), "0"), ".")
			case bool:
				if vv {
					return "true"
				}
				return "false"
			}
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func fmtFloat(f float64) string {
	// Import strconv at the top if needed
	return fmt.Sprintf("%g", f) // using fmt.Sprintf as fallback
}

func (rq *Request) InputInt(key string, defaultValue ...int) int {
	v := rq.Input(key)
	if v == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	var n int
	_, _ = fmt.Sscanf(v, "%d", &n)
	if n == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return n
}
func (rq *Request) InputBool(key string, defaultValue ...bool) bool {
	v := strings.ToLower(strings.TrimSpace(rq.Input(key)))
	switch v {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}
func (rq *Request) All() map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range rq.query {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}
	var m map[string]interface{}
	_ = json.Unmarshal(rq.body, &m)
	for k, v := range m {
		out[k] = v
	}
	return out
}
func (rq *Request) Only(keys ...string) map[string]interface{} {
	all := rq.All()
	res := map[string]interface{}{}
	for _, k := range keys {
		if v, ok := all[k]; ok {
			res[k] = v
		}
	}
	return res
}
func (rq *Request) Except(keys ...string) map[string]interface{} {
	all := rq.All()
	ex := map[string]struct{}{}
	for _, k := range keys {
		ex[k] = struct{}{}
	}
	for k := range ex {
		delete(all, k)
	}
	return all
}
func (rq *Request) Has(key string) bool    { _, ok := rq.All()[key]; ok = ok; return ok }
func (rq *Request) Filled(key string) bool { return strings.TrimSpace(rq.Input(key)) != "" }

// JSON / Body
func (rq *Request) Json(target interface{}) error { return json.Unmarshal(rq.body, target) }
func (rq *Request) Body() ([]byte, error)         { return rq.body, nil }
func (rq *Request) BodyString() (string, error) {
	if rq.bodyString != "" {
		return rq.bodyString, nil
	}
	return string(rq.body), nil
}
func (rq *Request) BodyReader() io.Reader { return bytes.NewReader(rq.body) }

// Files (adapters should populate using multipart parsing)
func (rq *Request) File(key string) (*multipart.FileHeader, error)    { return nil, nil }
func (rq *Request) Files(key string) ([]*multipart.FileHeader, error) { return nil, nil }
func (rq *Request) AllFiles() (map[string][]*multipart.FileHeader, error) {
	return map[string][]*multipart.FileHeader{}, nil
}
func (rq *Request) HasFile(key string) bool { return false }

// Headers
func (rq *Request) Header(key string, defaultValue ...string) string {
	v := rq.headers.Get(key)
	if v != "" {
		return v
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}
func (rq *Request) Headers() http.Header      { return rq.headers }
func (rq *Request) HasHeader(key string) bool { return rq.headers.Get(key) != "" }
func (rq *Request) Bearer() string {
	a := rq.headers.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(a), "bearer ") {
		return strings.TrimSpace(a[7:])
	}
	return ""
}

// Metadata
func (rq *Request) Method() string { return strings.ToUpper(rq.method) }
func (rq *Request) URL() *url.URL  { return rq.url }
func (rq *Request) Path() string   { return rq.path }
func (rq *Request) FullURL() string {
	if rq.url == nil {
		return rq.path
	}
	return rq.url.String()
}
func (rq *Request) Scheme() string {
	if rq.secure {
		return "https"
	}
	return "http"
}
func (rq *Request) Host() string {
	if rq.url != nil && rq.url.Host != "" {
		return rq.url.Host
	}
	return rq.headers.Get("Host")
}
func (rq *Request) Hostname() string {
	h := rq.Host()
	if i := strings.Index(h, ":"); i > 0 {
		return h[:i]
	}
	return h
}
func (rq *Request) Port() int {
	h := rq.Host()
	if i := strings.Index(h, ":"); i > 0 {
		var p int
		_, _ = fmt.Sscanf(h[i+1:], "%d", &p)
		return p
	}
	if rq.secure {
		return 443
	}
	return 80
}
func (rq *Request) IP() string           { return rq.ip }
func (rq *Request) UserAgent() string    { return rq.headers.Get("User-Agent") }
func (rq *Request) Referer() string      { return rq.headers.Get("Referer") }
func (rq *Request) ContentType() string  { return rq.headers.Get("Content-Type") }
func (rq *Request) ContentLength() int64 { return rq.contentLength }

// Cookies
func (rq *Request) Cookie(name string, defaultValue ...string) string {
	for _, c := range rq.cookies {
		if c.Name == name {
			return c.Value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}
func (rq *Request) Cookies() []*http.Cookie { return rq.cookies }
func (rq *Request) HasCookie(name string) bool {
	for _, c := range rq.cookies {
		if c.Name == name {
			return true
		}
	}
	return false
}

// Validation helpers
func (rq *Request) IsJson() bool {
	return strings.Contains(strings.ToLower(rq.ContentType()), "application/json")
}
func (rq *Request) IsXml() bool {
	ct := strings.ToLower(rq.ContentType())
	return strings.Contains(ct, "application/xml") || strings.Contains(ct, "text/xml")
}
func (rq *Request) IsForm() bool {
	return strings.Contains(strings.ToLower(rq.ContentType()), "application/x-www-form-urlencoded")
}
func (rq *Request) IsMultipart() bool {
	return strings.Contains(strings.ToLower(rq.ContentType()), "multipart/form-data")
}
func (rq *Request) IsSecure() bool { return rq.secure }
func (rq *Request) IsAjax() bool {
	return strings.ToLower(rq.headers.Get("X-Requested-With")) == "xmlhttprequest"
}
func (rq *Request) Accepts(contentType string) bool {
	accept := strings.ToLower(rq.headers.Get("Accept"))
	return strings.Contains(accept, strings.ToLower(contentType))
}
func (rq *Request) WantsJson() bool { return rq.Accepts("application/json") }

// Context
func (rq *Request) Context() interface{}                     { return rq.context }
func (rq *Request) SetContext(key string, value interface{}) { rq.context[key] = value }
func (rq *Request) GetContext(key string) interface{}        { return rq.context[key] }

// Caching headers
func (rq *Request) Fresh() bool                { return false }
func (rq *Request) Stale() bool                { return !rq.Fresh() }
func (rq *Request) IfModifiedSince() time.Time { return time.Time{} }
func (rq *Request) IfNoneMatch() string        { return rq.headers.Get("If-None-Match") }

// Compile-time interface compliance check
var _ interfaces.RequestInterface = (*Request)(nil)
