package huddles

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

type HuddleInteracted struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Sender   string
	HuddleId int
	Created  int64 `gorm:"autoCreateTime"`
}

func (HuddleInteracted) TableName() string {
	return "huddles_interacted"
}

type Interact struct {
	HuddleId int
	Message  string
	Receiver string
}

type UserInteracted struct {
	Username     string `json:"username"`
	Firstname    string `json:"name"`
	ProfilePhoto string `json:"profilePhoto"`
}

// HuddleInteract to huddles_interacted table
func HuddleInteract(db *gorm.DB, username string, t *Interact) error {
	interaction := HuddleInteracted{
		Sender:   username,
		HuddleId: t.HuddleId,
	}
	if err := db.Table("huddles_interacted").Create(&interaction).Error; err != nil {
		return err
	}

	if t.Receiver == username {
		return nil
	}

	var interactionsNotifications int
	if err := db.
		Table("users").
		Select("interactions_notifications").
		Where("username = ?", t.Receiver).
		Find(&interactionsNotifications).
		Error; err != nil {
		return err
	}

	if interactionsNotifications != 1 {
		return nil
	}

	var name string
	if err := db.
		Table("users").
		Select("firstname").
		Where("username = ?", username).
		Find(&name).
		Error; err != nil {
		return err
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return nil
	}

	fcmNotification := service.FcmNotification{
		Data: map[string]interface{}{
			"type":     huddleType,
			"huddleId": t.HuddleId,
		},
		Title:   name + " liked your huddle",
		Body:    t.Message,
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// GetHuddleInteractions from huddles_interacted table
func GetHuddleInteractions(db *gorm.DB, huddleId string) ([]UserInteracted, error) {
	var usersInteracted []UserInteracted

	var interactions []string
	if err := db.
		Table("huddles_interacted").
		Select("sender").
		Where("huddle_id = ?", huddleId).
		Find(&interactions).
		Error; err != nil {
		return usersInteracted, err
	}

	if err := db.
		Table("users").
		Where("username IN ?", interactions).
		Find(&usersInteracted).
		Error; err != nil {
		return usersInteracted, err
	}

	return usersInteracted, nil
}

// RemoveHuddleInteraction from huddles_interacted table
func RemoveHuddleInteraction(db *gorm.DB, username string, huddleId string) error {
	return db.
		Table("huddles_interacted").
		Where("sender = ? AND huddle_id = ?", username, huddleId).
		Delete(&HuddleInteracted{}).
		Error
}
