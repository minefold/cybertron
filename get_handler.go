package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type GetHandler struct {
	Store ArchiveStore
}

func (h *GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch {
	case strings.HasSuffix(req.URL.Path, ".tar.gz"):
		h.tarGz(w, req)

	case strings.HasSuffix(req.URL.Path, ".json"):
		h.json(w, req)

	default:
		http.NotFound(w, req)
	}
}

func (h *GetHandler) tarGz(w http.ResponseWriter, req *http.Request) {
	parts := strings.Split(req.URL.Path, ".tar.gz")

	revs, err := h.Store.Revs(parts[0], 1)
	if err != nil {
		log.Println(err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

  if len(revs) == 0 {
		http.NotFound(w, req)
		return
  }

	rc, err := h.Store.Get(parts[0], revs[0].Rev)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	defer rc.Close()
	io.Copy(w, rc)
}

func (h *GetHandler) json(w http.ResponseWriter, req *http.Request) {
	parts := strings.Split(req.URL.Path, ".json")

	revs, err := h.Store.Revs(parts[0], 0)
	if err != nil {
		log.Println(err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err = enc.Encode(revs); err != nil {
		log.Println(err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

}
