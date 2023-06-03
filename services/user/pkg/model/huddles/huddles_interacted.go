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
	Sender   string
	Receiver string
}

type Interaction struct {
	Sender    string
	Confirmed int
}

type UserInteracted struct {
	Username     string `json:"username"`
	Firstname    string `json:"name"`
	ProfilePhoto string `json:"profilePhoto"`
}

type RemoveConfirm struct {
	Id int
}

// Add Huddle interaction to huddles_interacted table
func HuddleInteract(db *gorm.DB, t *Interact) error {
	interaction := HuddleInteracted{
		Sender:   t.Sender,
		HuddleId: t.HuddleId,
	}
	if err := db.Table("huddles_interacted").Create(&interaction).Error; err != nil {
		return err
	}

	var interactionNotification int
	if err := db.
		Table("users").
		Select("interactions_notifications").
		Where("username = ?", t.Receiver).
		Find(&interactionNotification).
		Error; err != nil {
		return err
	}

	if interactionNotification != 1 {
		return nil
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return nil
	}

	fcmNotification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    huddleType,
		Body:    t.Sender + " interacted with your Huddle 👋",
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// Get Huddle interactions from huddles_interacted table
func GetHuddleInteractions(db *gorm.DB, huddleId int) ([]UserInteracted, error) {
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

// Remove Huddle interaction from huddles_interacted table
func RemoveHuddleInteraction(db *gorm.DB, username string, huddleId uint) error {
	return db.
		Table("huddles_interacted").
		Where("sender = ? AND huddle_id = ?", username, huddleId).
		Delete(&HuddleInteracted{}).
		Error
}
