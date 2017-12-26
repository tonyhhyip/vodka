package route

import "github.com/tonyhhyip/vodka"

type Router interface {
	Match(method vodka.Method, path string) vodka.Handler
}

func WrapRouter(router Router) vodka.Handler {
	return func(c vodka.Context) {
		handler := router.Match(c.GetMethod(), c.GetPath())
		handler(c)
	}
}
