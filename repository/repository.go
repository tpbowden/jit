package repository

import (
	"path/filepath"

	"github.com/tpbowden/jit/core"
	"github.com/tpbowden/jit/database"
	"github.com/tpbowden/jit/index"
)

type Repository struct {
	Index     *index.Index
	Workspace core.Workspace
	Database  database.Database
	Refs      core.Refs
}

func New(path string) *Repository {
	return &Repository{
		Index:     index.New(filepath.Join(path, "index")),
		Workspace: core.NewWorkspace(filepath.Dir(path)),
		Database:  database.New(filepath.Join(path, "objects")),
		Refs:      core.NewRefs(path),
	}
}
