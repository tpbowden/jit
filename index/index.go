package index

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tpbowden/jit/core"
)

type Index struct {
	indexPath string
	entries   map[string]IndexEntry
	order     *SortedSet
	lockfile  *core.Lockfile
	changed   bool
	parents   map[string]*Set
}

type IndexHeader struct {
	Signature [4]byte
	Version   uint32
	Entries   uint32
}

func readHeader(f *os.File) (int, error) {
	headerBytes := make([]byte, 12)
	header := &IndexHeader{}

	if _, err := f.Read(headerBytes); err != nil {
		return 0, err
	}
	if err := binary.Read(bytes.NewBuffer(headerBytes), binary.BigEndian, header); err != nil {
		return 0, err
	}

	if header.Signature != [4]byte{'D', 'I', 'R', 'C'} {
		return 0, errors.New("Invalid index header signature")
	}

	if header.Version != 2 {
		return 0, errors.New("Invalid index version")
	}

	return int(header.Entries), nil
}

func parentDirs(path string) (result []string) {
	path = filepath.Dir(path)
	for path != "." {
		result = append(result, path)
		path = filepath.Dir(path)
	}
	return result
}

func (i *Index) storeEntry(path string, e IndexEntry) {
	i.order.Add(path)
	i.entries[path] = e

	for _, dir := range parentDirs(path) {
		if _, exists := i.parents[dir]; !exists {
			i.parents[dir] = NewSet()
		}
		i.parents[dir].Add(path)
	}
}

func (i *Index) removeEntry(path string) {
	if exists := i.order.Remove(path); !exists {
		return
	}
	delete(i.entries, path)

	for _, parent := range parentDirs(path) {
		i.parents[parent].Remove(path)
	}
}

func (i *Index) discardConflicts(e IndexEntry) {
	for _, dir := range parentDirs(e.Path()) {
		i.removeEntry(dir)
	}

	parents, exists := i.parents[e.Path()]
	if exists {
		for _, parent := range parents.Entries() {
			i.removeEntry(parent)
		}
	}
}

func (i *Index) ReleaseLock() error {
	if err := i.lockfile.Rollback(); err != nil {
		return err
	}
	return nil
}

func (i *Index) Add(path, oid string, stat os.FileInfo) error {
	entry, err := NewIndexEntry(path, oid, stat)
	if err != nil {
		return err
	}
	i.discardConflicts(entry)
	i.storeEntry(path, entry)
	i.changed = true
	return nil
}

func (i Index) Entries() (entries []IndexEntry) {
	for _, path := range i.order.Entries() {
		entries = append(entries, i.entries[path])
	}
	return entries
}

func (i *Index) Load() error {
	f, err := os.OpenFile(i.indexPath, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	entries, err := readHeader(f)
	if err != nil {
		return err
	}

	for j := 0; j < entries; j++ {
		entry := make([]byte, 64)
		if _, err := f.Read(entry); err != nil {
			return err
		}

		for entry[len(entry)-1] != 0 {
			next := make([]byte, 8)
			if _, err := f.Read(next); err != nil {
				return err
			}
			entry = append(entry, next...)
		}

		info := &IndexFileInfo{}
		if err := binary.Read(bytes.NewBuffer(entry), binary.BigEndian, info); err != nil {
			return err
		}

		oid := entry[40:60]
		flags := int16(0)
		if err := binary.Read(bytes.NewBuffer(entry[60:62]), binary.BigEndian, &flags); err != nil {
			return err
		}

		path := bytes.Trim(entry[62:], "\000")
		i.storeEntry(string(path), IndexEntry{
			fileInfo: *info,
			path:     string(path),
			oid:      fmt.Sprintf("%x", oid),
			flags:    flags,
		})
	}
	data, err := i.data()
	digest := sha1.Sum(data)
	checksum := make([]byte, 20)

	if _, err := f.Read(checksum); err != nil {
		return err
	}

	if bytes.Compare(digest[:], checksum) != 0 {
		return errors.New("Checksum does not match value stored on disk")
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func (i *Index) LoadForUpdate() error {
	if err := i.lockfile.HoldForUpdate(); err != nil {
		return err
	}
	i.Clear()
	if err := i.Load(); err != nil {
		return err
	}

	return nil
}

func (i *Index) Clear() {
	i.entries = map[string]IndexEntry{}
	i.order = NewSortedSet()
	i.changed = false
	i.parents = map[string]*Set{}
}

func (i *Index) data() ([]byte, error) {
	result := new(bytes.Buffer)
	header := IndexHeader{
		Signature: [4]byte{'D', 'I', 'R', 'C'},
		Version:   2,
		Entries:   uint32(len(i.entries)),
	}
	if err := binary.Write(result, binary.BigEndian, header); err != nil {
		return nil, err
	}

	for _, entry := range i.Entries() {
		data, err := entry.Data()
		if err != nil {
			return nil, err
		}
		if _, err := result.Write(data); err != nil {
			return nil, err
		}
	}
	return result.Bytes(), nil

}

func (i *Index) WriteUpdates() error {
	if !i.changed {
		if err := i.lockfile.Rollback(); err != nil {
			return err
		}
		return nil
	}

	if err := i.lockfile.HoldForUpdate(); err != nil {
		return err
	}

	result, err := i.data()
	if err != nil {
		return err
	}
	sha := sha1.Sum(result)
	result = append(result, sha[:]...)

	if err := i.lockfile.Write(result); err != nil {
		return err
	}
	if err := i.lockfile.Commit(); err != nil {
		return err
	}

	i.changed = false

	return nil
}

func New(indexPath string) *Index {
	return &Index{
		indexPath: indexPath,
		entries:   map[string]IndexEntry{},
		lockfile:  core.NewLockfile(indexPath),
		order:     NewSortedSet(),
		changed:   false,
		parents:   map[string]*Set{},
	}
}
