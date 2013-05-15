package main

import (
	"fmt"
	"github.com/minefold/cybertron/cybertron"
	"io"
	"os"
	"os/exec"
)

var client *cybertron.Client

func main() {
	client = cybertron.NewClient(os.Getenv("CYBERTRON_URL"))

	local := os.Args[1]
	remote := os.Args[2]

	fmt.Println(local, remote)

	list(remote)
	create(local, remote, 1)
	list(remote)
  update(local, remote, 1, 2)
  list(remote)
  delete(remote, 1)
  list(remote)
  delete(remote, 2)
  list(remote)
}

func list(remote string) {
	revs, err := client.List(remote, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("revisions:", revs)
}

func create(local, remote string, rev int) {
	tar, err := tarGz(local)
	if err != nil {
		panic(err)
	}

	if err := client.Create(remote, rev, cybertron.TarGz, tar); err != nil {
		panic(err)
	}
	fmt.Println("created initial revision", rev)
}

func update(local, remote string, from, to int) {
  tar, err := tarGz(local)
  if err != nil {
    panic(err)
  }

  if err := client.Update(remote, from, to, cybertron.TarGz, tar); err != nil {
    panic(err)
  }
  fmt.Println("updated", from, "to", to)
}

func delete(remote string, rev int) {
  if err := client.Delete(remote, rev); err != nil {
    panic(err)
  }
  fmt.Println("deleted", rev)
}

func tarGz(path string) (io.Reader, error) {
	tar := exec.Command("tar", "-cz", path)
	stdout, err := tar.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = tar.Start()
	if err != nil {
		return nil, err
	}

	go tar.Wait()

	return stdout, nil
}
