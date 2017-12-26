# Vodka

A web framework which is extensible using middleware

## Install
[Glide](https://glide.sh) is recommend to be used.

```bash
glide get github.com/tonyhhyip/vodka
```

## Usage

```go
package main
import "github.com/tonyhhyip/vodka"

func main() {
	app := vodka.NewVodka()
	app.AddHandler(...)
	app.Run(":3000")
}
```