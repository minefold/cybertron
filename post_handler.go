package main

import (
	"fmt"
	"io"
	// "launchpad.net/goamz/aws"
	"bytes"
	"launchpad.net/goamz/s3"
	"log"
	"net/http"
)

type PostHandler struct {
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
		fmt.Println(part.Header)
		fmt.Printf("%s => %s\n", part.FileName(), req.URL.String())

		// read file into a buffer so we can see how big it is. S3 needs to know the content-length
		// this can be made more efficient by using S3s multipart uploads and uploading chunks of
		// data
		var b bytes.Buffer
		length, err := io.Copy(&b, part)
		if err != nil {
			log.Println("Error copying to buffer", err)
		}

		log.Printf("Copied %d bytes", length)

		// upload buffer to S3
		err = bucket.PutReader(req.URL.String(), &b, length, "", s3.Private)
		if err != nil {
			log.Println(err)
		}
	}

}
