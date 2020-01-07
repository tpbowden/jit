package command

import (
	"fmt"
	"path/filepath"

	"github.com/tpbowden/jit/core"
	"github.com/tpbowden/jit/database"
	"github.com/tpbowden/jit/repository"
)

func (c *Command) cmdAdd() (int, error) {
	dir := c.Dir

	if len(c.Args) < 3 {
		return 1, fmt.Errorf("No file path supplied to add")
	}

	gitDir := filepath.Join(dir, ".git")
	repo := repository.New(gitDir)
	if err := repo.Index.LoadForUpdate(); err != nil {
		if ld, ok := err.(*core.LockDenied); ok {
			fmt.Fprintln(c.Stderr, "fatal:", ld.Error())
			return 128, nil
		}
		return 1, err
	}

	paths := []string{}
	for _, root := range c.Args[2:] {
		root, err := filepath.Abs(root)
		if err != nil {
			return 1, err
		}
		p, err := repo.Workspace.ListFilesRelative(root)
		if err != nil {
			if mf, ok := err.(*core.MissingFile); ok {
				if err := repo.Index.ReleaseLock(); err != nil {
					return 1, err
				}
				fmt.Fprintln(c.Stdout, "fatal:", mf.Error())
				return 128, nil
			}
			return 1, err
		}
		paths = append(paths, p...)
	}

	for _, path := range paths {
		data, err := repo.Workspace.ReadFile(path)
		if err != nil {
			if npe, ok := err.(*core.NoPermission); ok {
				if err := repo.Index.ReleaseLock(); err != nil {
					return 1, err
				}
				fmt.Fprintln(c.Stdout, "fatal:", npe.Error())
				return 128, nil
			}
			return 1, err
		}
		stat, err := repo.Workspace.StatFile(path)
		if err != nil {
			if npe, ok := err.(*core.NoPermission); ok {
				if err := repo.Index.ReleaseLock(); err != nil {
					return 1, err
				}
				fmt.Fprintln(c.Stdout, "fatal:", npe.Error())
				return 128, nil
			}
			return 1, err
		}
		blob := database.NewBlob(data)
		if err := repo.Database.Store(blob); err != nil {
			return 1, err
		}

		repo.Index.Add(path, database.ObjectID(blob), stat)
	}

	if err := repo.Index.WriteUpdates(); err != nil {
		return 1, err
	}
	return 0, nil
}
