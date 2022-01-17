package router

import (
	"github.com/obgnail/http-server-toy/context"
	"net/http"
	"strings"
)

const (
	PathSeparator    = "/"
	HandlerSeparator = "-"
)

type HandlerFunc func(response *context.Context)

type routeNode struct {
	word     string
	path     string
	isFuzzy  bool
	children map[string]*routeNode // map[childWord]*routeNode
}

func newRouterNode(word string, isFuzzy bool) *routeNode {
	return &routeNode{word: word, isFuzzy: isFuzzy, children: make(map[string]*routeNode)}
}

func newBlankRouterNode() *routeNode {
	return newRouterNode("", false)
}

func (r *routeNode) IsFuzzy() bool {
	return r.isFuzzy
}

func (r *routeNode) IsLeaf() bool {
	return len(r.children) == 0
}

////////////////////////////////

type Router struct {
	node          map[string]*routeNode  // map[httpMethod]*routeNode
	nodeFunctions map[string]HandlerFunc // map[method-path]*routeNode
}

func NewRouter() *Router {
	return &Router{
		node:          make(map[string]*routeNode),
		nodeFunctions: make(map[string]HandlerFunc),
	}
}

func (r *Router) Add(method, path string, handler HandlerFunc) {
	if _, ok := r.node[method]; !ok {
		r.node[method] = newBlankRouterNode()
	}
	curNode := r.node[method]
	for _, word := range splitPath(path) {
		if _, ok := curNode.children[word]; !ok {
			curNode.children[word] = newRouterNode(word, word[0] == ':')
		}
		curNode = curNode.children[word]
	}
	curNode.path = path
	key := method + HandlerSeparator + path
	r.nodeFunctions[key] = handler
}

func (r *Router) Get(method, path string) (route *routeNode, params map[string]string) {
	curNode, ok := r.node[method]
	if !ok {
		return
	}

	testMatchNode := []*routeNode{curNode}
	paths := splitPath(path)
	for _, word := range paths {
		var tmpNode []*routeNode
		for _, node := range testMatchNode {
			if node.IsFuzzy() || word == node.word {
				tmpNode = append(tmpNode, mapToSlice(node.children)...)
			}
			testMatchNode = tmpNode
		}
	}

	for i := len(testMatchNode); i >= 0; i-- {
		if testMatchNode[i].IsLeaf() {
			route = testMatchNode[i]
		}
	}

	if route != nil {
		curNode := route
		params = make(map[string]string)
		for _, word := range paths {
			if curNode.IsFuzzy() {
				params[curNode.word[1:]] = word
			}
			curNode = curNode.children[word]
		}
	}
	return
}

func (r *Router) Handle(ctx *context.Context) {
	req := ctx.GetRequest()
	method := req.GetMethod()
	path := req.GetPath()

	node, params := r.Get(method, path)
	if node != nil {
		ctx.SetParams(params)
		k := method + HandlerSeparator + path
		f, ok := r.nodeFunctions[k]
		if ok {
			f(ctx)
		} else {
			ctx.GetResponse().SetStatusCode(http.StatusNotFound)
		}
	} else {
		ctx.GetResponse().SetStatusCode(http.StatusNotFound)
	}
	ctx.GetResponse().SetContentLength()
}

func splitPath(path string) []string {
	paths := strings.Split(path, PathSeparator)
	var ret []string
	for _, p := range paths {
		if p != "" {
			ret = append(ret, p)
		}
	}
	return ret
}

func mapToSlice(m map[string]*routeNode) (l []*routeNode) {
	for _, v := range m {
		l = append(l, v)
	}
	return
}
