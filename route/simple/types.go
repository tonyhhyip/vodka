package simple

import (
	"github.com/tonyhhyip/vodka"
	"github.com/tonyhhyip/vodka/route"
)

func CreateBasicRoute() route.EssentialRouter {
	return &simpleRoute{}
}

type simpleRoute struct {
	handlers map[vodka.Method][]*routeHandler
	fallback vodka.Handler
}

type routeHandler struct {
	handler  []vodka.Handler
	route    string
	redirect bool
}
