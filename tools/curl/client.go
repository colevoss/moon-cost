package curl

import (
	"context"
	"errors"
	"fmt"
	"moon-cost/tools/env"
	"net/http"
	"os"
	"time"
)

var (
	ErrRequestNotFound = errors.New("Request not found in file")
)

type Client struct {
	Client      *http.Client
	Curl        Curl
	Request     Request
	RequestName string
	Env         env.Env
}

func (c *Client) LoadCurl(curlFile string) error {
	file, err := os.Open(curlFile)

	if err != nil {
		return err
	}

	defer file.Close()

	if err := c.Curl.Read(file); err != nil {
		return err
	}

	return nil
}

func (c *Client) LoadEnv(envFile string) error {
	env := env.Env{}

	file, err := os.Open(envFile)

	if err != nil {
		return err
	}

	defer file.Close()

	if err := env.Read(file); err != nil {
		return fmt.Errorf("%s:%s", envFile, err)
	}

	c.Env = env

	c.Env.AddEnviron(os.Environ())

	return nil
}

func (c *Client) Use(requestName string) error {
	request, ok := c.Curl.Req(requestName)

	if !ok {
		return ErrRequestNotFound
	}

	c.RequestName = requestName
	c.Request = request

	return nil
}

func (c *Client) Execute(ctx context.Context) (Result, error) {
	req, err := c.buildRequest(ctx)

	var result Result

	if err != nil {
		return result, err
	}

	httpClient := http.DefaultClient

	if c.Client != nil {
		httpClient = c.Client
	}

	start := time.Now()

	response, err := httpClient.Do(req)

	if err != nil {
		return result, err
	}

	end := time.Now()
	duration := end.Sub(start)

	result, err = NewResult(c.RequestName, response, c.Request, duration)

	return result, err
}

func (c *Client) buildRequest(ctx context.Context) (*http.Request, error) {
	params := Params{
		Env:    c.Env,
		Params: c.Request.Params,
	}

	var builder Builder

	builder.Method(c.Curl.Method)

	if err := builder.URL(params, c.Curl.URL); err != nil {
		return nil, err
	}

	if err := builder.Headers(params, c.Curl.Headers, c.Request.Headers); err != nil {
		return nil, err
	}

	if err := builder.Query(params, c.Curl.Query, c.Request.Query); err != nil {
		return nil, err
	}

	if err := builder.Body(params, c.Request.Body); err != nil {
		return nil, err
	}

	request, err := builder.Build(ctx)

	if err != nil {
		return nil, err
	}

	return request, nil
}
