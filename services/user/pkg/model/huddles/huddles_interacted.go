package huddles

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

const huddleInteractedType = "huddle_interacted"
const huddleConfirmedType = "huddle_confirmed"

type HuddleInteracted struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	Sender    string
	HuddleId  uint
	Confirmed int   `gorm:"default:0"`
	Created   int64 `gorm:"autoCreateTime"`
}

func (HuddleInteracted) TableName() string {
	return "huddles_interacted"
}

type SenderConfirmed struct {
	Sender    string
	Confirmed int
}

type HuddleInteractedData struct {
	Id           uint   `json:"id"`
	Username     string `json:"username"`
	Firstname    string `json:"name"`
	ProfilePhoto string `json:"profilePhoto"`
}

// Add Huddle interaction to huddles_interacted table
func HuddleInteract(db *gorm.DB, t *HuddleNotification) error {
	t.Type = huddleInteractedType

	if err := db.Table("notifications_huddles").Create(&t).Error; err != nil {
		return err
	}

	interaction := HuddleInteracted{
		Sender:   t.Sender,
		HuddleId: t.HuddleId,
	}
	if err := db.Table("huddles_interacted").Create(&interaction).Error; err != nil {
		return err
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return nil
	}
	hangoutNotification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    huddleType,
		Body:    t.Sender + " interacted with your Huddle 👋",
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&hangoutNotification)
}

// Get Huddle interactions from huddles_interacted table
func GetHuddleInteractions(db *gorm.DB, huddleId uint) ([]HuddleInteractedData, *string, error) {
	var huddleInteractedData []HuddleInteractedData

	var huddlesInteracted []SenderConfirmed
	if err := db.
		Table("huddles_interacted").
		Select("sender, confirmed").
		Where("huddle_id = ?", huddleId).
		Find(&huddlesInteracted).
		Error; err != nil {
		return huddleInteractedData, nil, err
	}

	usernames := getUsernamesFromHuddlesInteracted(huddlesInteracted)

	if err := db.
		Table("users").
		Where("username IN ?", usernames).
		Find(&huddleInteractedData).
		Error; err != nil {
		return huddleInteractedData, nil, err
	}

	confirmedUser := getConfirmedUser(huddlesInteracted)

	return huddleInteractedData, confirmedUser, nil
}

// Confirm Huddle interaction, add notification to notifications_huddles table
// and update confirmed value in huddles_interacted table
func ConfirmHuddle(db *gorm.DB, t *HuddleNotification) error {
	t.Type = huddleConfirmedType

	if err := db.Table("notifications_huddles").Create(&t).Error; err != nil {
		return err
	}

	if err := db.
		Table("huddles_interacted").
		Where("sender = ? AND huddle_id = ?", t.Receiver, t.HuddleId).
		Update("confirmed", 1).
		Error; err != nil {
		return err
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return nil
	}
	hangoutNotification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    huddleType,
		Body:    t.Sender + " confirmed a Huddle ✅",
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&hangoutNotification)
}

// Remove Huddle interaction from huddles_interacted table
func RemoveHuddleInteraction(db *gorm.DB, username string, huddleId uint) error {
	return db.
		Table("huddles_interacted").
		Where("sender = ? AND huddle_id = ?", username, huddleId).
		Delete(&HuddleInteracted{}).
		Error
}

func getUsernamesFromHuddlesInteracted(huddlesInteracted []SenderConfirmed) []string {
	usernames := make([]string, 0)
	for _, huddleInteracted := range huddlesInteracted {
		usernames = append(usernames, huddleInteracted.Sender)
	}

	return usernames
}

func getConfirmedUser(huddlesInteracted []SenderConfirmed) *string {
	for _, huddleInteracted := range huddlesInteracted {
		if huddleInteracted.Confirmed == 1 {
			return &huddleInteracted.Sender
		}
	}

	return nil
}
