package engine

import (
	"github.com/juju/errors"
	"github.com/obgnail/http-server-toy/connection"
	"github.com/obgnail/http-server-toy/context"
	"github.com/obgnail/http-server-toy/router"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Engine struct {
	*router.Router
}

func Default() *Engine {
	return &Engine{Router: router.NewRouter()}
}

func (e *Engine) GET(path string, handler router.HandlerFunc) {
	e.Add(http.MethodGet, path, handler)
}

func (e *Engine) POST(path string, handler router.HandlerFunc) {
	e.Add(http.MethodPost, path, handler)
}

func (e *Engine) process(clientConn *connection.Conn) {
	t := time.AfterFunc(time.Minute, func() { clientConn.Close() })
	defer t.Stop()

	for {
		req, err := clientConn.GetRequest()
		if err != nil {
			log.Warnf("get req err:", errors.Trace(err))
			break
		}
		ctx := context.NewContext(req)
		e.Handle(ctx)
		resp := ctx.GetResponse()
		err = clientConn.SendResponse(resp)
		if err != nil {
			log.Warnf("write resp err:", errors.Trace(err))
			break
		}

		if !req.KeepAlive() {
			break
		}
	}
	clientConn.Close()
}

func (e *Engine) Run(bindAddr string, bindPort int64) {
	l, err := connection.NewListener(bindAddr, bindPort)
	if err != nil {
		panic("init listener error")
	}
	for {
		clientConn, err := l.GetConn()
		if err != nil {
			log.Warn("get conn err:", errors.Trace(err))
			continue
		}
		go e.process(clientConn)
	}
}
