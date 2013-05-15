package main

import (
	"log"
	"net/http"
)

type PatchHandler struct {
	Store ArchiveStore
}

func (h *PatchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fromRev, toRev, err := getRevParams(req)
	if err != nil {
		log.Println(err)
		http.Error(w, "400", http.StatusBadRequest)
		return
	}

	part, err := firstPart(req)
	if err != nil {
		log.Println(err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	// TODO lock.

  revs, err := h.Store.Revs(req.URL.Path, 1)
  if err != nil {
    log.Println(err)
    http.Error(w, "500", http.StatusInternalServerError)
    return
  }
	if len(revs) == 0 {
		http.NotFound(w, req)
		return
	}

	if revs[0].Rev != fromRev {
		http.Error(w, "409 Not Head", http.StatusConflict)
		return
	}

	defer part.Close()
	h.Store.Store(part, req.URL.Path, toRev)
}