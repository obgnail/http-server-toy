package context

import "fmt"

type Response struct {
	Status     string
	StatusCode int
	Proto      string
	Header     map[string][]string
	Body       []byte
}

func (resp *Response) SetContentLength() {
	if _, ok := resp.Header[ContentLength]; ok {
		return
	}
	if bodyLen := len(resp.Body); bodyLen != 0 {
		value := fmt.Sprintf("%d", bodyLen)
		resp.Header[ContentLength] = []string{value}
	}
}

func NewResponse() *Response {
	return &Response{}
}

func (resp *Response) SetHeader(key, value string) {
	if resp.Header == nil {
		resp.Header = make(map[string][]string)
	}
	_, ok := resp.Header[key]
	if !ok {
		resp.Header[key] = []string{value}
	} else {
		resp.Header[key] = append(resp.Header[key], value)
	}
}

func (resp *Response) SetStatusCode(code int) {
	resp.StatusCode = code
}

func (resp *Response) SetBody(body []byte) {
	resp.Body = body
}

func (resp *Response) SetProto(proto string) {
	resp.Proto = proto
}
