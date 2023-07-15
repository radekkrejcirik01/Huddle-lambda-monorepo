package people

import (
	"gorm.io/gorm"
)

type Blocked struct {
	Id      uint `gorm:"primary_key;auto_increment;not_null"`
	User    string
	Blocked string
}

func (Blocked) TableName() string {
	return "blocked"
}

// BlockUser in blocked table
func BlockUser(db *gorm.DB, t *Blocked) error {
	if err := db.
		Table("blocked").
		Where("user = ? AND blocked = ?", t.User, t.Blocked).
		FirstOrCreate(&t).
		Error; err != nil {
		return err
	}

	return db.
		Table("invites").
		Where("(sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)",
			t.User, t.Blocked, t.Blocked, t.User).
		Delete(&Invite{}).
		Error
}
