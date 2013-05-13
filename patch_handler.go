package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
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

func firstPart(req *http.Request) (io.ReadCloser, error) {
	mr, err := req.MultipartReader()
	if err != nil {
		return nil, err
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}

		return part, nil
	}

	return nil, nil
}

func getRevParams(req *http.Request) (int, int, error) {
	fromRev, err := strconv.Atoi(req.URL.Query().Get("rev"))
	if err != nil {
		return 0, 0, err
	}
	toRev, err := strconv.Atoi(req.Header.Get("X-Rev"))
	if err != nil {
		return 0, 0, err
	}

	if fromRev >= toRev {
		return fromRev, toRev, errors.New("rev not newer")
	}

	return fromRev, toRev, nil
}
