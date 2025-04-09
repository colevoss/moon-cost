package api

import (
	"fmt"
	"moon-cost/router"
	"moon-cost/services/auth"
	"net/http"
)

type AuthController struct {
	Route *router.Route
	Auth  *auth.Service
}

func (a *AuthController) Init(api *API) {
	a.Route = api.Server.Route("/auth")

	a.Route.Post("/signup", a.Signup)
}

func (a *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HELLO")
}
