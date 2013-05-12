package main

import (
	"io"
	// "launchpad.net/goamz/s3"
	// "github.com/kr/pretty"
	"log"
	"net/http"
)

type PostHandler struct {
  Store ArchiveStore
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mr, err := req.MultipartReader()
	if err != nil {
		return
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
		  log.Println(err)
      continue
		}

		if h.Store.Exists(req.URL.String()) {
			http.Error(w, "406 archive exists", http.StatusNotAcceptable)
			return

		} else {
      defer part.Close()
			h.Store.Store(part, req.URL.String(), 0)
		}
	}
}
