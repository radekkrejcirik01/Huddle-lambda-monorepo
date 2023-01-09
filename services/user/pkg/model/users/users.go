package users

import (
	"gorm.io/gorm"
)

type User struct {
	Id             uint   `gorm:"primary_key;auto_increment;not_null" json:"id"`
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	ProfilePicture string `json:"profileName"`
}

type UserGet struct {
	User     User  `json:"user"`
	People   int64 `json:"people"`
	Hangouts int64 `json:"hangouts"`
}

func (User) TableName() string {
	return "users"
}

// Create new User in DB
func CreateUser(db *gorm.DB, t *User) error {
	return db.Create(t).Error
}

// Get User from DB
func GetUser(db *gorm.DB, t *User) (UserGet, error) {
	err := db.Where("username = ?", t.Username).First(&t).Error
	if err != nil {
		return UserGet{}, err
	}

	var peopleCount int64
	if err := db.Table("people").Where("user = ?", t.Username).Count(&peopleCount).Error; err != nil {
		return UserGet{}, err
	}

	result := UserGet{User: *t, People: peopleCount, Hangouts: 1}

	return result, nil
}
