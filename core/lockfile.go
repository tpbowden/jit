package core

import (
	"fmt"
	"os"
)

type Lockfile struct {
	filePath string
	lockPath string
	lock     *os.File
}

func (l *Lockfile) HoldForUpdate() error {
	if l.lock != nil {
		return nil
	}
	lockfile, err := os.OpenFile(l.lockPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return lockDenied(l.lockPath)
		}
		if os.IsPermission(err) {
			return noPermission(l.lockPath)
		}
		return err
	}
	l.lock = lockfile
	return nil
}

func (l *Lockfile) Rollback() error {
	if l.lock == nil {
		return fmt.Errorf("No lock held for '%s'", l.filePath)
	}
	l.lock.Close()
	if err := os.Remove(l.lockPath); err != nil {
		return err
	}
	l.lock = nil
	return nil
}

func (l *Lockfile) Write(data []byte) error {
	if l.lock == nil {
		return fmt.Errorf("No lock held for '%s'", l.filePath)
	}
	if _, err := l.lock.Write(data); err != nil {
		return err
	}

	return nil
}

func (l *Lockfile) Commit() error {
	if l.lock == nil {
		return fmt.Errorf("No lock held for '%s'", l.filePath)
	}
	if err := l.lock.Close(); err != nil {
		return err
	}
	if err := os.Rename(l.lockPath, l.filePath); err != nil {
		return err
	}
	l.lock = nil

	return nil
}

func NewLockfile(path string) *Lockfile {
	return &Lockfile{
		filePath: path,
		lockPath: fmt.Sprintf("%s.lock", path),
	}
}
