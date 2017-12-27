package group

import (
	"net/http"

	"github.com/tonyhhyip/vodka"
)

type simpleRoute struct {
	handlers map[vodka.Method][]*routeHandler
	fallback vodka.Handler
}

func (r *simpleRoute) Match(method vodka.Method, path string) ([]vodka.Handler, map[string]string) {
	for _, handler := range r.handlers[method] {
		if params, ok := handler.try(path); ok {
			return handler.handler, params
		}
	}

	if r.fallback != nil {
		return []vodka.Handler{r.fallback}, nil
	}

	return []vodka.Handler{notFound}, nil
}

func notFound(c vodka.Context) {
	c.Status(http.StatusNotFound)
	c.Data([]byte("Not Found"))
}

func (r *simpleRoute) Any(path string, handlers ...vodka.Handler) Route {
	r.GET(path, handlers...)
	r.POST(path, handlers...)
	r.PATCH(path, handlers...)
	r.PUT(path, handlers...)
	r.DELETE(path, handlers...)
	return r
}

func (r *simpleRoute) HEAD(path string, handlers ...vodka.Handler) Route {
	r.Handle(vodka.Head, path, handlers...)
	return r
}

func (r *simpleRoute) GET(path string, handlers ...vodka.Handler) Route {
	r.HEAD(path, handlers...)
	r.Handle(vodka.Get, path, handlers...)
	return r
}

func (r *simpleRoute) POST(path string, handlers ...vodka.Handler) Route {
	r.Handle(vodka.Head, path, handlers...)
	return r
}

func (r *simpleRoute) DELETE(path string, handlers ...vodka.Handler) Route {
	r.Handle(vodka.Delete, path, handlers...)
	return r
}

func (r *simpleRoute) PATCH(path string, handlers ...vodka.Handler) Route {
	r.Handle(vodka.Patch, path, handlers...)
	return r
}

func (r *simpleRoute) PUT(path string, handlers ...vodka.Handler) Route {
	r.Handle(vodka.Put, path, handlers...)
}

func (r *simpleRoute) OPTIONS(path string, handlers ...vodka.Handler) Route {
	r.Handle(vodka.Options, path, handlers...)
	return r
}

func (r *simpleRoute) Handle(method vodka.Method, path string, handlers ...vodka.Handler) Route {
	r.add(false, method, path, handlers...)
	return r
}

func (r *simpleRoute) add(redirect bool, method vodka.Method, path string, h ...vodka.Handler) {
	handlers := r.handlers[method]
	for _, handler := range handlers {
		if handler.route == path {
			if !redirect {
				handler.handler = append(handler.handler, h...)
			}
			return
		}
	}

	handler := &routeHandler{
		route:    path,
		handler:  h,
		redirect: redirect,
	}

	r.handlers[method] = append(handlers, handler)

	n := len(path)
	if n > 0 && path[n-1] == '/' {
		r.add(true, method, path[:n-1], addSlashRedirect)
	}
}

func addSlashRedirect(c vodka.Context) {
	u := *c.GetRequest().URL
	u.Path += "/"
	c.Status(http.StatusMovedPermanently)
	c.Header("Location", u.String())
}
