package main

import (
	"fmt"
	"net/http"
)

type Router struct {
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method, req.URL)

  switch req.Method {
  case "GET":
    new(GetHandler).ServeHTTP(w, req)

  case "POST":
    new(PostHandler).ServeHTTP(w, req)
  }
}