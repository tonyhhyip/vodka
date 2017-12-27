package route

import (
	"github.com/tonyhhyip/vodka"
	"github.com/tonyhhyip/vodka/pkg/routes"
)

type BasicRouter interface {
	routes.Matchable
	routes.RouteApplyAble
}

type EssentialRouter interface {
	routes.Matchable
	routes.RouteTable
}

type EnhanceRouter interface {
	EssentialRouter
	routes.MiddlewareAble
}

type GroupRouter interface {
	EnhanceRouter
	routes.Groupable
}

func WrapRouter(router routes.Matchable) vodka.Handler {
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
