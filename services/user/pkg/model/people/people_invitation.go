package people

import (
	"time"

	"gorm.io/gorm"
)

const timeFormat = "2006-01-02 15:04:05"

type PeopleInvitationTable struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	User     string
	Username string
	Time     string
}

func (PeopleInvitationTable) TableName() string {
	return "people_invitation"
}

// Create new user invitatiom in DB
func CreatePeopleInvitation(db *gorm.DB, t *PeopleInvitationTable) (string, error) {
	var exists bool
	err := db.Table("users").Select("count(*) > 0").Where("username = ?", t.Username).Find(&exists).Error

	if exists {
		time := time.Now()
		t.Time = time.Format(timeFormat)
		return "User succesfully invited! 🎉", db.Create(t).Error
	}
	return "Sorry, this user does not exists 😔", err
}
