package command

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/tpbowden/jit/database"
	"github.com/tpbowden/jit/repository"
)

func (c *Command) cmdCommit() (int, error) {
	dir := c.Dir
	gitDir := filepath.Join(dir, ".git")
	repo := repository.New(gitDir)
	if err := repo.Index.Load(); err != nil {
		return 1, err
	}
	entries := repo.Index.Entries()
	dbEntries := make([]database.DatabaseEntry, len(entries))

	for i, entry := range entries {
		dbEntries[i] = entry
	}

	tree := database.BuildTree(dbEntries)
	tree.Traverse(func(t database.Tree) {
		if err := repo.Database.Store(t); err != nil {
			panic(err)
		}
	})
	authorName := c.Env["GIT_AUTHOR_NAME"]
	authorEmail := c.Env["GIT_AUTHOR_EMAIL"]
	author := database.NewAuthor(authorName, authorEmail)
	message, err := ioutil.ReadAll(c.Stdin)
	if err != nil {
		return 1, err
	}
	parent, err := repo.Refs.ReadHead()
	if err != nil {
		return 1, err
	}
	commit := database.NewCommit(
		author,
		database.ObjectID(tree),
		parent,
		string(message),
		time.Now(),
	)
	if err := repo.Database.Store(commit); err != nil {
		return 1, err
	}

	if err := repo.Refs.UpdateHead(database.ObjectID(commit)); err != nil {
		return 1, err
	}

	isRoot := ""
	if parent == "" {
		isRoot = "(root-commit) "
	}

	fmt.Fprintf(
		c.Stdout,
		"[%s%s] %s\n",
		isRoot,
		database.ObjectID(commit),
		strings.Split(string(commit.Message), "\n")[0],
	)
	return 0, nil
}
