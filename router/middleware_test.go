package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// middleware get registered and called
func TestMiddlewareCall(t *testing.T) {
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("middleware-header", "called")
		})
	}

	server := NewServer()
	useRouter := server.Route("/use")
	withRouter := server.Route("/with")

	useRouter.Use(middleware)
	useRouter.Get("", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	})

	withRouter.With(middleware).Get("", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	})

	ts := httptest.NewServer(server.Mux)
	defer ts.Close()

	useRes, err := http.Get(ts.URL + "/use")

	if err != nil {
		t.Fatal(err)
	}

	useHeader := useRes.Header.Get("middleware-header")

	if useHeader != "called" {
		t.Errorf("Expected middleware to set header when added with Used (%s). Got %s", "called", useHeader)
	}

	withRes, err := http.Get(ts.URL + "/with")

	if err != nil {
		t.Fatal(err)
	}

	withHeader := withRes.Header.Get("middleware-header")

	if withHeader != "called" {
		t.Errorf("Expected middleware to set header when added with With (%s). Got %s", "called", useHeader)
	}
}

// a middleware can terminate a request
func TestMiddlewareCanTerminateRequest(t *testing.T) {
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("test-terminate")

			if header == "yes" {
				http.Error(w, "error", http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	server := NewServer()
	route := server.Route("")

	route.Use(middleware)
	route.Get("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "success")
	})

	ts := httptest.NewServer(server.Mux)
	defer ts.Close()

	shouldFailReq, err := http.NewRequest(http.MethodGet, ts.URL, nil)

	if err != nil {
		t.Fatal(err)
	}

	shouldFailReq.Header.Add("test-terminate", "yes")
	shouldFailRes, err := http.DefaultClient.Do(shouldFailReq)

	if err != nil {
		t.Fatal(err)
	}

	if shouldFailRes.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected response with status %d. Got %s", http.StatusBadRequest, shouldFailRes.Status)
	}

	shouldSucceedRes, err := http.Get(ts.URL)

	if err != nil {
		t.Fatal(err)
	}

	if shouldSucceedRes.StatusCode != http.StatusOK {
		t.Errorf("Expected request with status %d. Got %s", http.StatusOK, shouldSucceedRes.Status)
	}
}

// Multiple middleware get chained in order
func TestChainedMiddlewareCalledInOrder(t *testing.T) {
	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("test-header", "1")

			next.ServeHTTP(w, r)
		})
	}

	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			prev := w.Header().Get("test-header")
			newHeader := prev + "2"
			w.Header().Set("test-header", newHeader)

			next.ServeHTTP(w, r)
		})
	}

	server := NewServer()
	route := server.Route("")

	route.Use(m1)
	route.Use(m2)

	route.Get("", func(w http.ResponseWriter, r *http.Request) {
		prev := w.Header().Get("test-header")
		newHeader := prev + "3"
		w.Header().Set("test-header", newHeader)

		fmt.Fprintf(w, "success")
	})

	ts := httptest.NewServer(server.Mux)
	defer ts.Close()

	res, err := http.Get(ts.URL)

	if err != nil {
		t.Fatal(err)
	}

	header := res.Header.Get("test-header")

	if header != "123" {
		t.Errorf("Expected header to be %s. Got %s", "123", header)
	}
}

// Middleware get inheritted from a route
func TestMiddlewareIsInherited(t *testing.T) {
	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("middleware-1", "yes")

			next.ServeHTTP(w, r)
		})
	}

	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("middleware-2", "yes")

			next.ServeHTTP(w, r)
		})
	}

	server := NewServer()
	router1 := server.Route("/one")
	router1.Use(m1)

	router2 := router1.Route("/two")
	router2.Use(m2)

	router1.Get("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "one")
	})

	router2.Get("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "two")
	})

	ts := httptest.NewServer(server.Mux)
	defer ts.Close()

	resOne, err := http.Get(ts.URL + "/one")

	if err != nil {
		t.Fatal(err)
	}

	resOneHeaderOne := resOne.Header.Get("middleware-1")
	resOneHeaderTwo := resOne.Header.Get("middleware-2")

	if resOneHeaderOne != "yes" {
		t.Errorf("Expected header to be %s. Got %s", "yes", resOneHeaderOne)
	}

	if resOneHeaderTwo != "" {
		t.Errorf("Expected header to be not be set. Got %s", resOneHeaderTwo)
	}

	resTwo, err := http.Get(ts.URL + "/one/two")

	if err != nil {
		t.Fatal(err)
	}

	resTwoHeaderOne := resTwo.Header.Get("middleware-1")
	resTwoHeaderTwo := resTwo.Header.Get("middleware-2")

	if resTwoHeaderOne != "yes" {
		t.Errorf("Expected header to be %s. Got %s", "yes", resTwoHeaderOne)
	}

	if resTwoHeaderTwo != "yes" {
		t.Errorf("Expected header to be %s. Got %s", "yes", resTwoHeaderOne)
	}
}
