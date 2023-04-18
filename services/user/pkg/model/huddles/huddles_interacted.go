package huddles

import "gorm.io/gorm"

type HuddleInteracted struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Username string
	HuddleId uint
	Created  int64 `gorm:"autoCreateTime"`
}

func (HuddleInteracted) TableName() string {
	return "huddles_interacted"
}

// Add Huddle interaction to huddles_interacted table
func HuddleInteract(db *gorm.DB, t *HuddleInteracted) error {
	return db.
		Table("huddles_interacted").
		Where("username = ? AND huddle_id = ?", t.Username, t.HuddleId).
		FirstOrCreate(&t).
		Error
}
