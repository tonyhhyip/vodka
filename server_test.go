package vodka

import (
	"bytes"
	"testing"
)

func TestServer_SetLogger(t *testing.T) {
	var logger Logger
	srv := New("")
	srv.SetLogger(logger)
	if srv.logger != logger {
		t.Error("failed to set logger")
	}
}

func TestServerInit(t *testing.T) {
	body := []byte("foo")
	handler := HandlerFunc(func(ctx *Context) {
		ctx.Response.Write(body)
	})

	srv := New("")
	srv.StandBy(handler)

	resp := &mockResponseWriter{}
	srv.Server.Handler.ServeHTTP(resp, nil)

	if !bytes.Equal(resp.body, body) {
		t.Errorf("expected response body %q, got %q", body, resp.body)
	}
}
