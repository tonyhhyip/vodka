package route

import "github.com/tonyhhyip/vodka"

type Router interface {
	Match(method string, path string) vodka.Handler
}
