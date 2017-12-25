package vodka

import (
	"mime/multipart"
	"net/http"

	"github.com/tonyhhyip/vodka/errors"
	ctx "golang.org/x/net/context"
)

type Method string

const (
	Head    Method = "HEAD"
	Get     Method = "GET"
	Post    Method = "POST"
	Put     Method = "PUT"
	Delete  Method = "DELETE"
	Patch   Method = "PATCH"
	Options Method = "OPTIONS"
)

type Handler func(c Context)

type runNext interface {
	Next()
}

type Engine interface {
	http.Handler
	runNext

	New() Engine
	Default() Engine
	Run(addr string) error
	RunTLS(addr string, certFile string, keyFile string) error
	HandleContext(c Context)
	AddHandler(handler Handler)
}

type Context interface {
	ctx.Context
	runNext

	GetRequest() *http.Request
	GetResponse() http.ResponseWriter

	IsAborted() bool
	Abort()
	Error(err error) errors.Error

	Set(key string, value interface{})
	Get(key string) (value interface{}, exists bool)

	Param(key string) string
	Query(key string) (string, bool)
	QueryArray(key string) ([]string, bool)
	PostForm(key string) (string, bool)
	PostFormArray(key string) ([]string, bool)
	FormFile(name string) (*multipart.FileHeader, error)
	MultipartForm() (*multipart.Form, error)

	ClientIP() string
	ContentType() string
	GetMethod() Method
	GetPath() string
	GetHeader(key string) string

	Status(code int)
	Header(key, value string)

	SetCookie(
		name string,
		value string,
		maxAge int,
		path string,
		domain string,
		secure bool,
		httpOnly bool,
	)
	Cookie(name string) (string, error)

	Data(code int, contentType string, data []byte)
}
