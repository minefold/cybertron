package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/minefold/cybertron/cybertron"
	"os"
)

var client *cybertron.Client

func main() {
	client = cybertron.NewClient(os.Getenv("CYBERTRON_URL"))

	list()
  // create()
}

func create() {
	// Create tar.gz archive
	archive := new(bytes.Buffer)
	gzW := gzip.NewWriter(archive)
	defer gzW.Close()

	tarW := tar.NewWriter(gzW)
	defer tarW.Close()

	// add file to archive
	body := "initial revision"
	hdr := &tar.Header{
		Name: "readme.txt",
		Size: int64(len(body)),
	}
	if err := tarW.WriteHeader(hdr); err != nil {
		panic(err)
	}
	if _, err := tarW.Write([]byte(body)); err != nil {
		panic(err)
	}

	// Upload initial archive
	if err := client.Create("/test/archive", 10, cybertron.TarGz, archive); err != nil {
		panic(err)
	}
}

func list() {
  revs, err := client.ListRevs("/test/archive", 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(revs)
}
