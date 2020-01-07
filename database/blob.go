package database

type Blob struct {
	data []byte
}

func (b Blob) Type() string {
	return "blob"
}

func (b Blob) Data() []byte {
	return b.data
}

func NewBlob(data []byte) Blob {
	return Blob{data}
}
