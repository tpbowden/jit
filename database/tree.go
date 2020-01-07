package database

import (
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type Node struct {
	entry *DatabaseEntry
	tree  *Tree
}

func (n Node) oid() string {
	if n.entry != nil {
		return (*n.entry).OID()
	}
	return ObjectID(n.tree)
}

func (n Node) mode() int32 {
	if n.entry != nil {
		return (*n.entry).Mode()
	}
	return n.tree.Mode()
}

type Tree struct {
	nodes map[string]Node
	order []string
}

func (t Tree) Type() string {
	return "tree"
}

func (t Tree) Mode() int32 {
	return 040000
}

func (t Tree) Data() (result []byte) {
	for _, key := range t.order {
		node := t.nodes[key]
		hash, err := hex.DecodeString(node.oid())
		if err != nil {
			panic(err)
		}
		mode := strconv.FormatInt(int64(node.mode()), 8)
		metadata := []byte(fmt.Sprintf("%s %s\x00", mode, key))
		result = append(result, metadata...)
		result = append(result, hash...)
	}
	return result
}

func (t *Tree) addNode(key string, value Node) {
	t.nodes[key] = value
	t.order = append(t.order, key)
}

func (t *Tree) addEntry(parents []string, name string, entry DatabaseEntry) {
	if len(parents) == 0 {
		t.addNode(name, Node{
			entry: &entry,
		})
	} else {
		node, found := t.nodes[parents[0]]
		if !found {
			node = Node{tree: NewTree()}
			t.addNode(parents[0], node)
		}

		node.tree.addEntry(parents[1:], name, entry)
	}
}

func (t Tree) Traverse(f func(Tree)) {
	for _, key := range t.order {
		node := t.nodes[key]
		if node.tree != nil {
			node.tree.Traverse(f)
		}
	}
	f(t)
}

func NewTree() *Tree {
	return &Tree{
		order: []string{},
		nodes: map[string]Node{},
	}
}

type DatabaseEntry interface {
	Path() string
	OID() string
	Mode() int32
}

func BuildTree(entries []DatabaseEntry) Tree {
	root := NewTree()
	for _, entry := range entries {
		path := strings.Split(entry.Path(), string(filepath.Separator))
		name, path := path[len(path)-1], path[:len(path)-1]
		root.addEntry(path, name, entry)
	}
	return *root
}
