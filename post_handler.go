package main

import (
	"fmt"
	"github.com/whatupdave/s3/s3util"
	"io"
	// "launchpad.net/goamz/s3"
	// "github.com/kr/pretty"
	"log"
	"net/http"
	"strings"
)

type PostHandler struct {
}

func archiveExists(url string) bool {
	prefix := strings.TrimLeft(url, "/") // strip off leading / to get prefix

	list, err := s3util.List(s3url, prefix, "", 1, nil)
	if err != nil {
		log.Println(err)
		return false
	}
	return len(list.Contents) > 0
}

func uploadPart(part io.ReadCloser, url string, revision int) {
	uploader, err := s3util.Create(fmt.Sprintf("%s.%d", url, revision), nil, nil)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = io.Copy(uploader, part)
	if err != nil {
		log.Println(err)
		return
	}

	defer uploader.Close()
	defer part.Close()
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
		}

		fmt.Printf("%s => %s\n", part.FileName(), req.URL.String())

		url := s3url + req.URL.String()
		if archiveExists(req.URL.String()) {
			http.Error(w, "406 archive exists", http.StatusNotAcceptable)
			return

		} else {
			uploadPart(part, url, 0)
		}
	}
}
