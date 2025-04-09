package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	RequestNotFoundError = errors.New("Request not found in file")
)

type Manager struct {
	Client *http.Client
	Curl   CurlFile
	Env    Env
}

func (m *Manager) Request(ctx context.Context, name string) error {
	reqCfg, ok := m.Curl.Req(name)

	if !ok {
		return RequestNotFoundError
	}

	cfg := Config{
		Env:    m.Env,
		Params: reqCfg.Params,
	}

	url, err := m.Curl.ParsedURL(cfg)

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, m.Curl.Method, url, nil)

	if err != nil {
		return err
	}

	client := m.client()

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	fmt.Printf("data: %v\n", string(data))

	return nil
}

func (m *Manager) client() *http.Client {
	if m.Client == nil {
		return http.DefaultClient
	}

	return m.Client
}
