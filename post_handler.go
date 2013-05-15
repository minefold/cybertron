package main

import (
	"log"
	"net/http"
	"strconv"
)

type PostHandler struct {
	Store ArchiveStore
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rev, err := strconv.Atoi(req.URL.Query().Get("rev"))
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
		log.Println("Store part")
		err := h.Store.Store(part, req.URL.Path, rev)
		if err != nil {
			log.Println(err)
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
		log.Println("Stored part")
	}
}
