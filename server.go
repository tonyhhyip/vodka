// Copyright 2018 Tony Yip. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
package vodka

import (
	"context"
	"net/http"
	"time"
)

// New return a Server instance by the given address.
func New(addr string) *Server {
	return &Server{
		Server: &http.Server{
			Addr: addr,
		},
		Timeout: 30 * time.Second,
		logger:  &NullLogger{},
	}
}

// Server contains *http.Server.
type Server struct {
	Server  *http.Server
	logger  Logger
	Timeout time.Duration
}

// SetTimeout set timeout.
func (srv *Server) SetTimeout(duration time.Duration) {
	srv.Timeout = duration
}

// SetLogger set logger.
func (srv *Server) SetLogger(logger Logger) {
	srv.logger = logger
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
func (srv *Server) ListenAndServe(handler Handler) error {
	srv.StandBy(handler)

	return srv.Server.ListenAndServe()
}

// ListenAndServeTLS listens on the TCP network address srv.Addr and
// then calls Serve to handle requests on incoming TLS connections.
// Accepted connections are configured to enable TCP keep-alives.
func (srv *Server) ListenAndServeTLS(certFile, keyFile string, handler Handler) error {
	srv.StandBy(handler)

	return srv.Server.ListenAndServeTLS(certFile, keyFile)
}

// StandBy ready everything required for running a server
func (srv *Server) StandBy(handler Handler) {
	srv.Server.Handler = srv.WrapHandler(handler)
}

// WrapHandler wrap handler into standard http.Handler
func (srv *Server) WrapHandler(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			base   context.Context
			cancel context.CancelFunc
		)
		base, cancel = context.WithTimeout(r.Context(), srv.Timeout)
		defer cancel()
		ctx := NewContext(base, srv, w, r)

		done := make(chan struct{})

		go func(ctx *Context, done chan struct{}) {
			handler.Handle(ctx)
			close(done)
		}(ctx, done)

		select {
		case <-done:
			break
		case <-ctx.Done():
			if err := ctx.Err(); err == context.DeadlineExceeded {
				http.Error(w, "Timeout", http.StatusGatewayTimeout)
			} else if err != nil {
				srv.logger.Error(err)
			}
			break
		}
	})
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
func ListenAndServe(addr string, handler Handler) error {
	srv := New(addr)

	return srv.ListenAndServe(handler)
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it
// expects HTTPS connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate, any intermediates, and the CA's certificate.
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error {
	srv := New(addr)

	return srv.ListenAndServeTLS(certFile, keyFile, handler)
}
