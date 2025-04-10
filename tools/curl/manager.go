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
	Curl   Curl
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

	reqBody, err := m.reqBody(reqCfg)
	if err != nil {
		return err
	}
	defer reqBody.Close()

	req, err := http.NewRequestWithContext(ctx, m.Curl.Method, url, reqBody.Data())

	m.headers(req.Header, m.Curl, reqCfg)

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

func (m *Manager) headers(headers http.Header, curl Curl, req ReqConfig) {
	for k, v := range curl.Headers {
		headers.Add(k, v)
	}

	if len(req.Headers) == 0 {
		return
	}

	for k, v := range req.Headers {
		headers.Set(k, v)
	}
}

func (m *Manager) reqBody(cfg ReqConfig) (ReqBody, error) {
	if cfg.Body.File != "" {
		return FileReqBody(cfg.Body.File)
	}

	if cfg.Body.JSON != nil {
		return JSONReqBody(cfg.Body.JSON)
	}

	return NoneReqBody(), nil
}

func (m *Manager) client() *http.Client {
	if m.Client == nil {
		return http.DefaultClient
	}

	return m.Client
}
