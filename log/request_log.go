package log

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/tonyhhyip/vodka"
	"github.com/tonyhhyip/vodka/pkg/plugins"
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
		ctx := &sizeCountContext{
			ContextWrapper: plugins.WrapContext(c),
		}
		c.Next(ctx)
		r.logger.WithField("context", ctx).Info()
	}
}

type sizeCountContext struct {
	plugins.ContextWrapper
	size int
}

func (c *sizeCountContext) Data(data []byte) {
	c.size += len(data)
	c.ContextWrapper.Data(data)
}

func (c *sizeCountContext) Abort() {
	c.Header("Content-Length", strconv.FormatInt(int64(c.size), 10))
	c.ContextWrapper.Abort()
}
