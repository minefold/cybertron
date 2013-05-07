package main

import (
	"io"
	"net/http"
)

type GetHandler struct {
}

func (h *GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rc, err := bucket.GetReader(req.URL.String())
	if err != nil {
		http.NotFoundHandler().ServeHTTP(w, req)
		return
	}

	defer rc.Close()
	io.Copy(w, rc)
}
