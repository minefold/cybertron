package duplicity

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestCRUD(t *testing.T) {
	prepare("tmp/testing")
	createFile("tmp/testing/src/a (crazy) 'readme' [file].txt", "oh hai!")
  defer remove("tmp/testing")

	// full update
	dup := new(Duplicity)
	stats, _ := dup.Full("tmp/testing/src", "file://tmp/testing/repo")
	if stats.NewFiles != 2 {
		t.Fatalf("Expected 2 was %d\n", stats.NewFiles)
	}

	// new file update
	createFile("tmp/testing/src/2.txt", "sup")
	dup.Incr("tmp/testing/src", "file://tmp/testing/repo")

  // list backups
  coll, _ := dup.Status("file://tmp/testing/repo")
  if len(coll.Backups) != 2 {
    t.Fatalf("Expected 2 was %d\n", len(coll.Backups))
  }

  // list files
  files, _ := dup.Files("file://tmp/testing/repo")
  if len(files) != 3 {
    t.Fatalf("Expected 3 was %d\n", len(files))
  }
}

func prepare(path string) {
	cmd := fmt.Sprintf("mkdir -p %s && cd %s && rm -rf * && mkdir src dst repo", path, path)
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		fmt.Println(cmd, string(out))
		panic(err)
	}
}

func createFile(path, contents string) {
	if err := ioutil.WriteFile(path, []byte(contents), 0777); err != nil {
		panic(err)
	}
}

func remove(path string) {
	err := exec.Command("sh", "-c", fmt.Sprintf("rm -rf %s", path)).Run()
	if err != nil {
		panic(err)
	}
}
