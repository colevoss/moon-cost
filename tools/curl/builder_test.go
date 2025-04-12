package curl

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestBuilderMethod(t *testing.T) {
	builder := Builder{}
	builder.Method(http.MethodGet)

	r, err := builder.Build(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if r.Method != http.MethodGet {
		t.Errorf("request.Method = %s. want %s", r.Method, http.MethodGet)
	}
}

func TestBuilderURL(t *testing.T) {
	params := Params{
		Env: Env{
			"url": "www.foo.com",
		},
		Params: ReqestParams{
			"param1": "1",
			"param2": "2",
		},
	}

	tests := []struct {
		tmpl     string
		expected string
	}{
		{"{{.Env.url}}/{{.Params.param1}}", "www.foo.com/1"},
		{"http://localhost:8080/{{.Params.param2}}", "http://localhost:8080/2"},
		{"{{.Env.url}}/123", "www.foo.com/123"},
		{"www.test.com", "www.test.com"},
	}

	for _, test := range tests {
		b := Builder{}

		if err := b.URL(params, test.tmpl); err != nil {
			t.Fatal(err)
		}

		r, err := b.Build(context.Background())

		if err != nil {
			t.Fatal(err)
		}

		url := r.URL

		if url.String() != test.expected {
			t.Errorf("request.URL = %s. want %s", url, test.expected)
		}
	}
}

func TestBuilderBody(t *testing.T) {
	params := Params{
		Env: Env{
			"ENV_FOO": "BAR",
		},
		Params: ReqestParams{
			"one": "1",
			"two": "2",
		},
	}

	tests := []struct {
		expected string
		body     Body
	}{
		{
			`{ "id": "BAR-1-2" }`,
			Body{
				File: "./test-fixtures/test-simple-body.json",
			},
		},
		{
			// compacting via JSON at the moment
			`{"id":"BAR-1-2"}`,
			Body{
				JSON: json.RawMessage(`{ "id": "{{.Env.ENV_FOO}}-{{.Params.one}}-{{.Params.two}}" }`),
			},
		},
	}

	for _, test := range tests {
		b := Builder{}
		b.URL(params, "www.test.com")
		b.Method("GET")

		if err := b.Body(params, test.body); err != nil {
			t.Fatal(err)
		}

		r, err := b.Build(context.Background())

		if err != nil {
			t.Fatal(err)
		}

		data, err := io.ReadAll(r.Body)

		if err != nil {
			t.Fatal(err)
		}

		dataStr := strings.TrimSpace(string(data))

		if dataStr != test.expected {
			t.Errorf("request.Body = %s. want %s", dataStr, test.expected)
		}
	}
}

func TestBuilderNoneBody(t *testing.T) {
	params := Params{}
	body := Body{}
	b := Builder{}

	b.URL(params, "www.test.com")
	b.Method("GET")

	if err := b.Body(params, body); err != nil {
		t.Fatal(err)
	}

	r, err := b.Build(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if r.Body != nil {
		t.Errorf("request.Body = %v. want nil", r.Body)
	}
}

func TestBuilderHeaders(t *testing.T) {
	secret := "53CR57"
	params := Params{
		Env: Env{
			"secret": secret,
		},
		Params: map[string]string{
			"param": "test-param",
		},
	}

	tests := []struct {
		name        string
		curlHeaders map[string]string
		reqHeaders  map[string]string
		expected    map[string]string
	}{
		{
			name: "combines passed headers",
			curlHeaders: map[string]string{
				"Curl-Header": "curl-{{.Env.secret}}",
			},
			reqHeaders: map[string]string{
				"Request-Header": "request-{{.Env.secret}}",
			},
			expected: map[string]string{
				"Curl-Header":    "curl-" + secret,
				"Request-Header": "request-" + secret,
			},
		},
		{
			name: "handles nil headers",
			reqHeaders: map[string]string{
				"Request-Header": "request-{{.Env.secret}}",
			},
			expected: map[string]string{
				"Request-Header": "request-" + secret,
			},
		},
		{
			name: "request header overrised curl headers",
			curlHeaders: map[string]string{
				"Curl-Header": "curl-{{.Env.secret}}",
				"Both-Header": "curl-both-{{.Env.secret}}",
			},
			reqHeaders: map[string]string{
				"Request-Header": "request-{{.Env.secret}}",
				"Both-Header":    "request-both-{{.Env.secret}}",
			},
			expected: map[string]string{
				"Curl-Header":    "curl-" + secret,
				"Request-Header": "request-" + secret,
				"Both-Header":    "request-both-" + secret,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			b := Builder{}

			b.URL(params, "www.test.com")
			b.Method("GET")

			if err := b.Headers(params, test.curlHeaders, test.reqHeaders); err != nil {
				t.Fatal(err)
			}

			r, err := b.Build(context.Background())

			if err != nil {
				t.Fatal(err)
			}

			reqHeaders := r.Header

			for k, v := range test.expected {
				reqHeader := reqHeaders.Get(k)

				if reqHeader != v {
					t.Errorf("header[%s] = %s. want %s", k, reqHeader, v)
				}
			}
		})
	}
}
