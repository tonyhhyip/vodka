package group

import "github.com/tonyhhyip/vodka"

type Route interface {
	Use(handlers ...vodka.Handler) Route
	Handle(method string, path string, handlers ...vodka.Handler) Route
	Any(path string, handlers ...vodka.Handler) Route
	GET(path string, handlers ...vodka.Handler) Route
	POST(path string, handlers ...vodka.Handler) Route
	DELETE(path string, handlers ...vodka.Handler) Route
	PATCH(path string, handlers ...vodka.Handler) Route
	PUT(path string, handlers ...vodka.Handler) Route
	OPTIONS(path string, handlers ...vodka.Handler) Route
	HEAD(path string, handlers ...vodka.Handler) Route
}

type RouteGroup interface {
	Group(prefix string, handlers ...vodka.Handler) RouteGroup
}