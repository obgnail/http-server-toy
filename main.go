package main

import (
	"fmt"
	"github.com/obgnail/http-server-toy/context"
	"github.com/obgnail/http-server-toy/engine"
	"net/http"
)

func main() {
	eng := engine.Default()
	eng.GET("/echo/:name", func(ctx *context.Context) {
		ctx.GetResponse().SetHeader("testHeader", "hello")
		ret := fmt.Sprintf("hello, %s", ctx.GetParams()["name"])
		ctx.JSON(http.StatusOK, ret)
	})

	eng.Run("0.0.0.0", 6666)
}
