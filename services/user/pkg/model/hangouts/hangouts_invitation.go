package hangouts

import "gorm.io/gorm"

type HangoutsInvitationTable struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	HangoutId uint
	User      string
	Username  string
	Time      string
	Confirmed int `gorm:"default:0"`
	Seen      int `gorm:"default:0"`
}

type AcceptInvite struct {
	Id    uint
	Value int
}

func (HangoutsInvitationTable) TableName() string {
	return "hangouts_invitations"
}

// Accept hangout invitation from DB
func AcceptHangout(db *gorm.DB, t *AcceptInvite) error {
	return db.Table("hangouts_invitations").Where("hangout_id = ?", t.Id).Update("confirmed", t.Value).Error
}
