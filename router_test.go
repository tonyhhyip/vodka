package vodka

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockResponseWriter struct {
	header     http.Header
	statusCode int
	body       []byte
}

func (w *mockResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = http.Header{}
	}
	return w.header
}

func (w *mockResponseWriter) Write(p []byte) (n int, err error) {
	w.body = append(w.body, p...)
	return len(p), nil
}

func (w *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (w *mockResponseWriter) WriteHeader(code int) {
	w.statusCode = code
}

type mockResponseWriter2 struct {
	*mockResponseWriter
}

func (w *mockResponseWriter2) Push(target string, opts *http.PushOptions) error {
	return nil
}

func TestRouter(t *testing.T) {
	router := NewRouter()

	routed := false
	router.HandleFunc("GET", "/user/:name", func(ctx *Context) {
		routed = true
		want := "gopher"
		if !reflect.DeepEqual(ctx.UserValue("name"), want) {
			t.Fatalf("wrong wildcard values: want %v, got %v", want, ctx.UserValue("name"))
		}
	})

	w := new(mockResponseWriter)

	req, _ := http.NewRequest("GET", "/user/gopher", nil)
	router.Handler().Handle(NewContext(context.Background(), nil, w, req))

	if !routed {
		t.Fatal("routing failed")
	}
}

func TestRouterAPI(t *testing.T) {
	var get, head, options, post, put, patch, delete bool

	router := NewRouter()
	router.GETFunc("/GET", func(ctx *Context) {
		get = true
	})
	router.HEADFunc("/GET", func(ctx *Context) {
		head = true
	})
	router.OPTIONSFunc("/GET", func(ctx *Context) {
		options = true
	})
	router.POSTFunc("/POST", func(ctx *Context) {
		post = true
	})
	router.PUTFunc("/PUT", func(ctx *Context) {
		put = true
	})
	router.PATCHFunc("/PATCH", func(ctx *Context) {
		patch = true
	})
	router.DELETEFunc("/DELETEFunc", func(ctx *Context) {
		delete = true
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET", nil)
	router.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !get {
		t.Error("routing GET failed")
	}

	r, _ = http.NewRequest("HEAD", "/GET", nil)
	router.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !head {
		t.Error("routing HEAD failed")
	}

	r, _ = http.NewRequest("OPTIONS", "/GET", nil)
	router.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !options {
		t.Error("routing OPTIONS failed")
	}

	r, _ = http.NewRequest("POST", "/POST", nil)
	router.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !post {
		t.Error("routing POST failed")
	}

	r, _ = http.NewRequest("PUT", "/PUT", nil)
	router.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !put {
		t.Error("routing PUT failed")
	}

	r, _ = http.NewRequest("PATCH", "/PATCH", nil)
	router.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !patch {
		t.Error("routing PATCH failed")
	}

	r, _ = http.NewRequest("DELETEFunc", "/DELETEFunc", nil)
	router.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !delete {
		t.Error("routing DELETEFunc failed")
	}
}

func TestRouterRoot(t *testing.T) {
	router := NewRouter()
	recv := catchPanic(func() {
		router.GET("noSlashRoot", nil)
	})
	if recv == nil {
		t.Fatal("registering path not beginning with '/' did not panic")
	}
}

func TestRouterChaining(t *testing.T) {
	router1 := NewRouter()
	router2 := NewRouter()
	router1.NotFound = router2.Handler()

	fooHit := false
	router1.POSTFunc("/foo", func(ctx *Context) {
		fooHit = true
		ctx.Response.WriteHeader(http.StatusOK)
	})

	barHit := false
	router2.POSTFunc("/bar", func(ctx *Context) {
		barHit = true
		ctx.Response.WriteHeader(http.StatusOK)
	})

	r, _ := http.NewRequest("POST", "/foo", nil)
	w := httptest.NewRecorder()
	router1.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusOK && fooHit) {
		t.Errorf("Regular routing failed with router chaining.")
		t.FailNow()
	}

	r, _ = http.NewRequest("POST", "/bar", nil)
	w = httptest.NewRecorder()
	router1.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusOK && barHit) {
		t.Errorf("Chained routing failed with router chaining.")
		t.FailNow()
	}

	r, _ = http.NewRequest("POST", "/qax", nil)
	w = httptest.NewRecorder()
	router1.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusNotFound) {
		t.Errorf("NotFound behavior failed with router chaining.")
		t.FailNow()
	}
}

