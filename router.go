package cherry

import (
	"strings"
)

type Router struct {
	tree map[string]*tree

	Routes []*Route
}

func NewRouter() *Router {
	r := &Router{}
	r.tree = make(map[string]*tree)
	return r
}

type Route struct {
	Method   string
	Pattern  string
	Business Business
}

func (this *Router) FindRoute(method string, urlPath string) *Route {
	if this.tree[method] == nil {
		return nil
	}
	return this.tree[method].findRoute(urlPath)
}

func (this *Router) AddRoute(method string, pattern string, b Business) {
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}
	r := &Route{Method: method, Pattern: pattern, Business: b}
	this.Routes = append(this.Routes, r)
	if this.tree[r.Method] == nil {
		this.tree[r.Method] = newTree()
	}
	this.tree[r.Method].addRoute(r)
}

type tree struct {
	root *treeNode
}

type treeNode struct {
	name  string
	end   *Route
	block bool
	next  map[string]*treeNode
}

func newTreeNode(name string) *treeNode {
	t := &treeNode{}
	t.name = name
	t.block = false
	t.next = make(map[string]*treeNode)
	return t
}

func newTree() *tree {
	t := &tree{}
	t.root = newTreeNode("/")
	return t
}

func (this *tree) addRoute(route *Route) {
	s := strings.Split(route.Pattern, "/")
	var node *treeNode = this.root
	for _, name := range s {
		if name == "" {
			continue
		}
		if node.block {
			break
		}
		if name == "*" {
			node.block = true
			node.end = route
			break
		}
		if strings.HasPrefix(name, ":") {
			name = ":"
		}
		if node.next[name] == nil {
			node.next[name] = newTreeNode(name)
		}
		node = node.next[name]
	}
	if node.end == nil {
		node.end = route
	}
}

func (this *tree) findRoute(pattern string) *Route {
	s := strings.Split(pattern, "/")
	var node *treeNode = this.root
	for _, name := range s {
		if name == "" {
			continue
		}
		if node.block {
			break
		}
		if node.next[name] == nil {
			if node.next[":"] == nil {
				return nil
			}
			node = node.next[":"]
		} else {
			node = node.next[name]
		}
	}
	return node.end
}
