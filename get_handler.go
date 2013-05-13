package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type GetHandler struct {
	Store ArchiveStore
}

func (h *GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// check for .tar.gz extension

	parts := strings.Split(req.URL.Path, ".tar.gz")
	if len(parts) < 2 {
		http.NotFound(w, req)
		return
	}

	revs, err := h.Store.ListRevs(parts[0], 0)
	if err != nil {
		log.Println(err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	rev := revs[len(revs)-1]
	rc, err := h.Store.Get(parts[0], rev.Rev)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	defer rc.Close()
	io.Copy(w, rc)
}
