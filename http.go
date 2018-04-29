// Copyright 2018 Tony Yip. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
package vodka

import (
	"context"
	"net/http"
)

// CastHandlerForHTTP convert vodka.Handle to http.Handler
func CastHandlerForHTTP(h Handler, logger Logger) *handler {
	return &handler{
		Handler: h,
		Logger:  logger,
	}
}

type handler struct {
	Handler
	Logger Logger
}

// ServeHTTP serve HTTP request
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := r.Context()
	base, cancel := context.WithCancel(base)
	defer cancel()
	ctx := NewContextWithLogger(base, h.Logger, w, r)

	done := make(chan bool)

	go func(ctx *Context, done chan bool) {
		h.Handle(ctx)
		done <- true
	}(ctx, done)

	select {
	case <-done:
		break
	case <-ctx.Done():
		if err := ctx.Err(); err == context.DeadlineExceeded {
			http.Error(w, "Timeout", http.StatusGatewayTimeout)
		}
		break
	}
}
