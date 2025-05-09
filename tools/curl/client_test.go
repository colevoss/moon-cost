package curl

import (
	"context"
	"fmt"
	"moon-cost/tools/env"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientLoadCurl(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var client Client

		err := client.LoadCurl("./test-fixtures/test-client-load-curl.json")

		if err != nil {
			t.Errorf("client.LoadCurl() = %s. want nil error", err)
		}

		expectedUrl := "test-client-load-curl"

		if client.Curl.URL != expectedUrl {
			t.Errorf("loaded curl url = %s. want %s", client.Curl.URL, expectedUrl)
		}
	})

	t.Run("no file", func(t *testing.T) {
		var client Client

		err := client.LoadCurl("./test-fixtures/file-doesnot-exist")

		if err == nil {
			t.Error("client.LoadCurl() = nil. want missing file error")
		}
	})
}

func TestClientLoadEnv(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var client Client

		err := client.LoadEnv("./test-fixtures/test-client-load-env.env")

		if err != nil {
			t.Errorf("client.LoadEnv() = %s. want nil error", err)
		}

		expectedVar := "bar"
		envVar := client.Env["FOO"]

		if envVar != expectedVar {
			t.Errorf("loaded env[FOO] = %s. want %s", envVar, expectedVar)
		}
	})

	t.Run("no file", func(t *testing.T) {
		var client Client

		err := client.LoadEnv("./test-fixtures/does-not-exist.env")

		if err == nil {
			t.Error("client.Env() = nil. want missing file error")
		}
	})
}

func TestClientUse(t *testing.T) {
	curl := Curl{
		Requests: map[string]Request{
			"a": {
				Expect: Expect{
					Status: 123,
				},
			},
		},
	}

	client := Client{
		Curl: curl,
	}

	err := client.Use("a")

	if err != nil {
		t.Errorf("client.Use(a) = %s. want nil error", err)
	}

	if client.Request.Expect.Status != 123 {
		t.Error("client.Request not set correctly")
	}
}

func TestClientExecute(t *testing.T) {
	curl := Curl{
		Method: http.MethodGet,
		URL:    "{{.Env.url}}",
	}

	request := Request{}

	expected := "success"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))

	defer ts.Close()

	env := env.Env{
		"url": ts.URL,
	}

	client := Client{
		Client:  ts.Client(),
		Curl:    curl,
		Env:     env,
		Request: request,
	}

	result, err := client.Execute(context.Background())

	if err != nil {
		t.Errorf("client.Execute() = _, %s. want nil err", err)
	}

	if result.Status != 200 {
		t.Errorf("result.Status = %d. want %d", result.Status, 200)
	}

	body := result.BodyString()

	if body != expected {
		t.Errorf("result body = %s. want %s", body, expected)
	}
}
