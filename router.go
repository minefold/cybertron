package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"
	"strconv"
	"strings"
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

	fmt.Println(req.Method, req.URL.String(), "DONE")
}

func (r *Router) getHandler(req *http.Request) http.Handler {
	url := req.URL.String()
	switch {
	case strings.HasPrefix(url, "/debug/pprof/cmdline"):
		return http.HandlerFunc(pprof.Cmdline)
	case strings.HasPrefix(url, "/debug/pprof/profile"):
		return http.HandlerFunc(pprof.Profile)
	case strings.HasPrefix(url, "/debug/pprof/"):
		return http.HandlerFunc(pprof.Index)
	case strings.HasPrefix(url, "/debug/pprof/symbol"):
		return http.HandlerFunc(pprof.Symbol)
	}

	switch req.Method {
	case "GET":
		return &GetHandler{Store: r.Store}

	case "POST":
		return &PostHandler{Store: r.Store}

	case "PATCH":
		return &PatchHandler{Store: r.Store}

	case "DELETE":
		return &DeleteHandler{Store: r.Store}

	}
	return nil
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
