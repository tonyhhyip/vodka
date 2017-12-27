package simple

import (
	"net/http"

	"github.com/tonyhhyip/vodka"
	"github.com/tonyhhyip/vodka/pkg/routes"
)

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

func (r *simpleRoute) Any(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.GET(path, handlers...)
	r.POST(path, handlers...)
	r.PATCH(path, handlers...)
	r.PUT(path, handlers...)
	r.DELETE(path, handlers...)
	return r
}

func (r *simpleRoute) HEAD(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Head, path, handlers...)
	return r
}

func (r *simpleRoute) GET(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.HEAD(path, handlers...)
	r.Handle(vodka.Get, path, handlers...)
	return r
}

func (r *simpleRoute) POST(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Head, path, handlers...)
	return r
}

func (r *simpleRoute) DELETE(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Delete, path, handlers...)
	return r
}

func (r *simpleRoute) PATCH(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Patch, path, handlers...)
	return r
}

func (r *simpleRoute) PUT(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Put, path, handlers...)
	return r
}

func (r *simpleRoute) OPTIONS(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Options, path, handlers...)
	return r
}

func (r *simpleRoute) Handle(method vodka.Method, path string, handlers ...vodka.Handler) routes.RouteApplyAble {
	r.add(false, method, path, handlers...)
	return r
}

func (r *simpleRoute) add(redirect bool, method vodka.Method, path string, h ...vodka.Handler) {
	handlers, exists := r.handlers[method]
	if !exists {
		r.handlers[method] = make([]*routeHandler, 0)
		handlers = r.handlers[method]
	}
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
	u := c.GetRequest().URL
	u.Path += "/"
	c.Status(http.StatusMovedPermanently)
	c.Header("Location", u.String())
}
