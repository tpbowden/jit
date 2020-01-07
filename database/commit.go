package database

import (
	"fmt"
	"time"
)

type Commit struct {
	Author    Author
	TreeID    string
	ParentID  string
	Message   string
	Timestamp time.Time
}

func (c Commit) Type() string {
	return "commit"
}

func (c Commit) Data() (result []byte) {
	authorString := fmt.Sprintf(
		"%s <%s> %d +0000",
		c.Author.name,
		c.Author.email,
		c.Timestamp.Unix(),
	)
	tree := fmt.Sprintf("tree %s\n", c.TreeID)
	parent := fmt.Sprintf("parent %s\n", c.ParentID)
	author := fmt.Sprintf("author %s\n", authorString)
	committer := fmt.Sprintf("committer %s\n", authorString)
	message := fmt.Sprintf("\n%s", c.Message)

	result = append(result, tree...)
	if c.ParentID != "" {
		result = append(result, parent...)
	}
	result = append(result, author...)
	result = append(result, committer...)
	result = append(result, message...)

	return result
}

func NewCommit(
	author Author,
	treeID string,
	parentID string,
	message string,
	timestamp time.Time,
) Commit {
	return Commit{
		Author:    author,
		TreeID:    treeID,
		ParentID:  parentID,
		Message:   message,
		Timestamp: timestamp,
	}
}