func TestRouterOPTIONS(t *testing.T) {
	handlerFunc := func(ctx *Context) {}

	router := NewRouter()
	router.POSTFunc("/path", handlerFunc)

	// test not allowed
	// * (server)
	r, _ := http.NewRequest("OPTIONS", "*", nil)
	w := httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusOK) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", w.Code, w.Header())
	} else if allow := w.Header().Get("Allow"); allow != "POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}

	// path
	r, _ = http.NewRequest("OPTIONS", "/path", nil)
	w = httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusOK) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", w.Code, w.Header())
	} else if allow := w.Header().Get("Allow"); allow != "POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}

	r, _ = http.NewRequest("OPTIONS", "/doesnotexist", nil)
	w = httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusNotFound) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", w.Code, w.Header())
	}

	// add another method
	router.GETFunc("/path", handlerFunc)

	// test again
	// * (server)
	r, _ = http.NewRequest("OPTIONS", "*", nil)
	w = httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusOK) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", w.Code, w.Header())
	} else if allow := w.Header().Get("Allow"); allow != "POST, GET, OPTIONS" && allow != "GET, POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}

	// path
	r, _ = http.NewRequest("OPTIONS", "/path", nil)
	w = httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusOK) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", w.Code, w.Header())
	} else if allow := w.Header().Get("Allow"); allow != "POST, GET, OPTIONS" && allow != "GET, POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}

	// custom handler
	var custom bool
	router.OPTIONSFunc("/path", func(ctx *Context) {
		custom = true
	})

	// test again
	// * (server)
	r, _ = http.NewRequest("OPTIONS", "*", nil)
	w = httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusOK) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", w.Code, w.Header())
	} else if allow := w.Header().Get("Allow"); allow != "POST, GET, OPTIONS" && allow != "GET, POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}
	if custom {
		t.Error("custom handler called on *")
	}

	// path
	r, _ = http.NewRequest("OPTIONS", "/path", nil)
	w = httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusOK) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", w.Code, w.Header())
	}
	if !custom {
		t.Error("custom handler not called")
	}
}

func TestRouterNotAllowed(t *testing.T) {
	handlerFunc := func(_ *Context) {}

	router := NewRouter()
	router.POSTFunc("/path", handlerFunc)

	// test not allowed
	r, _ := http.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusMethodNotAllowed) {
		t.Errorf("NotAllowed handling failed: Code=%d, Header=%v", w.Code, w.Header())
	} else if allow := w.Header().Get("Allow"); allow != "POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}

	// add another method
	router.DELETEFunc("/path", handlerFunc)
	router.OPTIONSFunc("/path", handlerFunc) // must be ignored

	// test again
	r, _ = http.NewRequest("GET", "/path", nil)
	w = httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == http.StatusMethodNotAllowed) {
		t.Errorf("NotAllowed handling failed: Code=%d, Header=%v", w.Code, w.Header())
	} else if allow := w.Header().Get("Allow"); allow != "POST, DELETEFunc, OPTIONS" && allow != "DELETEFunc, POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}

	// test custom handler
	w = httptest.NewRecorder()
	responseText := "custom method"
	router.MethodNotAllowed = HandlerFunc(func(ctx *Context) {
		ctx.Response.WriteHeader(http.StatusTeapot)
		ctx.Response.Write([]byte(responseText))
	})
	router.handle(NewContext(context.Background(), nil, w, r))
	if got := w.Body.String(); !(got == responseText) {
		t.Errorf("unexpected response got %q want %q", got, responseText)
	}
	if w.Code != http.StatusTeapot {
		t.Errorf("unexpected response code %d want %d", w.Code, http.StatusTeapot)
	}
	if allow := w.Header().Get("Allow"); allow != "POST, DELETEFunc, OPTIONS" && allow != "DELETEFunc, POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}
}

func TestRouterNotFound(t *testing.T) {
	handlerFunc := func(_ *Context) {}

	router := NewRouter()
	router.GETFunc("/path", handlerFunc)
	router.GETFunc("/dir/", handlerFunc)
	router.GETFunc("/", handlerFunc)

	testRoutes := []struct {
		route    string
		code     int
		location string
	}{
		{"/path/", 301, "/path"},   // TSR -/
		{"/dir", 301, "/dir/"},     // TSR +/
		{"", 301, "/"},             // TSR +/
		{"/PATH", 301, "/path"},    // Fixed Case
		{"/DIR/", 301, "/dir/"},    // Fixed Case
		{"/PATH/", 301, "/path"},   // Fixed Case -/
		{"/DIR", 301, "/dir/"},     // Fixed Case +/
		{"/../path", 301, "/path"}, // CleanPath
		{"/nope", 404, ""},         // NotFound
	}
	for _, tr := range testRoutes {
		r, _ := http.NewRequest("GET", tr.route, nil)
		w := httptest.NewRecorder()
		router.handle(NewContext(context.Background(), nil, w, r))
		if !(w.Code == tr.code && (w.Code == 404 || w.Header().Get("Location") == tr.location)) {
			t.Errorf("NotFound handling route %s failed: Code=%d, Header=%v", tr.route, w.Code, w.Header())
		}
	}

	// Test custom not found handler
	var notFound bool
	router.NotFound = HandlerFunc(func(ctx *Context) {
		ctx.Response.WriteHeader(404)
		notFound = true
	})
	r, _ := http.NewRequest("GET", "/nope", nil)
	w := httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == 404 && notFound == true) {
		t.Errorf("Custom NotFound handler failed: Code=%d, Header=%v", w.Code, w.Header())
	}

	// Test other method than GET (want 307 instead of 301)
	router.PATCHFunc("/path", handlerFunc)
	r, _ = http.NewRequest("PATCH", "/path/", nil)
	w = httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == 307 && fmt.Sprint(w.Header()) == "map[Location:[/path]]") {
		t.Errorf("Custom NotFound handler failed: Code=%d, Header=%v", w.Code, w.Header())
	}

	// Test special case where no node for the prefix "/" exists
	router = NewRouter()
	router.GETFunc("/a", handlerFunc)
	r, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	router.handle(NewContext(context.Background(), nil, w, r))
	if !(w.Code == 404) {
		t.Errorf("NotFound handling route / failed: Code=%d", w.Code)
	}
}

