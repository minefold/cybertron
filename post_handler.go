package main

import (
	"log"
	"net/http"
)

type PostHandler struct {
  Store ArchiveStore
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	part, err := firstPart(req)
	if err != nil {
		log.Println(err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

  exists, err := h.Store.Exists(req.URL.Path)
  if err != nil {
		log.Println(err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "406 archive exists", http.StatusNotAcceptable)
		return

	} else {
    defer part.Close()
		h.Store.Store(part, req.URL.Path, 0)
	}
}
