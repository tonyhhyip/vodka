package tree

import (
	"github.com/tonyhhyip/vodka"
	"github.com/tonyhhyip/vodka/pkg/routes"
)

func (r *router) Match(method vodka.Method, path string) ([]vodka.Handler, map[string]string) {
	if root := r.trees[method]; root != nil {
		if handlers, p, redirect := root.getValue(path); handlers != nil && len(handlers) > 0 {
			return handlers, p.asMap()
		} else if path != "/" {
			code := 301
			if method != vodka.Get {
				code = 307
			}

			if redirect {
				var url string
				if len(path) > 1 && path[len(path)-1] == '/' {
					url = path[:len(path)-1]
				} else {
					url = path + "/"
				}
				return []vodka.Handler{makeRedirect(code, url)}, nil
			}
		}
	}

	if allowed := r.allowed(method, path); len(allowed) > 0 {
		allowHandler := listAllowed(allowed)
		if method == vodka.Options {
			return []vodka.Handler{allowHandler}, nil
		} else {
			return []vodka.Handler{allowHandler, sendMethodNotAllowed}, nil
		}
	}

	return nil, nil
}

func (r *router) allowed(method vodka.Method, path string) (allow []vodka.Method) {
	allow = make([]vodka.Method, 0)
	if path == "*" { // server-wide
		for method := range r.trees {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			allow = append(allow, method)
		}
	} else { // specific path
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == method || method == "OPTIONS" {
				continue
			}

			handle, _, _ := r.trees[method].getValue(path)
			if handle != nil {
				// add request method to list of allowed methods
				allow = append(allow, method)
			}
		}
	}
	if len(allow) > 0 {
		allow = append(allow, vodka.Options)
	}
	return
}

func (r *router) Any(path string, handlers ...vodka.Handler) routes.RouteTable {
	return r.
		GET(path, handlers...).
		POST(path, handlers...).
		DELETE(path, handlers...).
		PATCH(path, handlers...).
		PUT(path, handlers...).
		OPTIONS(path, handlers...)
}

func (r *router) GET(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.HEAD(path, handlers...)
	r.Handle(vodka.Get, path, handlers...)
	return r
}

func (r *router) POST(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Post, path, handlers...)
	return r
}

func (r *router) DELETE(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Delete, path, handlers...)
	return r
}

func (r *router) PATCH(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Patch, path, handlers...)
	return r
}

func (r *router) PUT(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Put, path, handlers...)
	return r
}

func (r *router) OPTIONS(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Options, path, handlers...)
	return r
}

func (r *router) HEAD(path string, handlers ...vodka.Handler) routes.RouteTable {
	r.Handle(vodka.Head, path, handlers...)
	return r
}

func (r *router) Handle(method vodka.Method, path string, handlers ...vodka.Handler) routes.RouteApplyAble {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.trees == nil {
		r.trees = make(map[vodka.Method]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}

	root.addRoute(path, handlers)
	return r
}
