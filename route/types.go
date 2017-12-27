package route

import (
	"github.com/tonyhhyip/vodka"
)

type Router interface {
	Match(method vodka.Method, path string) ([]vodka.Handler, map[string]string)
}

func WrapRouter(router Router) vodka.Handler {
	return func(c vodka.Context) {
		handlers, params := router.Match(c.GetMethod(), c.GetPath())
		for k, v := range params {
			c.SetParam(k, v)
		}
		for i := 0; i < len(handlers) && !c.IsAborted(); i++ {
			handlers[i](c)
		}

		c.Next(c)
	}
}
