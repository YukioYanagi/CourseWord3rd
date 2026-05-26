package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"exchange-gateway/internal/storage"
	"exchange-gateway/internal/transform"
	"exchange-gateway/internal/validate"
)

type API struct {
	Store     *storage.Store
	Transform *transform.Client
	Log       *slog.Logger
	Version   string
}

func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"version": a.Version,
	})
}

func (a *API) Send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	const maxMultipart = 32 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxMultipart+(1<<20))
	if err := r.ParseMultipartForm(maxMultipart); err != nil { // #nosec G120 -- bounded by MaxBytesReader above
		http.Error(w, "bad multipart", http.StatusBadRequest)
		return
	}
	format := strings.ToLower(r.FormValue("format"))
	if format == "" {
		format = "json"
	}
	transformTo := strings.ToLower(r.FormValue("transform_to"))

	file, hdr, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			a.Log.Warn("close upload failed", "err", closeErr)
		}
	}()

	body, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "read error", http.StatusInternalServerError)
		return
	}
	body = bytes.TrimPrefix(body, []byte{0xEF, 0xBB, 0xBF})

	f := storage.Format(format)
	switch f {
	case storage.FormatJSON:
		err = validate.JSON(body)
	case storage.FormatXML:
		err = validate.XML(body)
	case storage.FormatSOAP:
		err = validate.SOAP(body)
	default:
		http.Error(w, "format must be json, xml or soap", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "validation: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	saveBody := body
	saveFormat := f
	origName := hdr.Filename
	if transformTo != "" && transformTo != string(f) {
		res, terr := a.Transform.Transform(string(f), transformTo, string(body))
		if terr != nil {
			a.Log.Error("transform failed", "err", terr)
			http.Error(w, "transform: "+terr.Error(), http.StatusBadGateway)
			return
		}
		saveBody = []byte(res)
		saveFormat = storage.Format(transformTo)
		if origName != "" {
			ext := filepath.Ext(origName)
			base := strings.TrimSuffix(origName, ext)
			origName = base + "_converted." + transformTo
		}
	}

	rec, err := a.Store.SavePayload(saveFormat, origName, saveBody)
	if err != nil {
		a.Log.Error("save failed", "err", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(rec)
}

func (a *API) ListReceived(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	list := a.Store.List()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (a *API) Download(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	path, rec, err := a.Store.GetPath(id)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.NotFound(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}
	if err := a.Store.AssertPathInsideReceived(path); err != nil {
		a.Log.Error("download path rejected", "err", err)
		http.NotFound(w, r)
		return
	}

	ct := "application/octet-stream"
	switch rec.Format {
	case storage.FormatJSON:
		ct = "application/json"
	case storage.FormatXML, storage.FormatSOAP:
		ct = "application/xml"
	}
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, rec.Filename))
	f, err := os.Open(path) // #nosec G703 G304 -- path checked by Store.AssertPathInsideReceived
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			a.Log.Warn("close file failed", "err", closeErr)
		}
	}()
	st, err := f.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, rec.Filename, st.ModTime(), f)
}
