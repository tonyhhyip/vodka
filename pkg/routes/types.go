package routes

import "github.com/tonyhhyip/vodka"

type Matchable interface {
	Match(method vodka.Method, path string) ([]vodka.Handler, map[string]string)
}

type RouteApplyAble interface {
	Handle(method vodka.Method, path string, handlers ...vodka.Handler) RouteApplyAble
}

type RouteTable interface {
	RouteApplyAble
	Any(path string, handlers ...vodka.Handler) RouteTable
	GET(path string, handlers ...vodka.Handler) RouteTable
	POST(path string, handlers ...vodka.Handler) RouteTable
	DELETE(path string, handlers ...vodka.Handler) RouteTable
	PATCH(path string, handlers ...vodka.Handler) RouteTable
	PUT(path string, handlers ...vodka.Handler) RouteTable
	OPTIONS(path string, handlers ...vodka.Handler) RouteTable
	HEAD(path string, handlers ...vodka.Handler) RouteTable
}

type Groupable interface {
	Group(prefix string, handlers ...vodka.Handler) Groupable
}

type MiddlewareAble interface {
	Use(handlers ...vodka.Handler) MiddlewareAble
}
