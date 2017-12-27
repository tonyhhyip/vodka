package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyhhyip/vodka"
)

func TestSimpleRouteHandle(t *testing.T) {
	assert := assert.New(t)
	r := CreateBasicRoute().(*simpleRoute)
	r.GET("/foo", notFound)
	assert.Equal(1, len(r.handlers[vodka.Get]))
	assert.Equal(false, r.handlers[vodka.Get][0].redirect)
	assert.Equal(1, len(r.handlers[vodka.Get][0].handler))
}
