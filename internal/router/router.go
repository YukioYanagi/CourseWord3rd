package router

import (
	"net/http"

	"exchange-gateway/internal/handlers"
)

// RegisterAPI registers versioned HTTP routes on mux (Go 1.22+ patterns).
func RegisterAPI(mux *http.ServeMux, apiVersion string, api *handlers.API) {
	p := "/api/" + apiVersion
	mux.HandleFunc("GET "+p+"/health", api.Health)
	mux.HandleFunc("POST "+p+"/send", api.Send)
	mux.HandleFunc("GET "+p+"/received", api.ListReceived)
	mux.HandleFunc("GET "+p+"/received/{id}/download", api.Download)
}

// RegisterStatic serves files from webDir at "/".
func RegisterStatic(mux *http.ServeMux, webDir string) {
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
}
