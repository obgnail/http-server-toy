package context

import (
	"net/http"
	"strings"
)

const (
	ContentType   = "content-type"
	ContentLength = "content-length"
	UserAgent     = "user-agent"
	Connection    = "connection"
)

type Request struct {
	Method string
	URL    string
	Path   string
	Proto  string // "HTTP/1.0"
	Header map[string][]string
	Query  map[string][]string
	Body   []byte
}

func NewRequest(method string, url string, proto string, header map[string][]string, body []byte) *Request {
	req := &Request{
		Method: method,
		URL:    url,
		Header: header,
		Body:   body,
		Proto:  proto,
	}
	req.ParsePath()
	return req
}

func (req *Request) GetMethod() string {
	return req.Method
}

func (req *Request) GetPath() string {
	return req.Path
}

func (req *Request) KeepAlive() bool {
	v, ok := req.Header[Connection]
	return ok && v[0] == "keep-alive"
}

func (req *Request) ParsePath() {
	// /liangbo?message=hello&author=liangbo
	req.Query = make(map[string][]string)
	idx := strings.Index(req.URL, "?")
	if idx == -1 {
		req.Path = req.URL
		return
	}
	req.Path = req.URL[:idx]
	queryStr := req.URL[idx+1:]
	args := strings.Split(queryStr, "&")
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		k := parts[0]
		v := parts[1]
		req.Query[k] = append(req.Query[k], v)
	}
}

var Methods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodPost:    {},
	http.MethodHead:    {},
	http.MethodConnect: {},
	http.MethodPut:     {},
	http.MethodDelete:  {},
	http.MethodOptions: {},
	http.MethodTrace:   {},
	http.MethodPatch:   {},
}

var BodyType = map[string]struct{}{
	"application/json": {},
	"application/form": {},
}

func ParseFirstLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return strings.ToUpper(line[:s1]), line[s1+1 : s2], strings.ToUpper(line[s2+1:]), true
}

func ValidMethod(method string) bool {
	_, ok := Methods[method]
	return ok
}

func ValidBodyType(bodyType string) bool {
	_, ok := Methods[bodyType]
	return ok
}
