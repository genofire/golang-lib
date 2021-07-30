package file

import (
	"dev.sum7.eu/genofire/golang-lib/web/file/fs"
	"dev.sum7.eu/genofire/golang-lib/web/file/s3"
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
	FSType fsType `toml:"type"`
	// file system
	Root string `toml:",omitempty"`
	// s3
	Endpoint string `toml:",omitempty"`
	Secure   bool   `toml:",omitempty"`
	ID       string `toml:",omitempty"`
	Secret   string `toml:",omitempty"`
	Bucket   string `toml:",omitempty"`
	Location string `toml:",omitempty"`
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
