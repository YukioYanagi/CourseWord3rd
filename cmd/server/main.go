package main

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"exchange-gateway/internal/config"
	"exchange-gateway/internal/handlers"
	"exchange-gateway/internal/router"
	"exchange-gateway/internal/storage"
	"exchange-gateway/internal/transform"
)

func main() {
	cfg := config.Load()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	st, err := storage.New(cfg.DataDir)
	if err != nil {
		log.Error("storage init", "err", err)
		os.Exit(1)
	}

	api := &handlers.API{
		Store:     st,
		Transform: transform.New(cfg.PythonURL),
		Log:       log,
		Version:   cfg.APIVersion,
	}

	mux := http.NewServeMux()
	router.RegisterAPI(mux, cfg.APIVersion, api)
	router.RegisterStatic(mux, "web")

	handler := loggingMiddleware(log, cfg.APIVersion, versionHeader(cfg.APIVersion, mux))

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	log.Info("gateway listening", "addr", cfg.Addr, "api", "/api/"+cfg.APIVersion)
	if err := srv.ListenAndServe(); err != nil {
		log.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func loggingMiddleware(log *slog.Logger, apiVersion string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &respWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lw, r)
		log.Info("http",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"status", lw.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"client", r.RemoteAddr,
			"api_version", apiVersion,
		)
	})
}

type respWriter struct {
	http.ResponseWriter
	status int
}

func (rw *respWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func versionHeader(apiVersion string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("X-API-Version", apiVersion)
		}
		next.ServeHTTP(w, r)
	})
}
