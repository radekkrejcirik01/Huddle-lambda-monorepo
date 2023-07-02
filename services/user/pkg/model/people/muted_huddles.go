package people

import (
	"gorm.io/gorm"
)

type MutedHuddle struct {
	Id    uint `gorm:"primary_key;auto_increment;not_null"`
	User  string
	Muted string
}

func (MutedHuddle) TableName() string {
	return "muted_huddles"
}

// MuteHuddles in muted_huddles table
func MuteHuddles(db *gorm.DB, t *MutedHuddle) error {
	return db.
		Table("muted_huddles").
		Create(&t).
		Error
}

// GetMutedHuddles from muted_huddles table
func GetMutedHuddles(db *gorm.DB, username string) ([]Person, error) {
	var mutedUsernames []string
	var mutedPeople []Person

	if err := db.
		Table("muted_huddles").
		Select("muted").
		Where("user = ?", username).
		Find(&mutedUsernames).
		Error; err != nil {
		return nil, err
	}

	if err := db.
		Table("users").
		Where("username IN ?", mutedUsernames).
		Find(&mutedPeople).
		Error; err != nil {
		return nil, err
	}

	return mutedPeople, nil
}

// RemoveMutedHuddle from muted_huddles table
func RemoveMutedHuddles(db *gorm.DB, username string, user string) error {
	return db.
		Table("muted_huddles").
		Where("user = ? AND muted = ?", username, user).
		Delete(&MutedHuddle{}).
		Error
}
