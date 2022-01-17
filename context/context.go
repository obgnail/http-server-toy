package context

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	params   map[string]string
	request  *Request
	response *Response
}

func NewContext(request *Request) *Context {
	ctx := &Context{
		params:   make(map[string]string),
		request:  request,
		response: NewResponse(),
	}
	return ctx
}

func (ctx *Context) GetRequest() *Request {
	return ctx.request
}

func (ctx *Context) SetRequest(req *Request) {
	ctx.request = req
}

func (ctx *Context) GetResponse() *Response {
	return ctx.response
}

func (ctx *Context) SetResponse(resp *Response) {
	ctx.response = resp
}

func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

func (ctx *Context) GetParams() map[string]string {
	return ctx.params
}

func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.GetResponse().SetStatusCode(code)
	data, err := json.Marshal(obj)
	if err != nil {
		ctx.GetResponse().SetStatusCode(http.StatusInternalServerError)
	}
	ctx.GetResponse().SetBody(data)
}
