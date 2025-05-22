package router

import "net/http"

type Server struct {
	Mux *http.ServeMux
}

func New() *Server {
	return &Server{
		Mux: http.NewServeMux(),
	}
}

func (s *Server) Route(path string) *Route {
	return &Route{
		Path:   path,
		Server: s,
	}
}
