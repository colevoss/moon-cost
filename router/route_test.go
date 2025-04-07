package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpMethods(t *testing.T) {
	type testData struct {
		method   string
		hasBody  bool
		expected string
	}

	tests := []testData{
		{http.MethodGet, true, "test-GET"},
		{http.MethodHead, false, "test-HEAD"},
		{http.MethodPost, true, "test-POST"},
		{http.MethodPut, true, "test-PUT"},
		{http.MethodPatch, true, "test-PATCH"},
		{http.MethodDelete, true, "test-DELETE"},
		// {http.MethodConnect, false},
		{http.MethodOptions, false, "test-OPTIONS"},
		{http.MethodTrace, false, "test-TRACE"},
	}

	server := NewServer()
	route := server.Route("")

	createHandler := func(td testData) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("test", "test-"+td.method)
			fmt.Fprint(w, "test-"+td.method)
		}
	}

	for _, test := range tests {
		switch test.method {
		case http.MethodGet:
			route.Get("", createHandler(test))
		case http.MethodHead:
			route.Head("", createHandler(test))
		case http.MethodPost:
			route.Post("", createHandler(test))
		case http.MethodPut:
			route.Put("", createHandler(test))
		case http.MethodPatch:
			route.Patch("", createHandler(test))
		case http.MethodDelete:
			route.Delete("", createHandler(test))
		// case http.MethodConnect:
		case http.MethodOptions:
			route.Options("", createHandler(test))
		case http.MethodTrace:
			route.Trace("", createHandler(test))
		}
	}

	ts := httptest.NewServer(server.Mux)
	defer ts.Close()

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s Request", test.method), func(t *testing.T) {
			req, err := http.NewRequest(test.method, ts.URL, nil)

			if err != nil {
				t.Fatal(err)
			}

			res, err := http.DefaultClient.Do(req)

			if err != nil {
				t.Fatal(err)
			}

			if test.hasBody {
				data, err := io.ReadAll(res.Body)

				if err != nil {
					t.Fatal(err)
				}

				res.Body.Close()

				dataStr := string(data)

				if dataStr != test.expected {
					t.Errorf("Expected response body to be %s. Got %s", test.expected, dataStr)
				}
			} else {
				testHeader := res.Header.Get("test")

				if testHeader != test.expected {
					t.Errorf("Expected response header to be %s. Got %s", test.expected, testHeader)
				}
			}
		})
	}
}

func TestRoutes(t *testing.T) {
	server := NewServer()

	baseRouter := server.Route("/base")
	baseRouter.Get("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "base")
	})

	testRouter := baseRouter.Route("/test")
	testRouter.Get("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test-base")
	})

	baseRouter.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "base-%s", id)
	})

	ts := httptest.NewServer(server.Mux)
	defer ts.Close()

	tests := []struct {
		path string
		body string
	}{
		{"/base", "base"},
		{"/base/1", "base-1"},
		{"/base/100", "base-100"},
		{"/base/test", "test-base"}, // priority route since test is registered first
	}

	for _, test := range tests {
		res, err := http.Get(ts.URL + test.path)

		if err != nil {
			t.Fatal(err)
		}

		data, err := io.ReadAll(res.Body)

		if err != nil {
			t.Fatal(err)
		}

		dataStr := string(data)

		if dataStr != test.body {
			t.Errorf("Expected body %s. Got %s", test.body, dataStr)
		}
	}
}
