package router

import (
	"net/http"
)

// func testMiddlware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("BIG TEST")
//
// 		next.ServeHTTP(w, r)
// 	})
// }
//
// func testMiddlware2(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("BIG TEST 2")
//
// 		next.ServeHTTP(w, r)
// 	})
// }

type Middleware func(http.Handler) http.Handler

func ComposeMiddleware(middleware ...Middleware) Composed {
	return Composed(func(next http.Handler) http.Handler {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}

		return next
	})
}

type Composed func(http.Handler) http.Handler

func (c Composed) Handle(handler http.HandlerFunc) http.Handler {
	return c(http.HandlerFunc(handler))
}
