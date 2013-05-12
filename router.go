package main

import (
	"fmt"
	"net/http"
)

type Router struct {
	Store ArchiveStore
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method, req.URL)

	switch req.Method {
	case "GET":
		new(GetHandler).ServeHTTP(w, req)

	case "POST":
		handler := &PostHandler{Store: r.Store}
		handler.ServeHTTP(w, req)
	}
}
