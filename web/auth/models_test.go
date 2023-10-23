package auth

import (
	"testing"

	gormigrate "github.com/genofire/gormigrate/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"codeberg.org/genofire/golang-lib/database"

	"github.com/stretchr/testify/assert"
)

var (
	TestUser1ID = uuid.MustParse("88078ec0-2135-445f-bf05-632701c77695")
)

func SetupMigration(db *database.Database) {
	db.AddMigration([]*gormigrate.Migration{
		{
			ID: "01-schema-0008-01-user",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
		{
			ID: "10-data-0008-01-user",
			Migrate: func(tx *gorm.DB) error {
				PasswordHashCost = bcrypt.MinCost
				user, err := NewUser("admin", "CHANGEME")
				if err != nil {
					return err
				}
				user.ID = TestUser1ID
				return tx.Create(user).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Delete(&User{
					ID: TestUser1ID,
				}).Error
			},
		},
	}...)
}

func TestUserPassword(t *testing.T) {
	assert := assert.New(t)
	password := "password"
	user, err := NewUser("admin", password)

	assert.Nil(err)
	assert.NotNil(user)

	assert.False(user.ValidatePassword("12346"))
	assert.True(user.ValidatePassword(password))
	assert.NotEqual(password, user.Password, "password should be hashed")
}
