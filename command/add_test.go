package command_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tpbowden/jit/command"
	"github.com/tpbowden/jit/repository"
)

type TestHelper struct {
	path string
	t    *testing.T
	repo *repository.Repository
	cmd  *command.Command
}

func NewTestHelper(test *testing.T) *TestHelper {
	path, err := ioutil.TempDir("", "jit_test")
	if err != nil {
		test.Fatal(err)
	}

	var stdout, stderr strings.Builder

	repo := repository.New(filepath.Join(path, ".git"))
	cmd := &command.Command{
		Dir:    path,
		Args:   os.Args,
		Env:    map[string]string{},
		Stderr: &stderr,
		Stdout: &stdout,
		Stdin:  strings.NewReader(""),
	}
	return &TestHelper{
		path: path,
		repo: repo,
		t:    test,
		cmd:  cmd,
	}
}

func TestAddingAFileToTheIndex(t *testing.T) {
	helper := NewTestHelper(t)
}
