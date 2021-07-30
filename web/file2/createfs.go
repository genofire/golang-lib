package file

import (
	"errors"

	"dev.sum7.eu/genofire/golang-lib/web/file2/fs"
	"dev.sum7.eu/genofire/golang-lib/web/file2/s3"
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
		return errors.New("invalid file store type")
	}
	*t = val
	return nil
}

// FSInfo is a TOML structure storing access information about a file store.
type FSInfo struct {
	fstype fsType `toml:"type"`
	// file system
	root string `toml:",omitempty"`
	// s3
	endpoint string `toml:",omitempty"`
	secure   bool   `toml:",omitempty"`
	id       string `toml:",omitempty"`
	secret   string `toml:",omitempty"`
	bucket   string `toml:",omitempty"`
	location string `toml:",omitempty"`
}

// Create creates a file store from the information provided.
func (i *FSInfo) Create() (FS, error) {
	switch i.fstype {
	case typeFS:
		if len(i.root) == 0 {
			return nil, errors.New("no file store root")
		}
		return &fs.FS{Root: i.root}, nil
	case typeS3:
		return s3.New(i.endpoint, i.secure, i.id, i.secret, i.bucket, i.location)
	default:
		return nil, errors.New("FSInfo.Create not implemented for provided file store type")
	}
}
