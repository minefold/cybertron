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

	rev, err := h.Store.Head(req.URL.Path)
	if err != nil {
		log.Println(err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	if rev == nil {
		http.NotFound(w, req)
		return
	}

	if rev.Rev != fromRev {
		http.Error(w, "409 Not Head", http.StatusConflict)
		return
	}

	defer part.Close()
	h.Store.Store(part, req.URL.Path, toRev)
}