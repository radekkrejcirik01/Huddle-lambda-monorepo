package hangouts

import (
	"time"

	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/notifications"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
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
	Name     string
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
		Type:     "accepted_hangout",
	}

	if rowsAffected := db.Table("accepted_invitations").Create(&acceptedInvitation).RowsAffected; rowsAffected == 0 {
		return nil
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Username); err != nil {
		return nil
	}

	acceptHangoutInviteNotification := service.FcmNotification{
		Sender:  t.User,
		Type:    "hangout",
		Body:    t.Name + " accepted hangout invite!",
		Sound:   "notification.wav",
		Devices: *tokens,
	}
	service.SendNotification(&acceptHangoutInviteNotification)

	return nil
}
