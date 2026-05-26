package router

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"exchange-gateway/internal/handlers"
	"exchange-gateway/internal/storage"
	"exchange-gateway/internal/transform"
)

func TestHealth(t *testing.T) {
	dir := t.TempDir()
	st, err := storage.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	api := &handlers.API{
		Store:     st,
		Transform: transform.New("http://127.0.0.1:9"),
		Log:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		Version:   "v1",
	}
	mux := http.NewServeMux()
	RegisterAPI(mux, "v1", api)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
}

func TestStaticServesIndex(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("ok"), 0o644); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	RegisterStatic(mux, dir)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
}
