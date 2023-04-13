package notify

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/notifications"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

const hangoutType = "hangout_notify"

type HangoutInvitation struct {
	Sender     string
	SenderName string
	Receiver   string
}

// Send a hangout notification in DB
func Notify(db *gorm.DB, t *HangoutInvitation) error {
	notification := notifications.Notifications{
		Sender:   t.Sender,
		Receiver: t.Receiver,
		Type:     hangoutType,
	}

	if err := db.Table("notifications").Create(&notification).Error; err != nil {
		return err
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return nil
	}
	hangoutNotification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    hangoutType,
		Title:   t.Sender + " sends a hangout!",
		Sound:   "notification.wav",
		Devices: *tokens,
	}
	service.SendNotification(&hangoutNotification)
	return nil
}
