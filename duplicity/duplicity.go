package duplicity

import (
	"os/exec"
)

const TimeLayout = "Mon Jan 2 15:04:05 2006"

type Duplicity struct{}

func (d *Duplicity) Full(src, repo string) (*OpStats, error) {
	out, err := cmd("full", src, repo).CombinedOutput()
	if err != nil {
		return nil, err
	}
	return DecodeOpStats(out)
}

func (d *Duplicity) Incr(src, repo string) (*OpStats, error) {
	out, err := cmd("incr", src, repo).CombinedOutput()
	if err != nil {
		return nil, err
	}
	return DecodeOpStats(out)
}

func (d *Duplicity) Status(repo string) (*CollStats, error) {
	out, err := cmd("collection-status", repo).CombinedOutput()
	if err != nil {
		return nil, err
	}
	return DecodeCollStats(out)
}

func (d *Duplicity) Files(repo string) (Files, error) {
	out, err := cmd("list-current-files", repo).CombinedOutput()
	if err != nil {
		return nil, err
	}
	return DecodeFiles(out)
}

func cmd(arg ...string) *exec.Cmd {
	args := append([]string{"--no-encryption"}, arg...)
	return exec.Command("duplicity", args...)
}
