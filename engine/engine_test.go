package engine

import (
	"fmt"
	"github.com/obgnail/http-server-toy/context"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewEngine(t *testing.T) {
	e := Default()
	e.GET("/echo/:name", func(ctx *context.Context) {
		ctx.GetResponse().SetHeader("testheader", "hello")
		ctx.JSON(http.StatusOK, "hello, world")
	})
	eng := Default()
	eng.Run("0.0.0.0", 6666)
}

func Dial() {
	resp, err := http.Get("echo/obgnail")
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
	fmt.Print(string(body))
}
