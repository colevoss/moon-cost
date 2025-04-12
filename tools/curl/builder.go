package curl

import (
	"bytes"
	"context"
	"io"
	"moon-cost/assert"
	"net/http"
	"text/template"
)

type Builder struct {
	url    string
	method string
	body   io.Reader
	header http.Header
}

func (b *Builder) Build(ctx context.Context) (*http.Request, error) {
	r, err := http.NewRequestWithContext(ctx, b.method, b.url, b.body)

	if err != nil {
		return nil, err
	}

	r.Header = b.header

	return r, nil
}

func (b *Builder) Method(method string) {
	b.method = method
}

func (b *Builder) URL(params Params, url string) error {
	builderUrl, err := parseString(url, params)

	if err != nil {
		return err
	}

	urlStr, err := io.ReadAll(builderUrl)

	if err != nil {
		return err
	}

	b.url = string(urlStr)

	return nil
}

func (b *Builder) Body(params Params, req Body) error {
	body, err := OpenBody(req)

	if err != nil {
		return err
	}

	if body.Type == ReqBodyTypeNone {
		return nil
	}

	parsedBody, err := parseReader(body.Data(), params)

	if err != nil {
		return err
	}

	b.body = parsedBody

	return nil
}

func (b *Builder) Headers(params Params, headers ...map[string]string) error {
	header := make(http.Header)

	for _, h := range headers {
		assert.Ensure(h, "Header map should not be nil")

		for k, v := range h {
			parsedV, err := parseString(v, params)

			if err != nil {
				return err
			}

			parsedStr, err := io.ReadAll(parsedV)

			if err != nil {
				return err
			}

			header.Set(k, string(parsedStr))
		}
	}

	b.header = header

	return nil
}

func parseString(str string, params Params) (io.Reader, error) {
	t, err := template.New("").Parse(str)

	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}

	if err := t.Execute(&buf, params); err != nil {
		return nil, err
	}

	return &buf, nil
}

func parseReader(reader io.Reader, params Params) (io.Reader, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return parseString(string(data), params)
}
