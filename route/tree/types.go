package tree

import "github.com/tonyhhyip/vodka"

type nodeType uint8

const (
	static nodeType = iota // default
	root
	paramNode
	catchAll
)

type node struct {
	path      string
	wildChild bool
	nType     nodeType
	maxParams uint8
	indices   string
	children  []*node
	handlers  []vodka.Handler
	priority  uint32
}

type param struct {
	key   string
	value string
}

type params []param

type router struct {
	trees map[vodka.Method]*node
}
