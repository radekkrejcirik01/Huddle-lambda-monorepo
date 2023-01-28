package hangouts

import (
	"time"

	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/notifications"
	"gorm.io/gorm"
)

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
	Id       uint
	Value    int
	User     string
	Username string
}

func (HangoutsInvitationTable) TableName() string {
	return "hangouts_invitations"
}

// Accept hangout invitation from DB
func AcceptHangout(db *gorm.DB, t *AcceptInvite) error {
	if err := db.Table("hangouts_invitations").Where("hangout_id = ?", t.Id).Update("confirmed", t.Value).Error; err != nil {
		return err
	}

	now := time.Now().Format(timeFormat)
	acceptedInvitation := notifications.AcceptedInvitations{
		EventId:  t.Id,
		User:     t.User,
		Username: t.Username,
		Time:     now,
		Type:     "hangout",
	}

	return db.Table("accepted_invitations").Create(acceptedInvitation).Error
}
