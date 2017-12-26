package formatter

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tonyhhyip/vodka"
)

var (
	ErrContextMissing = errors.New("context is missing")
)

func (*ApacheCommonFormatter) Format(e *logrus.Entry) ([]byte, error) {
	var buffer bytes.Buffer
	data := e.Data
	context, ok := data["context"].(vodka.Context)
	if !ok {
		return nil, ErrContextMissing
	}

	request := context.GetRequest()
	response := context.GetResponse()

	buffer.WriteString(request.RemoteAddr)
	buffer.WriteString(" - ")
	auth := context.GetHeader("Authorization")
	pieces := strings.Split(auth, " ")
	remoteUser := "-"
	if pieces[0] == "Basic" {
		authString, _ := base64.StdEncoding.DecodeString(pieces[1])
		pieces = strings.Split(string(authString), ":")
		remoteUser = pieces[0]
	}
	buffer.WriteString(remoteUser)
	buffer.WriteString(fmt.Sprintf(" [%s] ", e.Time.Format("02/Jan/2006:15:04:05 -0700")))
	buffer.WriteString(fmt.Sprintf("\"%s %s %s\"", request.Method, request.RequestURI, request.Proto))
	buffer.WriteString(fmt.Sprintf(" %d %d", context.GetStatus(), response.Header().Get("Content-Length")))
	return buffer.Bytes(), nil
}
