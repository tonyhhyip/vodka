package vodka

import (
	"mime/multipart"
	"net/http"
	"time"

	"github.com/tonyhhyip/vodka/errors"
	"github.com/tonyhhyip/vodka/route"
	"golang.org/x/net/context"
)

type Handler func (c Context)
type HandlersChain []Handler

type Engine interface {
	http.Handler

	New() Engine
	Default() Engine
	Run(addr string) error
	RunTLS(addr string, certFile string, keyFile string) error
	HandleContext(c Context)
	AddRouter(router route.Router)
}

type Context interface {
	context.Context

	Copy() Context
	HandlerName() string
	Handler() Handler
	Next()

	IsAborted() bool
	Abort()
	AbortWithStatus(code int)
	AbortWithStatusJSON(code int, jsonObj interface{})
	AbortWithError(code int, err error) errors.Error

	Error(err error) errors.Error

	Set(key string, value interface{})
	Get(key string) (value interface{}, exists bool)
	MustGet(key string) interface{}
	GetString(key string) (s string)
	GetBool(key string) (b bool)
	GetInt(key string) (i int)
	GetInt64(key string) (i64 int64)
	GetFloat64(key string) (f64 float64)
	GetTime(key string) (t time.Time)

	Param(key string) string
	Query(key string) string
	DefaultQuery(key, defaultValue string) string
	GetQuery(key string) (string, bool)
	QueryArray(key string) []string
	GetQueryArray(key string) ([]string, bool)
	PostForm(key string) string
	DefaultPostForm(key, defaultValue string) string
	GetPostForm(key string) (string, bool)
	PostFormArray(key string) []string
	GetPostFormArray(key string) ([]string, bool)
	FormFile(name string) (*multipart.FileHeader, error)
	MultipartForm() (*multipart.Form, error)

	Bind(obj interface{}) error
	BindJSON(obj interface{}) error

	ClientIP() string
	ContentType() string

	Status(code int)
	Header(key, value string)
	GetHeader(key string) string
	Push(path string) string

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

	String(code int, format string, values ...interface{})
	IndentedJSON(code int, obj interface{})
	JSON(code int, obj interface{})
	YAML(code int, obj interface{})
	XML(code int, obj interface{})
	Data(code int, contentType string, data []byte)

	Redirect(code int, location string)
	SendFile(filepath string)
}
