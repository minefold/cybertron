package main

import (
	"log"
	"net/http"
	"strconv"
)

type DeleteHandler struct {
	Store ArchiveStore
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rev, err := strconv.Atoi(req.URL.Query().Get("rev"))
	if err != nil {
		log.Println(err)
		http.Error(w, "400", http.StatusBadRequest)
		return
	}

	err = h.Store.Del(req.URL.Path, rev)
	if err != nil {
		log.Println(err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
}
