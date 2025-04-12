package curl

import (
	"context"
	"errors"
	"net/http"
)

var (
	RequestNotFoundError = errors.New("Request not found in file")
)

type Manager struct {
	Client *http.Client
	Curl   Curl
	Env    Env
}

func (m *Manager) Call(ctx context.Context, name string) (*http.Response, error) {
	req, ok := m.Curl.Req(name)

	if !ok {
		return nil, RequestNotFoundError
	}

	params := Params{
		Env:    m.Env,
		Params: req.Params,
	}

	var builder Builder

	builder.Method(m.Curl.Method)

	if err := builder.URL(params, m.Curl.URL); err != nil {
		return nil, err
	}

	if err := builder.Body(params, req.Body); err != nil {
		return nil, err
	}

	if err := builder.Headers(params, m.Curl.Headers, req.Headers); err != nil {
		return nil, err
	}

	request, err := builder.Build(ctx)

	if err != nil {
		return nil, err
	}

	client := m.client()
	return client.Do(request)
}

func (m *Manager) client() *http.Client {
	if m.Client == nil {
		return http.DefaultClient
	}

	return m.Client
}
