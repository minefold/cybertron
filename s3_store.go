package main

import (
	"fmt"
	"github.com/whatupdave/s3/s3util"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

type S3Store struct {
	BaseUrl string // eg. https://s3.amazonaws.com/cybertron
}

func NewS3Store(baseUrl string) *S3Store {
	return &S3Store{BaseUrl: baseUrl}
}

func (s3 *S3Store) Exists(url string) (bool, error) {
	revs, err := s3.ListRevs(url, 1)
	if err != nil {
		return false, err
	}

	return len(revs) > 0, nil
}

func (s3 *S3Store) Head(url string) (*Revision, error) {
	revs, err := s3.ListRevs(url, 0)
	if err != nil {
		return nil, err
	}

	if len(revs) == 0 {
		return nil, nil
	}

	return &(revs[len(revs)-1]), nil
}

func (s3 *S3Store) ListRevs(url string, max int) ([]Revision, error) {
	prefix := strings.TrimLeft(url, "/") // strip off leading / to get prefix

	if max < 1 {
		max = 1000
	}

	list, err := s3util.List(s3.BaseUrl, prefix, "", max, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	revs := make([]Revision, 0)
	for _, key := range list.Contents {
		parts := strings.Split(key.Key, "/")
		revPart := strings.TrimRight(parts[len(parts)-1], ".tar.gz")

		rev, err := strconv.Atoi(revPart)
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

func (s3 *S3Store) Store(archive io.Reader, url string, rev int) error {
	key := s3.key(url, rev)
	uploader, err := s3util.Create(key, nil, nil)
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

func (s3 *S3Store) Get(url string, rev int) (io.ReadCloser, error) {
	rc, err := s3util.Open(s3.key(url, rev), nil)
	if err != nil {
		return nil, err
	}

	return rc, nil
}

func (s3 *S3Store) key(url string, rev int) string {
	return s3.BaseUrl + fmt.Sprintf("%s/%d.tar.gz", url, rev)
}
