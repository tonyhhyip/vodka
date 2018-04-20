// Copyright 2018 Tony Yip. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
package vodka

import "net/http"

// Handler for processing incoming requests.
type Handler interface {
	Handle(*Context)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(*Context)

// Handle calls f(ctx).
func (f HandlerFunc) Handle(ctx *Context) {
	f(ctx)
}

// HandlerOption option for handler.
type HandlerOption struct {
	Middlewares []Middleware
}

// NewHandlerOption returns HandlerOption instance by the
// given middlewares.
func NewHandlerOption(middlewares ...Middleware) *HandlerOption {
	return &HandlerOption{
		Middlewares: middlewares,
	}
}

// HttpHandle for wrap http.Handler into vodka.Handler
type HttpHandler struct {
	Handler http.Handler
}

// Handle calls Handler.ServeHTTP(ctx).
func (h HttpHandler) Handle(ctx *Context) {
	h.Handler.ServeHTTP(ctx.Response, ctx.Request)
}
