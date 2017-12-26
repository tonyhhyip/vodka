package main

import (
	"github.com/tonyhhyip/vodka"
	"github.com/tonyhhyip/vodka/log"
)

func main() {
	app := vodka.NewVodka()
	app.AddHandler(log.RequestLogger(log.ApacheCommon))
	app.AddHandler(func(c vodka.Context) {
		c.Status(200)
		c.Header("Content-Type", "text/plains")
		c.Data([]byte("hello world"))
	})
	app.Run(":3000")
}
