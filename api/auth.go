package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
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
	a.Route.Post("/signin", a.Signin)
}

type Response struct {
	Hello string `json:"hello"`
}

type SignupBody struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (a *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body SignupBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	signup := auth.SignupInput{
		Email:     body.Email,
		Password:  body.Password,
		Firstname: body.Firstname,
		Lastname:  body.Lastname,
	}

	signupResult, err := a.Auth.Signup(ctx, signup)

	if err == auth.SignupAccountExistsError {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	if err := json.NewEncoder(w).Encode(signupResult); err != nil {
		fmt.Printf("err: %s\n", err)
	}
}

type SigninBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthController) Signin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body SigninBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signinInput := auth.SignInInput{
		Email:    body.Email,
		Password: body.Password,
	}

	signinResult, err := a.Auth.SignIn(ctx, signinInput)

	if err == auth.AccountNotFoundError {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(signinResult); err != nil {
		slog.Error("error encoding result", slog.String("err", err.Error()))
	}
}
