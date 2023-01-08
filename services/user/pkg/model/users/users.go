package users

import (
	"gorm.io/gorm"
)

type User struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Username       string
	Firstname      string
	ProfilePicture string
}

func (User) TableName() string {
	return "users"
}

// Create new User in DB
func CreateUser(db *gorm.DB, t *User) error {
	return db.Create(t).Error
}
