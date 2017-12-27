package group

import (
	"net/url"

	"github.com/tonyhhyip/vodka"
)

type routeHandler struct {
	handler  []vodka.Handler
	route    string
	redirect bool
}

func (r *routeHandler) try(path string) (map[string]string, bool) {
	values := make(map[string]string)
	j := 0

	for i := 0; i < len(path); {
		switch {
		case j >= len(r.route):
			if r.route != "/" && len(r.route) > 0 && r.route[len(r.route)-1] == '/' {
				return values, true
			}
			return nil, false
		case r.route[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(r.route, isAlnum, j+1)
			val, _, i = match(path, matchPart(nextc), i)
			escval, err := url.QueryUnescape(val)
			if err != nil {
				return nil, false
			}
			values[":"+name] = escval
		case path[i] == r.route[j]:
			i += 1
			j += 1
		default:
			return nil, false
		}
	}

	if j != len(r.route) {
		return nil, false
	}

	return values, true
}

func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

func matchPart(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

func isAlnum(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
