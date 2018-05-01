package data

import (
	"time"

	"github.com/jinzhu/gorm"
	//
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// UserEntity ..
type UserEntity struct {
	gorm.Model
	UUID            string `gorm:"type:uuid;primary_key;DEFAULT:uuid_generate_v4()"`
	Email           string `gorm:"unique"`
	PasswordHash    string `gorm:"type:varchar(128); not null"`
	FullName        string `gorm:"type:varchar(72); not null"`
	IsEmailVerified bool   `gorm:"type:boolean; default:false"`
	CreatedAt       time.Time
	LastModifiedAt  time.Time
}
