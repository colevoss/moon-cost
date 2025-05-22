package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicServer(t *testing.T) {
	server := New()
	route := server.Route("")

	expected := "hello world"

	route.Get("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})

	ts := httptest.NewServer(server.Mux)
	defer ts.Close()

	res, err := http.Get(ts.URL)

	if err != nil {
		t.Fatal(err)
	}

	response, err := io.ReadAll(res.Body)
	res.Body.Close()

	resStr := string(response)

	if resStr != expected {
		t.Errorf("Expectd response to be %s. Got %s", expected, resStr)
	}
}
