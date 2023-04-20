package notify

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

const hangoutType = "hangout_notify"

type NotifyNotification struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Sender   string
	Receiver string
	Type     string `gorm:"default:'hangout_notify'"`
	Seen     int    `gorm:"default:0"`
	Created  int64  `gorm:"autoCreateTime"`
}

func (NotifyNotification) TableName() string {
	return "notifications_notify"
}

type Notify struct {
	Sender     string
	SenderName string
	Receiver   string
}

// Send a hangout notification
func SendNotify(db *gorm.DB, t *Notify) error {
	notification := NotifyNotification{
		Sender:   t.Sender,
		Receiver: t.Receiver,
		Type:     hangoutType,
	}

	if err := db.Table("notifications_notify").Create(&notification).Error; err != nil {
		return err
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return nil
	}
	hangoutNotification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    hangoutType,
		Title:   t.SenderName + ": Let's hangout!",
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&hangoutNotification)
}
