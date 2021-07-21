package file

import "github.com/google/uuid"

// File to store information in database
type File struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid()" example:"32466d63-efa4-4f27-9f2b-a1f06c8e2e1d"`
	StorageType string    `json:"storage_type,omitempty"`
	Path        string    `json:"path,omitempty"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content-type"`
}
