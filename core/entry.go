package core

import "os"

type Entry struct {
	Name string
	Oid  string
	Stat os.FileInfo
}

func (e Entry) Mode() string {
	if e.Stat.Mode()&0111 == 0 {
		return "100644"
	} else {
		return "100755"
	}
}

func NewEntry(name string, oid string, stat os.FileInfo) Entry {
	return Entry{
		Name: name,
		Oid:  oid,
		Stat: stat,
	}
}
