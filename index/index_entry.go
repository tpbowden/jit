package index

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"os"
	"syscall"
)

type IndexFileInfo struct {
	Ctime     int32
	CtimeNsec int32
	Mtime     int32
	MtimeNsec int32
	Dev       int32
	Ino       uint32
	Mode      int32
	Uid       uint32
	Gid       uint32
	Size      int32
}

type IndexEntry struct {
	fileInfo IndexFileInfo
	oid      string
	flags    int16
	path     string
}

func (e IndexEntry) Path() string {
	return e.path
}

func (e IndexEntry) OID() string {
	return e.oid
}

func (e IndexEntry) Mode() int32 {
	return e.fileInfo.Mode
}

func (e IndexEntry) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, e.fileInfo); err != nil {
		return nil, err
	}

	hexOid, err := hex.DecodeString(e.oid)
	if err != nil {
		return nil, err
	}

	if _, err := buf.Write(hexOid); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, e.flags); err != nil {
		return nil, err
	}

	if _, err := buf.Write(append([]byte(e.path), 0)); err != nil {
		return nil, err
	}

	for buf.Len()%8 != 0 {
		if err := buf.WriteByte(0); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func NewIndexEntry(path, oid string, stat os.FileInfo) (result IndexEntry, err error) {
	info := stat.Sys().(*syscall.Stat_t)
	var flags int16
	var mode int32
	if len([]byte(path)) >= 0xfff {
		flags = 0xfff
	} else {
		flags = int16(len([]byte(path)))
	}

	if stat.Mode()&0111 == 0 {
		mode = 0100644
	} else {
		mode = 0100755
	}

	return IndexEntry{
		fileInfo: IndexFileInfo{
			Mtime:     int32(info.Mtimespec.Sec),
			MtimeNsec: int32(info.Mtimespec.Nsec),
			Ctime:     int32(info.Ctimespec.Sec),
			CtimeNsec: int32(info.Ctimespec.Nsec),
			Dev:       info.Dev,
			Ino:       uint32(info.Ino),
			Mode:      mode,
			Uid:       info.Uid,
			Gid:       info.Gid,
			Size:      int32(info.Size),
		},
		oid:   oid,
		flags: flags,
		path:  path,
	}, nil
}
