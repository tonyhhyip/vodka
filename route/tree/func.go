package tree

import (
	"net/http"
	"strings"

	"github.com/tonyhhyip/vodka"
)

func countParams(path string) uint8 {
	var n uint
	for i := 0; i < len(path); i++ {
		if path[i] != ':' && path[i] != '*' {
			continue
		}
		n++
	}
	if n >= 255 {
		return 255
	}
	return uint8(n)
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func makeRedirect(code int, path string) vodka.Handler {
	return func(c vodka.Context) {
		c.Status(code)
		c.Header("Location", path)
	}
}

func listAllowed(methods []vodka.Method) vodka.Handler {
	results := make([]string, len(methods))

	for i, method := range methods {
		results[i] = string(method)
	}

	return func(c vodka.Context) {
		c.Header("Allow", strings.Join(results, ", "))
	}
}

func sendMethodNotAllowed(c vodka.Context) {
	c.Status(http.StatusMethodNotAllowed)
}
