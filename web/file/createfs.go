package file

import (
	"codeberg.org/genofire/golang-lib/web/file/fs"
	"codeberg.org/genofire/golang-lib/web/file/s3"
)

// fsType represents a type of file store.
type fsType int

const (
	typeFS fsType = iota
	typeS3
)

var stringToType = map[string]fsType{
	"fs": typeFS,
	"s3": typeS3,
}

func (t *fsType) UnmarshalText(input []byte) error {
	val, ok := stringToType[string(input)]
	if !ok {
		return ErrInvalidFSType
	}
	*t = val
	return nil
}

// FSInfo is a TOML structure storing access information about a file store.
type FSInfo struct {
	FSType fsType `config:"type" toml:"type"`
	// file system
	Root string `config:",omitempty" toml:",omitempty"`
	// s3
	Endpoint string `config:",omitempty" toml:",omitempty"`
	Secure   bool   `config:",omitempty" toml:",omitempty"`
	ID       string `config:",omitempty" toml:",omitempty"`
	Secret   string `config:",omitempty" toml:",omitempty"`
	Bucket   string `config:",omitempty" toml:",omitempty"`
	Location string `config:",omitempty" toml:",omitempty"`
}

// Create creates a file store from the information provided.
func (i *FSInfo) Create() (FS, error) {
	switch i.FSType {
	case typeFS:
		if len(i.Root) == 0 {
			return nil, ErrNoFSRoot
		}
		return &fs.FS{Root: i.Root}, nil
	case typeS3:
		return s3.New(i.Endpoint, i.Secure, i.ID, i.Secret, i.Bucket, i.Location)
	default:
		return nil, ErrNotImplementedFSType
	}
}
