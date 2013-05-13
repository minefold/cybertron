package main

import (
	"io"
	"time"
)

type ArchiveStore interface {
	Store(archive io.Reader, url string, revision int) error
	Get(url string, revision int) (io.ReadCloser, error)
	Exists(url string) (bool, error)
	Head(url string) (*Revision, error)
	ListRevs(url string, max int) ([]Revision, error)
}

type Revision struct {
	Rev int
	At  time.Time
}
