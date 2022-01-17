package main

import (
	"github.com/obgnail/http-server-toy/context"
	"github.com/obgnail/http-server-toy/engine"
	"net/http"
)

func main() {
	eng := engine.Default()
	eng.GET("/echo/:name", func(ctx *context.Context) {
		ctx.GetResponse().SetHeader("testHeader", "hello")
		ctx.JSON(http.StatusOK, "hello, world")
	})

	eng.Run("0.0.0.0", 6666)
}
