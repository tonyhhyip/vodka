package tree

import "github.com/tonyhhyip/vodka/route"

func New() route.EssentialRouter {
	return new(router)
}
