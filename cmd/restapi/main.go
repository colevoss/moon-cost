package main

import (
	"fmt"
	"moon-cost/api"
	"moon-cost/services/auth"
	"net/http"
	"os"
)

func createAuth() *auth.Service {
	repo := auth.NoopRepo{}
	return auth.NewService(&repo)
}

func run() int {
	cfg := api.Config{
		Port: 8080,
	}

	restApi := api.New(cfg)

	authSvc := createAuth()

	authController := api.AuthController{
		Auth: authSvc,
	}

	authController.Init(restApi)

	if err := http.ListenAndServe(restApi.Port(), restApi.Server.Mux); err != nil {
		fmt.Printf("ERR")
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
