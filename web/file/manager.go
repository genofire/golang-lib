package file

import (
	"io"
	"mime/multipart"
)

type FileManager interface {
	Check(s *Service) error
	Save(s *Service, fileObj *File, file multipart.File) error
	Read(s *Service, fileObj *File) (io.Reader, error)
}

var (
	managers = make(map[string]FileManager)
)

func AddManager(typ string, m FileManager) {
	managers[typ] = m
}
