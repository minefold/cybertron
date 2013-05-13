package main

import (
	"fmt"
	"net/http"
)

type Router struct {
	Store ArchiveStore
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method, req.URL.String())

	handler := r.getHandler(req)
	if handler != nil {
		handler.ServeHTTP(w, req)
	}
}

func (r *Router) getHandler(req *http.Request) http.Handler {
	switch req.Method {
	case "GET":
		return &GetHandler{Store: r.Store}

  case "POST":
    return &PostHandler{Store: r.Store}

  case "PATCH":
    return &PatchHandler{Store: r.Store}
	}
	return nil
}
