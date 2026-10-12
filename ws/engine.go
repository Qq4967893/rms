package ws

import (
	"backend-rms/core/logger"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

type Engine struct {
	routes      map[string]HandlersChain
	middlewares HandlersChain
}

func (engine *Engine) GET(path string, handlers ...HandlerFunc) {
	handlers = append(engine.middlewares, handlers...)
	engine.routes[path] = append(engine.routes[path], handlers...)
}

func (engine *Engine) Use(middleware ...HandlerFunc) *Engine {
	engine.middlewares = append(engine.middlewares, middleware...)
	return engine
}

func (engine *Engine) Run(addr string) error {
	server := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return server.ListenAndServe()
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	handlers := engine.findHandlers(c.Request)
	if handlers == nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	conn, err := engine.accept(c.Writer, c.Request)
	if err != nil {
		logger.Warn(err.Error())
		c.Abort()
		return
	}

	defer func(conn *websocket.Conn, code websocket.StatusCode, reason string) {
		conn.Close(code, reason)
	}(conn, websocket.StatusNormalClosure, "")

	c.handlers = handlers
	c.Conn = conn
	c.Next()
}

func (engine *Engine) findHandlers(r *http.Request) HandlersChain {
	path := r.URL.Path
	if handlers, ok := engine.routes[path]; ok && len(handlers) > 0 {
		return handlers
	}
	return nil
}

func (engine *Engine) accept(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	acceptOpts := &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		Subprotocols:       []string{"rms"},
	}
	return websocket.Accept(w, r, acceptOpts)
}

func NewEngine() *Engine {
	return &Engine{
		routes: make(map[string]HandlersChain),
	}
}
