package router

import (
	"fmt"
	"net/http"
	"net/url"
)

type Route struct {
	Path       string
	Server     *Server
	Middleware []Middleware
}

func (r *Route) Route(path string) *Route {
	joinedPath, err := combineRoutes(r.Path, path)

	if err != nil {
		panic(err)
	}

	return &Route{
		Path:       joinedPath,
		Middleware: r.Middleware[:],
		Server:     r.Server,
	}
}

func (r *Route) Use(middleware ...Middleware) *Route {
	for _, m := range middleware {
		r.Middleware = append(r.Middleware, m)
	}

	return r
}

func (r *Route) With(middleware ...Middleware) *Route {
	newRoute := r.Route("")

	return newRoute.Use(middleware...)
}

func (r *Route) Get(path string, handler http.HandlerFunc) {
	r.register(http.MethodGet, path, handler)
}

func (r *Route) Head(path string, handler http.HandlerFunc) {
	r.register(http.MethodHead, path, handler)
}

func (r *Route) Post(path string, handler http.HandlerFunc) {
	r.register(http.MethodPost, path, handler)
}

func (r *Route) Put(path string, handler http.HandlerFunc) {
	r.register(http.MethodPut, path, handler)
}

func (r *Route) Patch(path string, handler http.HandlerFunc) {
	r.register(http.MethodPatch, path, handler)
}

func (r *Route) Delete(path string, handler http.HandlerFunc) {
	r.register(http.MethodDelete, path, handler)
}

// func (r *Route) Connect(path string, handler http.HandlerFunc) {
// 	r.register(http.MethodConnect, path, handler)
// }

func (r *Route) Options(path string, handler http.HandlerFunc) {
	r.register(http.MethodOptions, path, handler)
}

func (r *Route) Trace(path string, handler http.HandlerFunc) {
	r.register(http.MethodTrace, path, handler)
}

func (r *Route) register(method string, path string, handler http.HandlerFunc) {
	joinedPath, err := combineRoutes(r.Path, path)

	if err != nil {
		panic(err)
	}

	middleware := ComposeMiddleware(r.Middleware...)

	pattern := fmt.Sprintf("%s %s", method, joinedPath)

	r.Server.Mux.Handle(pattern, middleware.Handle(handler))
}

func combineRoutes(a, b string) (string, error) {
	joined, err := url.JoinPath("/", a, b)

	if err != nil {
		return "", err
	}

	return url.PathUnescape(joined)
}
