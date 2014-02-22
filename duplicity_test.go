package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

type Duplicity struct{}

func (d *Duplicity) Full(src, dst string) (*DupStats, error) {
	cmd := exec.Command("duplicity",
		"--no-encryption",
		"full",
		src,
		dst)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return DecodeStats(out)
}

type DupStats struct {
	NewFiles int
}

func DecodeStats(info []byte) (*DupStats, error) {
	stats := new(DupStats)
	scanner := bufio.NewScanner(bytes.NewReader(info))
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return stats, nil
}

func TestCreate(t *testing.T) {
	prepare("tmp/testing")
	defer remove("tmp/testing")

	dup := new(Duplicity)
	stats, err := dup.Full("tmp/testing/src", "file://tmp/testing/repo")
	if err != nil {
		t.Error(err)
	}
	if stats.NewFiles != 2 {
		t.Fatalf("Expected 2 was %d\n", stats.NewFiles)
	}
}

func prepare(path string) {
	exec.Command("sh", "-c", fmt.Sprintf("mkdir %s && cd %s && mkdir src dst repo", path)).Run()
}

func remove(path string) {
	exec.Command("sh", "-c", fmt.Sprintf("rm -rf %s", path)).Run()
}
