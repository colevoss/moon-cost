package curl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrRequestNotFound = errors.New("Request not found in file")
)

type Manager struct {
	Curl    Curl
	Request Request
	Env     Env
}

func BuildRequest(ctx context.Context, m Manager) (*http.Request, error) {
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

func CheckResponse(req Request, res *http.Response) error {
	// status is default
	if req.Expect.Status == 0 {
		return nil
	}

	if req.Expect.Status == res.StatusCode {
		return nil
	}

	return fmt.Errorf("Received status %d. Expected %s", req.Expect.Status, res.Status)
}
