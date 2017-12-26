package vodka

import (
	"mime/multipart"
	"net/http"

	"github.com/tonyhhyip/vodka/errors"
)

type ContextWrapper struct {
	Context Context
}

func (c *ContextWrapper) GetRequest() *http.Request {
	return c.Context.GetRequest()
}

func (c *ContextWrapper) GetResponse() http.ResponseWriter {
	return c.Context.GetResponse()
}

func (c *ContextWrapper) IsAborted() bool {
	return c.Context.IsAborted()
}

func (c *ContextWrapper) Abort() {
	c.Context.Abort()
}

func (c *ContextWrapper) Error(err error) errors.Error {
	return c.Context.Error(err)
}

func (c *ContextWrapper) Set(key string, value interface{}) {
	c.Context.Set(key, value)
}

func (c *ContextWrapper) Get(key string) (value interface{}, exists bool) {
	return c.Context.Get(key)
}

func (c *ContextWrapper) Param(key string) string {
	return c.Context.Param(key)
}

func (c *ContextWrapper) Query(key string) (string, bool) {
	return c.Context.Query(key)
}

func (c *ContextWrapper) QueryArray(key string) ([]string, bool) {
	return c.Context.QueryArray(key)
}

func (c *ContextWrapper) PostForm(key string) (string, bool) {
	return c.Context.PostForm(key)
}

func (c *ContextWrapper) PostFormArray(key string) ([]string, bool) {
	return c.Context.PostFormArray(key)
}

func (c *ContextWrapper) FormFile(name string) (*multipart.FileHeader, error) {
	return c.Context.FormFile(name)
}

func (c *ContextWrapper) MultipartForm() (*multipart.Form, error) {
	return c.Context.MultipartForm()
}

func (c *ContextWrapper) ClientIP() string {
	return c.Context.ClientIP()
}

func (c *ContextWrapper) ContentType() string {
	return c.Context.ContentType()
}

func (c *ContextWrapper) GetMethod() Method {
	return c.Context.GetMethod()
}

func (c *ContextWrapper) GetPath() string {
	return c.Context.GetPath()
}

func (c *ContextWrapper) GetHeader(key string) string {
	return c.Context.GetHeader(key)
}

func (c *ContextWrapper) Status(code int) {
	c.Context.Status(code)
}

func (c *ContextWrapper) GetStatus() int {
	return c.Context.GetStatus()
}

func (c *ContextWrapper) Header(key, value string) {
	c.Context.Header(key, value)
}

func (c *ContextWrapper) SetCookie(
	name string,
	value string,
	maxAge int,
	path string,
	domain string,
	secure bool,
	httpOnly bool,
) {
	c.Context.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (c *ContextWrapper) Cookie(name string) (string, error) {
	return c.Context.Cookie(name)
}

func (c *ContextWrapper) Data(data []byte) {
	c.Context.Data(data)
}
