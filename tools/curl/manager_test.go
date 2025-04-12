package curl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestManagerCall(t *testing.T) {
	curl := Curl{
		Method: http.MethodGet,
		URL:    "{{.Env.url}}",
		Requests: map[string]Request{
			"default": {},
		},
	}

	expected := "test"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))

	defer ts.Close()

	env := Env{
		"url": ts.URL,
	}

	manager := Manager{
		Client: ts.Client(),
		Curl:   curl,
		Env:    env,
	}

	ctx := context.Background()

	res, err := manager.Call(ctx, "default")

	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	res.Body.Close()

	bodyStr := string(body)

	if bodyStr != expected {
		t.Errorf("req.Body = %s, want %s", bodyStr, expected)
	}
}
