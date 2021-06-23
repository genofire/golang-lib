package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User struct - default User model which could be extended
type User struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid()" example:"88078ec0-2135-445f-bf05-632701c77695"`
	Username   string     `json:"username" gorm:"unique" example:"kukoon"`
	Password   string     `json:"-" example:"super secret password"`
	ForgetCode *uuid.UUID `json:"-" gorm:"forget_code;type:uuid"`
}

// NewUser by username and password
func NewUser(username, password string) (*User, error) {
	user := &User{
		Username: username,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	return user, nil
}

// SetPassword  - create new hash of password
func (u *User) SetPassword(password string) error {
	p, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = p
	return nil
}

// ValidatePassword - check if given password is equal to saved hash
func (u *User) ValidatePassword(password string) bool {
	return ValidatePassword(u.Password, password)
}

// HasPermission interface for middleware check in other models
type HasPermission interface {
	HasPermission(tx *gorm.DB, userID, objID uuid.UUID) (interface{}, error)
}
