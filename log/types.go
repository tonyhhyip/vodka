package log

import (
	"github.com/sirupsen/logrus"
	"github.com/tonyhhyip/vodka"
	"github.com/tonyhhyip/vodka/log/formatter"
)

func RequestLogger(logType RequestLogType) vodka.Handler {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	if logType == ApacheCommon {
		logger.Formatter = &formatter.ApacheCommonFormatter{}
	} else {
		logger.Formatter = &formatter.ApacheCombinedFormatter{
			ApacheCommonFormatter: formatter.ApacheCommonFormatter{},
		}
	}

	log := &requestLog{
		logger: logger,
	}
	return log.createHandler()
}
