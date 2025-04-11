package api

import (
	"fmt"
	"io"
	"moon-cost/router"
	"moon-cost/services/auth"
	"net/http"
	"time"
)

type AuthController struct {
	Route *router.Route
	Auth  *auth.Service
}

func (a *AuthController) Init(api *API) {
	a.Route = api.Server.Route("/auth")

	a.Route.Post("/signup/{id}", a.Signup)
}

func (a *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	body, err := io.ReadAll(r.Body)

	fmt.Printf("r.Header: %v\n", r.Header)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	time.Sleep(10 * time.Second)

	fmt.Printf("%v\n", string(body))
	fmt.Fprintf(w, "Hello, %s", id)
}
