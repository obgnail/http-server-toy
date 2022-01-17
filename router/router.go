package router

import (
	"fmt"
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

type Router struct {
	root          map[string]*routeNode  // map[httpMethod]*routeNode
	rootFunctions map[string]HandlerFunc // map[method-path]*routeNode
}

func NewRouter() *Router {
	return &Router{
		root:          make(map[string]*routeNode),
		rootFunctions: make(map[string]HandlerFunc),
	}
}

func (r *Router) Add(method, path string, handler HandlerFunc) {
	if _, ok := r.root[method]; !ok {
		r.root[method] = newBlankRouterNode()
	}
	curNode := r.root[method]
	for _, word := range splitPath(path) {
		wildcardConflict(word, curNode.children)
		if _, ok := curNode.children[word]; !ok {
			curNode.children[word] = newRouterNode(word, word[0] == ':')
		}
		curNode = curNode.children[word]
	}
	curNode.path = path
	key := method + HandlerSeparator + path
	r.rootFunctions[key] = handler
}

func (r *Router) Get(method, path string) (*routeNode, map[string]string) {
	node, ok := r.root[method]
	if !ok {
		return nil, nil
	}

	params := make(map[string]string, 0)
	parts := splitPath(path)

	for _, p := range parts {
		var temp string
		for _, child := range node.children {
			if child.word == p || child.IsFuzzy() {
				if child.word[0] == ':' {
					k := child.word[1:]
					v := p
					params[k] = v
				}
				temp = child.word
			}
		}
		node = node.children[temp]
	}
	return node, params
}

func (r *Router) Handle(ctx *context.Context) {
	req := ctx.GetRequest()
	method := req.GetMethod()
	path := req.GetPath()

	node, params := r.Get(method, path)
	if node != nil {
		ctx.SetParams(params)
		k := method + HandlerSeparator + path
		f, ok := r.rootFunctions[k]
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

func wildcardConflict(path string, nodes map[string]*routeNode) {
	if len(nodes) != 0 {
		for k := range nodes {
			if strings.HasPrefix(k, ":") {
				panic(fmt.Errorf("word %s conflicts with existing wildcard %s", path, k))
			}
		}
	}
}
