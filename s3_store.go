package main

import (
	"fmt"
	"github.com/whatupdave/s3/s3util"
	"io"
	"log"
	"strings"
)

type S3Store struct {
	BaseUrl string // eg. https://s3.amazonaws.com/cybertron
}

func NewS3Store(baseUrl string) *S3Store {
	return &S3Store{BaseUrl: baseUrl}
}

func (s3 *S3Store) Exists(url string) bool {
	prefix := strings.TrimLeft(url, "/") // strip off leading / to get prefix

	list, err := s3util.List(s3.BaseUrl, prefix, "", 1, nil)
	if err != nil {
		log.Println(err)
		return false
	}
	return len(list.Contents) > 0
}

func (s3 *S3Store) Store(archive io.Reader, url string, revision int) error {
	uploader, err := s3util.Create(fmt.Sprintf("%s.%d", s3.BaseUrl+url, revision), nil, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(uploader, archive)
	if err != nil {
		return err
	}

	defer uploader.Close()

	return nil
}
