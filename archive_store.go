package main

import (
	"io"
)

type ArchiveStore interface {
	Exists(url string) bool
	Store(archive io.Reader, url string, revision int) error
}
