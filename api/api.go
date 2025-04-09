package api

import (
	"fmt"
	"moon-cost/router"
)

type Config struct {
	Port int
}

type API struct {
	Server *router.Server
	Config Config
}

func New(config Config) *API {
	server := router.New()

	return &API{
		Server: server,
		Config: config,
	}
}

func (a *API) Port() string {
	return fmt.Sprintf(":%d", a.Config.Port)
}
