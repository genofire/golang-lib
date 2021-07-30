package file

import "errors"

// errors
var (
	ErrInvalidFSType        = errors.New("invalid file store type")
	ErrNoFSRoot             = errors.New("no file store root")
	ErrNotImplementedFSType = errors.New("FSInfo.Create not implemented for provided file store type")
)
