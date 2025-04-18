package curl

import (
	"bytes"
	"context"
	"io"
	"moon-cost/assert"
	"net/http"
	"net/url"
	"text/template"
)

type Builder struct {
	url    *url.URL
	method string
	body   io.Reader
	header http.Header
	query  url.Values
}

func (b *Builder) Build(ctx context.Context) (*http.Request, error) {
	assert.Ensure(b.url, "Url must not be nil")

	b.url.RawQuery = b.query.Encode()
	urlStr := b.url.String()

	r, err := http.NewRequestWithContext(ctx, b.method, urlStr, b.body)

	if err != nil {
		return nil, err
	}

	r.Header = b.header

	return r, nil
}

func (b *Builder) Method(method string) {
	b.method = method
}

func (b *Builder) URL(params Params, urlTmpl string) error {
	builderUrl, err := parseString("url", urlTmpl, params)

	if err != nil {
		return err
	}

	urlStr, err := io.ReadAll(builderUrl)

	if err != nil {
		return err
	}

	u, err := url.Parse(string(urlStr))

	if err != nil {
		return err
	}

	b.url = u

	return nil
}

func (b *Builder) Body(params Params, req Body) error {
	body, err := OpenBody(req)

	if err != nil {
		return err
	}

	defer body.Close()

	if body.Type == ReqBodyTypeNone {
		return nil
	}

	parsedBody, err := parseReader("body", body.Data(), params)

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
			parsedV, err := parseString("headers", v, params)

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

func (b *Builder) Query(params Params, queries ...map[string]string) error {
	values := url.Values{}

	for _, query := range queries {
		for k, v := range query {
			parsed, err := parseString("query", v, params)

			if err != nil {
				return err
			}

			value, err := io.ReadAll(parsed)

			if err != nil {
				return err
			}

			values.Set(k, string(value))
		}
	}

	b.query = values

	return nil
}

func parseString(name string, str string, params Params) (io.Reader, error) {
	t, err := template.New(name).Option("missingkey=error").Parse(str)

	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}

	if err := t.Execute(&buf, params); err != nil {
		return nil, err
	}

	return &buf, nil
}

func parseReader(name string, reader io.Reader, params Params) (io.Reader, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return parseString(name, string(data), params)
}
