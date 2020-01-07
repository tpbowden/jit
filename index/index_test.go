package index_test

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/tpbowden/jit/index"
)

var indexFile, tempFile, tempDir string
var stat os.FileInfo

func setup() error {
	tempDir, err := ioutil.TempDir("", "jit")
	if err != nil {
		return err
	}
	indexFile = filepath.Join(tempDir, "index")
	tempFile = filepath.Join(tempDir, "test.go")
	if err := ioutil.WriteFile(tempFile, []byte("some data"), 0644); err != nil {
		return err
	}
	stat, err = os.Stat(tempFile)
	if err != nil {
		return err
	}
	return nil
}

func cleanup() {
	os.RemoveAll(tempDir)
}

func sha() string {
	sha := sha1.Sum([]byte("some hash"))
	return fmt.Sprintf("%x", sha)
}

func compareFileList(t *testing.T, i *index.Index, expected []string) {
	entries := i.Entries()
	if len(entries) != len(expected) {
		t.Fatalf("Number of entries does not match. Expected %d, got %d", len(expected), len(entries))
	}
	for i, entry := range entries {
		if expected[i] != entry.Path() {
			t.Errorf("Paths do not match. Expected %s, got %s", expected[i], entry.Path())
		}
	}
}

func TestAddingAFile(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	i := index.New(indexFile)
	i.Add("test.go", sha(), stat)

	expected := []string{"test.go"}
	compareFileList(t, i, expected)
}

func TestReplacingAFiledWithADirectory(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	i := index.New(indexFile)
	i.Add("alice.txt", sha(), stat)
	i.Add("bob.txt", sha(), stat)
	i.Add("alice.txt/nested.txt", sha(), stat)

	expected := []string{"alice.txt/nested.txt", "bob.txt"}
	compareFileList(t, i, expected)
}

func TestReplacingADirectoryWithAFile(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	i := index.New(indexFile)
	i.Add("alice.txt", sha(), stat)
	i.Add("nested/bob.txt", sha(), stat)
	i.Add("nested", sha(), stat)

	expected := []string{"alice.txt", "nested"}
	compareFileList(t, i, expected)
}

func TestRecursivelyReplacingADirectoryWithAFile(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	i := index.New(indexFile)
	i.Add("alice.txt", sha(), stat)
	i.Add("nested/bob.txt", sha(), stat)
	i.Add("nested/inner/claire.txt", sha(), stat)
	i.Add("nested", sha(), stat)

	expected := []string{"alice.txt", "nested"}
	compareFileList(t, i, expected)
}
