package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]map[Method]http.HandlerFunc
	WebServerPort string
}

type Method string

const (
	POST Method = "POST"
	GET  Method = "GET"
)

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]map[Method]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method Method, path string, handler http.HandlerFunc) {
	if _, ok := s.Handlers[path]; !ok {
		s.Handlers[path] = make(map[Method]http.HandlerFunc)
	}

	s.Handlers[path][method] = handler
}

// Start Loop through the handlers and add them to the router
// register middleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, methodHandler := range s.Handlers {
		for method, handler := range methodHandler {
			if method == POST {
				s.Router.Post(path, handler)
			} else if method == GET {
				s.Router.Get(path, handler)
			}
		}
	}

	http.ListenAndServe(s.WebServerPort, s.Router)
}
