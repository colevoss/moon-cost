package api

import (
	"encoding/json"
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

	a.Route.Post("/signup/{id}", a.Signup)
}

type Response struct {
	Hello string `json:"hello"`
}

func (a *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	hello := Response{
		Hello: id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	if err := json.NewEncoder(w).Encode(hello); err != nil {
		fmt.Printf("err: %s\n", err)
	}
}
