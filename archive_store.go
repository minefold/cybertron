package main

import (
	"io"
	"time"
)

type ArchiveStore interface {
	Exists(url string) bool
	Store(archive io.Reader, url string, revision int) error
	ListRevs(url string, max int) ([]Revision, error)
}

type Revision struct {
	Rev int
	At  time.Time
}
