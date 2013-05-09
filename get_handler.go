package main

import (
	"fmt"
	"github.com/whatupdave/s3/s3util"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Revision struct {
	Rev int
	At  time.Time
}

type GetHandler struct {
}

func getRevisions(url string) ([]Revision, error) {
	prefix := strings.TrimLeft(url, "/") // strip off leading / to get prefix

	list, err := s3util.List(s3url, prefix, "", 1000, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	revs := make([]Revision, 0)
	for _, key := range list.Contents {
		parts := strings.Split(key.Key, ".")

		// file := strings.Join(parts[:len(parts)-1], ".")
		rev, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			log.Println("Bad revision number", key.Key)
			continue
		}

		at, err := time.Parse(time.RFC3339, key.LastModified)
		if err != nil {
			log.Println("Bad timestamp", key.Key, key.LastModified)
			continue
		}

		revs = append(revs, Revision{
			Rev: rev,
			At:  at,
		})
	}

	return revs, err
}

func (h *GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// check for .tar.gz extension

	parts := strings.Split(req.URL.String(), ".tar.gz")
	if len(parts) < 2 {
		http.NotFound(w, req)
		return
	}

	revs, err := getRevisions(parts[0])
	if err != nil {
		log.Println(err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	log.Println(revs)
	rev := revs[len(revs)-1]

	rc, err := s3util.Open(s3url+fmt.Sprintf("%s.%d", parts[0], rev.Rev), nil)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	defer rc.Close()
	io.Copy(w, rc)
}
