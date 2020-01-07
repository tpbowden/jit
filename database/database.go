package database

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Database struct {
	dbPath string
}

type PersistableObject interface {
	Type() string
	Data() []byte
}

func ObjectContent(object PersistableObject) string {
	return fmt.Sprintf("%s %d\x00%s", object.Type(), len(object.Data()), object.Data())

}

func ObjectID(object PersistableObject) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(ObjectContent(object))))
}

func (db Database) Store(object PersistableObject) error {
	hash := ObjectID(object)
	content := ObjectContent(object)
	path := filepath.Join(db.dbPath, hash[0:2], hash[2:])
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return nil
	}

	dirname := filepath.Dir(path)

	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
		log.Print("Failed to create object directory")
		return err
	}

	tmpFile, err := ioutil.TempFile(dirname, "tmp_object_*")
	if err != nil {
		log.Print("Failed to create temporary object")
		return err
	}

	var b bytes.Buffer
	w, err := zlib.NewWriterLevel(&b, zlib.BestSpeed)
	if err != nil {
		log.Print("Failed to compress object")
		return err
	}

	w.Write([]byte(content))
	w.Close()

	if _, err := tmpFile.Write(b.Bytes()); err != nil {
		log.Print("Failed to write temporary object")
		return err
	}

	if err := tmpFile.Close(); err != nil {
		log.Print("Failed to close temporary object")
		return err
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		log.Print("Failed to rename temporary object")
		return err
	}

	return nil
}

func New(dbPath string) Database {
	return Database{dbPath}
}
