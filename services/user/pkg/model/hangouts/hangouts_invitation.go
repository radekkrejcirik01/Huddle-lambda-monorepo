package hangouts

import "gorm.io/gorm"

type HangoutsInvitationTable struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	HangoutId uint
	Username  string
	Confirmed uint
}

func (HangoutsInvitationTable) TableName() string {
	return "hangouts_invitations"
}

// Create new hangout invitation in DB
func CreateHangoutInvitation(db *gorm.DB, t *HangoutsInvitationTable) error {
	return db.Create(t).Error
}
