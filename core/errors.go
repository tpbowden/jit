package core

import "fmt"

type MissingFile struct {
	path string
}

func (e *MissingFile) Error() string {
	return fmt.Sprintf("pathspec '%s' did not match any files", e.path)
}

func missingFile(path string) error {
	return &MissingFile{path}
}

type NoPermission struct {
	path string
}

func (e *NoPermission) Error() string {
	return fmt.Sprintf("permission denied: '%s'", e.path)
}

func noPermission(path string) error {
	return &NoPermission{path}
}

type LockDenied struct {
	path string
}

func (e *LockDenied) Error() string {
	return fmt.Sprintf("Lock denied: '%s' exists", e.path)
}

func lockDenied(path string) error {
	return &LockDenied{path}
}
