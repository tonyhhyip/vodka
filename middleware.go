// Copyright 2018 Tony Yip. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
package vodka

import "net/http"

// Middleware interface.
type Middleware interface {
	Wrap(next Handler) Handler
}

type MiddlewareFunc func(next Handler) Handler

func (f MiddlewareFunc) Wrap(next Handler) Handler {
	return f(next)
}

type HttpMiddlewareFunc func(next http.Handler) http.Handler

func (f HttpMiddlewareFunc) Wrap(next Handler) Handler {
	return &HttpHandler{
		Handler: f(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context().(*Context)
			next.Handle(ctx)
		})),
	}
}