func TestRouterPanicHandler(t *testing.T) {
	router := NewRouter()
	panicHandled := false

	router.PanicHandler = func(ctx *Context, _ interface{}) {
		panicHandled = true
	}

	router.HandleFunc("PUT", "/user/:name", func(_ *Context) {
		panic("oops!")
	})

	w := new(mockResponseWriter)
	req, _ := http.NewRequest("PUT", "/user/gopher", nil)

	defer func() {
		if rcv := recover(); rcv != nil {
			t.Fatal("handling panic failed")
		}
	}()

	router.handle(NewContext(context.Background(), nil, w, req))

	if !panicHandled {
		t.Fatal("simulating failed")
	}
}

func TestRouterLookup(t *testing.T) {
	routed := false
	wantHandle := func(_ *Context) {
		routed = true
	}
	wantParams := userValue{&userData{"name", "gopher"}}

	router := NewRouter()

	ctx := &Context{}

	// try empty router first
	handle, tsr := router.Lookup("GET", "/nope", ctx)
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}

	// insert route and try again
	router.GETFunc("/user/:name", wantHandle)

	handle, tsr = router.Lookup("GET", "/user/gopher", ctx)
	if handle == nil {
		t.Fatal("Got no handle!")
	} else {
		handle.Handle(nil)
		if !routed {
			t.Fatal("Routing failed!")
		}
	}

	for _, param := range wantParams {
		if !reflect.DeepEqual(ctx.UserValue(param.key), param.value) {
			t.Fatalf("Wrong parameter values: want %v, got %v", param.value, ctx.UserValue(param.key))
		}
	}

	handle, tsr = router.Lookup("GET", "/user/gopher/", ctx)
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if !tsr {
		t.Error("Got no TSR recommendation!")
	}

	handle, tsr = router.Lookup("GET", "/nope", ctx)
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}
}

type mockFileSystem struct {
	opened bool
}

func (mfs *mockFileSystem) Open(name string) (http.File, error) {
	mfs.opened = true
	return nil, errors.New("this is just a mock")
}

func TestRouterServeFiles(t *testing.T) {
	router := NewRouter()
	mfs := &mockFileSystem{}

	recv := catchPanic(func() {
		router.ServeFiles("/noFilepath", mfs)
	})
	if recv == nil {
		t.Fatal("registering path not ending with '*filepath' did not panic")
	}

	router.ServeFiles("/*filepath", mfs)
	w := new(mockResponseWriter)
	r, _ := http.NewRequest("GET", "/favicon.ico", nil)
	router.handle(NewContext(context.Background(), nil, w, r))
	if !mfs.opened {
		t.Error("serving file failed")
	}
}

type testMiddleware struct {
	handled bool
}

func (m *testMiddleware) Wrap(next Handler) Handler {
	return HandlerFunc(func(ctx *Context) {
		m.handled = true

		next.Handle(ctx)
	})
}

func TestRouter_Use(t *testing.T) {
	m := &testMiddleware{}

	router := NewRouter()
	router.Use(m)

	w := new(mockResponseWriter)
	r, _ := http.NewRequest("GET", "/", nil)
	router.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !m.handled {
		t.Error("use middleware failed")
	}
}

func TestRouter_Handle(t *testing.T) {
	m := &testMiddleware{}

	router := NewRouter()
	router.HandleFunc("GET", "/", func(ctx *Context) {}, &HandlerOption{
		Middlewares: []Middleware{m},
	})
	router.HandleFunc("GET", "/other", func(ctx *Context) {})

	w := new(mockResponseWriter)
	r, _ := http.NewRequest("GET", "/other", nil)
	router.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if m.handled {
		t.Error("use handler option failed: the ")
	}

	w = new(mockResponseWriter)
	r, _ = http.NewRequest("GET", "/", nil)
	router.Handler().Handle(NewContext(context.Background(), nil, w, r))
	if !m.handled {
		t.Error("use handler option failed")
	}
}
