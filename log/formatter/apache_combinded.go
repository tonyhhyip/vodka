package formatter

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tonyhhyip/vodka"
)

func (a *ApacheCombinedFormatter) Format(e *logrus.Entry) ([]byte, error) {
	content, err := a.ApacheCommonFormatter.Format(e)
	if err != nil {
		return content, err
	}

	var buffer bytes.Buffer
	buffer.Write(content)
	data := e.Data
	context := data["context"].(vodka.Context)
	referer := context.GetHeader("referer")
	if referer == "" {
		referer = "-"
	}
	buffer.WriteString(fmt.Sprintf(" \"%s\" \"", referer))
	buffer.WriteString(context.GetHeader("User-Agent"))
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}
