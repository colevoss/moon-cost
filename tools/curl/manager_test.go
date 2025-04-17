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
	}

	request := Request{}

	expected := "test"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))

	defer ts.Close()

	env := Env{
		"url": ts.URL,
	}

	manager := Manager{
		Curl:    curl,
		Env:     env,
		Request: request,
	}

	ctx := context.Background()

	req, err := BuildRequest(ctx, manager)
	res, err := ts.Client().Do(req)

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
