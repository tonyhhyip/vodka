// Copyright 2018 Tony Yip. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
package vodka

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"mime/multipart"
	"net/http"
	"net/url"
)

// MIME types
const (
	MIMEJSON = "application/json"
	MIMEXML  = "application/xml"
)

// Request methods.
const (
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodDelete  = "DELETE"
	MethodHead    = "HEAD"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodPatch   = "PATCH"
)

func newContext(base context.Context, s *Server, w http.ResponseWriter, r *http.Request) *Context {
	c := &Context{
		Context:  base,
		server:   s,
		Request:  r,
		Response: w,
	}

	c.Request = c.Request.WithContext(c)

	return c
}

type userValue []*userData

type userData struct {
	key   string
	value interface{}
}

// Context contains *http.Request and http.Response.
type Context struct {
	context.Context

	server    *Server
	userValue userValue

	Request  *http.Request
	Response http.ResponseWriter
}

// SetUserValue stores the given value under the given key in ctx.
func (ctx *Context) SetUserValue(key string, value interface{}) {
	data := userData{key: key, value: value}

	if ctx.userValue == nil {
		ctx.userValue = make(userValue, 0)
		ctx.userValue = append(ctx.userValue, &data)
		return
	}

	for i := 0; i < len(ctx.userValue); i++ {
		if ctx.userValue[i].key == key {
			ctx.userValue[i].value = value
			return
		}
	}

	ctx.userValue = append(ctx.userValue, &data)
}

// UserValue returns the value stored via SetUserValue* under the given key.
func (ctx *Context) UserValue(key string) interface{} {
	if ctx.userValue == nil {
		return nil
	}

	values := ctx.userValue
	for i := 0; i < len(values); i++ {
		if values[i].key == key {
			return values[i].value
		}
	}

	return nil
}

// Logger returns the server's logger.
func (ctx *Context) Logger() Logger {
	return ctx.server.logger
}

// SetServer for testing, do not use it in other places.
func (ctx *Context) SetServer(server *Server) {
	ctx.server = server
}

// Redirect replies to the request with a redirect to url,
// which may be a path relative to the request path.
func (ctx *Context) Redirect(url string, code int) {
	http.Redirect(ctx.Response, ctx.Request, url, code)
}

// IsDelete returns true if request method is DELETE.
func (ctx *Context) IsDelete() bool {
	return ctx.Request.Method == MethodDelete
}

// IsGet returns true if request method is GET.
func (ctx *Context) IsGet() bool {
	return ctx.Request.Method == MethodGet
}

// IsHead returns true if request method is HEAD.
func (ctx *Context) IsHead() bool {
	return ctx.Request.Method == MethodHead
}

// IsPost returns true if request method is POST.
func (ctx *Context) IsPost() bool {
	return ctx.Request.Method == MethodPost
}

// IsPut returns true if request method is PUT.
func (ctx *Context) IsPut() bool {
	return ctx.Request.Method == MethodPut
}

// IsAjax returns true if request is an AJAX (XMLHttpRequest) request.
func (ctx *Context) IsAjax() bool {
	return ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// Error is a shortcut of http.Error
func (ctx *Context) Error(error string, code int) {
	http.Error(ctx.Response, error, code)
}

// NotFound is a shortcut of http.NotFound
func (ctx *Context) NotFound() {
	http.NotFound(ctx.Response, ctx.Request)
}

// SetContentType set response Content-Type.
func (ctx *Context) SetContentType(v string) {
	ctx.Response.Header().Set("Content-Type", v)
}

// SetStatusCode is shortcut of http.ResponseWriter.WriteHeader.
func (ctx *Context) SetStatusCode(code int) {
	ctx.Response.WriteHeader(code)
}

// URL is shortcut of http.Request.URL.
func (ctx *Context) URL() *url.URL {
	return ctx.Request.URL
}

// FormValue is a shortcut of http.Request.FormValue.
func (ctx *Context) FormValue(key string) string {
	return ctx.Request.FormValue(key)
}

// PostFormValue is a shortcut of http.Request.PostFormValue.
func (ctx *Context) PostFormValue(key string) string {
	return ctx.Request.PostFormValue(key)
}

// ParseForm is a shortcut of http.Request.ParseForm.
func (ctx *Context) ParseForm() error {
	return ctx.Request.ParseForm()
}

// FormFile is a shortcut of http.Request.FormFile.
func (ctx *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return ctx.Request.FormFile(key)
}

// Host returns the value of http.Request.Host.
func (ctx *Context) Host() string {
	return ctx.Request.Host
}

// Referer is a shortcut of http.Request.Referer.
func (ctx *Context) Referer() string {
	return ctx.Request.Referer()
}

// Write is a shortcut of http.Response.Write.
func (ctx *Context) Write(p []byte) (n int, err error) {
	return ctx.Response.Write(p)
}

// Set response status code
func (ctx *Context) Status(code int) {
	ctx.Response.WriteHeader(code)
}

// JSON responses JSON data and custom status code to client.
func (ctx *Context) JSON(code int, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		ctx.Logger().Errorf("JSON error: %s\n", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.SetContentType(MIMEJSON)
	ctx.Status(code)
	ctx.Write(data)
}

// XML responses XML data and custom status code to client.
func (ctx *Context) XML(code int, v interface{}, headers ...string) {
	data, err := xml.Marshal(v)
	if err != nil {
		ctx.Logger().Errorf("XML error: %s\n", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	header := xml.Header
	if len(headers) > 0 {
		header = headers[0]
	}

	var bytes []byte
	bytes = append(bytes, header...)
	bytes = append(bytes, data...)

	ctx.SetContentType(MIMEXML)
	ctx.Status(code)
	ctx.Write(bytes)
}
