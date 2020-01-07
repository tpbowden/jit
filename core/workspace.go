package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Workspace struct {
	rootDir string
}

func (w Workspace) doListFiles(root string, fileNames []string) ([]string, error) {
	stat, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, missingFile(root)
		}
		return nil, err
	}

	if !stat.IsDir() {
		relative, err := filepath.Rel(w.rootDir, root)
		if err != nil {
			return nil, err
		}
		fileNames = append(fileNames, relative)
		return fileNames, nil
	}

	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err

	}

	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}
		if file.IsDir() {
			fileNames, err = w.doListFiles(filepath.Join(root, file.Name()), fileNames)
			if err != nil {
				return nil, err
			}
			continue
		}
		relative, err := filepath.Rel(w.rootDir, filepath.Join(root, file.Name()))
		if err != nil {
			return nil, err
		}
		fileNames = append(fileNames, relative)
	}

	return fileNames, err
}

func (w Workspace) ListFilesRelative(root string) (fileNames []string, err error) {
	return w.doListFiles(root, fileNames)
}

func (w Workspace) ListFiles() (fileNames []string, err error) {
	return w.doListFiles(w.rootDir, fileNames)
}

func (w Workspace) ReadFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Join(w.rootDir, path))
	if err != nil {
		if os.IsPermission(err) {
			return nil, noPermission(path)
		}
		return nil, err
	}
	return data, nil

}

func (w Workspace) StatFile(path string) (os.FileInfo, error) {
	info, err := os.Stat(filepath.Join(w.rootDir, path))
	if err != nil {
		if os.IsPermission(err) {
			return nil, noPermission(path)
		}
		return nil, err
	}
	return info, nil
}

func NewWorkspace(rootDir string) Workspace {
	return Workspace{
		rootDir: rootDir,
	}
}
