package gee

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup // root routerGroup
	router       *router
	groups       []*RouterGroup // store all groups
}

// RouterGroup struct
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middlewares
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

// ServeHTTP implement the interface of Handle
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// register middlewares
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// newContext create request context
	c := newContext(w, req)
	c.Handlers = middlewares
	e.router.handle(c)
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix:      g.prefix + prefix,
		middlewares: nil,
		parent:      g,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	g.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
