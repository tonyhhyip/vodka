package vodka

import (
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tonyhhyip/vodka/errors"
)

const (
	defaultMemory = 100 << 20 // 100 MB
)

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

type BasicContext struct {
	handlers []Handler
	index    uint

	request  *http.Request
	response http.ResponseWriter

	abort bool

	ctx    map[string]interface{}
	params map[string]string
	errors []errors.Error
}

func (c *BasicContext) Next() {
	c.index++
	if int(c.index) > len(c.handlers) {
		if !c.IsAborted() {
			c.Abort()
		}
		return
	}

	c.handlers[c.index](c)
}

func (c *BasicContext) GetRequest() *http.Request {
	return c.request
}

func (c *BasicContext) GetResponse() http.ResponseWriter {
	return c.response
}

func (c *BasicContext) IsAborted() bool {
	return c.abort
}

func (c *BasicContext) Abort() {
	c.abort = true
}

func (c *BasicContext) Error(err error) errors.Error {
	var parsedError errors.Error
	switch err.(type) {
	case errors.Error:
		parsedError = err.(errors.Error)
	default:
		parsedError = errors.NewError(err, errors.ErrorTypePrivate, nil)
	}

	c.errors = append(c.errors, parsedError)

	return parsedError
}

func (c *BasicContext) Set(key string, value interface{}) {
	c.ctx[key] = value
}

func (c *BasicContext) Get(key string) (value interface{}, exists bool) {
	value, exists = c.ctx[key]
	return
}

func (c *BasicContext) Param(key string) string {
	return c.params[key]
}

func (c *BasicContext) Query(key string) (string, bool) {
	if values, ok := c.QueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *BasicContext) QueryArray(key string) ([]string, bool) {
	req := c.request
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

func (c *BasicContext) PostForm(key string) (string, bool) {
	if values, ok := c.PostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *BasicContext) PostFormArray(key string) ([]string, bool) {
	req := c.request
	req.ParseForm()
	req.ParseMultipartForm(defaultMemory)
	if values := req.PostForm[key]; len(values) > 0 {
		return values, true
	}
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		if values := req.MultipartForm.Value[key]; len(values) > 0 {
			return values, true
		}
	}
	return []string{}, false
}

func (c *BasicContext) FormFile(name string) (*multipart.FileHeader, error) {
	_, fh, err := c.request.FormFile(name)
	return fh, err
}

func (c *BasicContext) MultipartForm() (*multipart.Form, error) {
	err := c.request.ParseMultipartForm(defaultMemory)
	return c.request.MultipartForm, err
}

func (c *BasicContext) ClientIP() string {
	clientIP := c.GetHeader("X-Forwarded-For")
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = strings.TrimSpace(c.GetHeader("X-Real-Ip"))
	if len(clientIP) > 0 {
		return clientIP
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.request.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func (c *BasicContext) ContentType() string {
	return filterFlags(c.GetHeader("Content-Type"))
}

func (c *BasicContext) GetHeader(key string) string {
	if values, _ := c.request.Header[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}

func (c *BasicContext) GetMethod() Method {
	return Method(c.request.Method)
}

func (c *BasicContext) GetPath() string {
	return c.request.URL.Path
}

func (c *BasicContext) Status(code int) {
	c.response.WriteHeader(code)
}

func (c *BasicContext) Header(key, value string) {
	if len(value) == 0 {
		c.response.Header().Del(key)
	} else {
		c.response.Header().Set(key, value)
	}
}

func (c *BasicContext) SetCookie(
	name string,
	value string,
	maxAge int,
	path string,
	domain string,
	secure bool,
	httpOnly bool,
) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.response, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (c *BasicContext) Cookie(name string) (string, error) {
	cookie, err := c.request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (c *BasicContext) Data(code int, contentType string, data []byte) {
	c.Status(code)
	c.Header("Content-Type", contentType)
	_, _ = c.response.Write(data)
}

func (c *BasicContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *BasicContext) Done() <-chan struct{} {
	return nil
}

func (c *BasicContext) Err() error {
	return nil
}

func (c *BasicContext) Value(key interface{}) interface{} {
	if key == 0 {
		return c.request
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}
