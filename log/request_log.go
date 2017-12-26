package log

import (
	"github.com/sirupsen/logrus"
	"github.com/tonyhhyip/vodka"
)

type RequestLogType int

const (
	ApacheCommon RequestLogType = iota
	ApacheCombinded
)

type requestLog struct {
	logger *logrus.Logger
}

func (r *requestLog) createHandler() vodka.Handler {
	return func(c vodka.Context) {
		c.Next()
		r.logger.WithField("context", c).Info()
	}
}
