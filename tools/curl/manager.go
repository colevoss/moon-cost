package curl

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrRequestNotFound = errors.New("Request not found in file")
)

type Manager struct {
	Client  *http.Client
	Request Request
	Curl    Curl
	Env     Env
}

func (m *Manager) Build(ctx context.Context) (*http.Request, error) {
	params := Params{
		Env:    m.Env,
		Params: m.Request.Params,
	}

	var builder Builder

	builder.Method(m.Curl.Method)

	if err := builder.URL(params, m.Curl.URL); err != nil {
		return nil, err
	}

	if err := builder.Headers(params, m.Curl.Headers, m.Request.Headers); err != nil {
		return nil, err
	}

	if err := builder.Query(params, m.Curl.Query, m.Request.Query); err != nil {
		return nil, err
	}

	if err := builder.Body(params, m.Request.Body); err != nil {
		return nil, err
	}

	request, err := builder.Build(ctx)

	if err != nil {
		return nil, err
	}

	return request, nil
}
