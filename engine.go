package vodka

import (
	"net/http"

	"github.com/tonyhhyip/vodka/errors"
)

type engine struct {
	handlers []Handler
}

func (e *engine) Next() {}

func (e *engine) AddHandler(handler Handler) {
	e.handlers = append(e.handlers, handler)
}

func (e *engine) HandleContext(c Context) {
	e.handlers[0](c)
	c.Abort()
}

func (e *engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	c := &BasicContext{
		handlers: e.handlers,
		index:    0,
		engine:   e,
		request:  req,
		response: resp,
		abort:    false,
		ctx:      make(map[string]interface{}),
		params:   make(map[string]string),
		errors:   make([]errors.Error, 0),
	}

	e.HandleContext(c)
}

func (e *engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
