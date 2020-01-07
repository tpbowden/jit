package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Refs struct {
	gitDir string
}

func (r Refs) UpdateHead(oid string) error {
	head := filepath.Join(r.gitDir, "HEAD")
	lockfile := NewLockfile(head)
	if err := lockfile.HoldForUpdate(); err != nil {
		return err
	}

	content := fmt.Sprintf("%s\n", oid)
	if err := lockfile.Write([]byte(content)); err != nil {
		return err
	}

	if err := lockfile.Commit(); err != nil {
		return err
	}

	return nil
}

func (r Refs) ReadHead() (string, error) {
	head := filepath.Join(r.gitDir, "HEAD")

	oid, err := ioutil.ReadFile(head)
	if os.IsNotExist(err) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(oid)), nil
}

func NewRefs(gitDir string) Refs {
	return Refs{
		gitDir: gitDir,
	}
}
