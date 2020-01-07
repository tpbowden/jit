package command

import (
	"fmt"
	"os"
	"path/filepath"
)

func (c *Command) cmdInit() (int, error) {
	dir := c.Dir
	if len(c.Args) >= 3 {
		var err error
		dir, err = filepath.Abs(c.Args[2])
		if err != nil {
			return 1, err
		}
	}

	gitDir := filepath.Join(dir, ".git")

	if err := os.MkdirAll(gitDir, os.ModePerm); err != nil {
		return 1, err
	}

	refsDir := filepath.Join(gitDir, "refs")
	if err := os.MkdirAll(refsDir, os.ModePerm); err != nil {
		return 1, err
	}

	objectsDir := filepath.Join(gitDir, "objects")
	if err := os.MkdirAll(objectsDir, os.ModePerm); err != nil {
		return 1, err
	}

	fmt.Fprintln(c.Stdout, "Initialised empty Jit repository in", gitDir)
	return 0, nil
}
