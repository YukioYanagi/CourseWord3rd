package transform

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientTransform_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/transform" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"<x/>"}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	out, err := c.Transform("json", "xml", `{"a":1}`)
	if err != nil {
		t.Fatal(err)
	}
	if out != "<x/>" {
		t.Fatalf("got %q", out)
	}
}

func TestClientTransform_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"detail":"bad"}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	_, err := c.Transform("json", "xml", `{}`)
	if err == nil {
		t.Fatal("expected error")
	}
}
