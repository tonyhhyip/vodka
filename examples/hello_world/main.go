package hello_world

import (
	"github.com/tonyhhyip/vodka"
	"github.com/tonyhhyip/vodka/log"
)

func main() {
	app := vodka.NewVodka()
	app.AddHandler(log.RequestLogger(log.ApacheCommon))
	app.AddHandler(func(c vodka.Context) {
		c.Data(200, "text/plain", []byte("hello world"))
	})
	app.Run(":3000")
}
