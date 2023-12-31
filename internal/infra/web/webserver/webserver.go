package webserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HandlerProps struct {
	Method string
	Path   string
	Func   http.HandlerFunc
}

type WebServer struct {
	WebServerPort string
	Router        chi.Router
	Handlers      []HandlerProps
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make([]HandlerProps, 0),
		WebServerPort: fmt.Sprintf(":%s", serverPort),
	}
}

func (s *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	s.Handlers = append(s.Handlers, HandlerProps{
		Method: method,
		Path:   path,
		Func:   handler,
	})
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for _, h := range s.Handlers {
		s.Router.Method(h.Method, h.Path, h.Func)
	}
	err := http.ListenAndServe(s.WebServerPort, s.Router)
	if err != nil {
		return
	}
}
