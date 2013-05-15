package main

import (
	"io"
	"time"
)

type ArchiveStore interface {
	Exists(url string) (bool, error)
	Revs(url string, max int) ([]Revision, error)
	Get(url string, rev int) (io.ReadCloser, error)
	Store(archive io.Reader, url string, rev int) error
  Del(url string, rev int) error
}

type Revision struct {
	Rev int
	At  time.Time
}
