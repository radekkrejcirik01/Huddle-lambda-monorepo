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

type Interaction struct {
	Sender    string
	Confirmed int
}

type UserInteracted struct {
	Id           *uint  `json:"id"`
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
		Body:    t.Sender + " tapped to a Huddle 👋",
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&hangoutNotification)
}

// Get Huddle interactions from huddles_interacted table
func GetHuddleInteractions(db *gorm.DB, huddleId uint) ([]UserInteracted, *UserInteracted, error) {
	var usersInteracted []UserInteracted

	var interactions []Interaction
	if err := db.
		Table("huddles_interacted").
		Select("sender, confirmed").
		Where("huddle_id = ?", huddleId).
		Find(&interactions).
		Error; err != nil {
		return usersInteracted, nil, err
	}

	usernames := getUsernamesFromHuddlesInteracted(interactions)

	if err := db.
		Table("users").
		Where("username IN ?", usernames).
		Find(&usersInteracted).
		Error; err != nil {
		return usersInteracted, nil, err
	}

	confirmedUser := getConfirmedUser(interactions, usersInteracted)

	if confirmedUser != nil {
		usersInteracted = removeConfirmedUser(usersInteracted, confirmedUser.Username)
	}

	return usersInteracted, confirmedUser, nil
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

	if err := db.
		Table("huddles").
		Where("id = ?", t.HuddleId).
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
	var interaction HuddleInteracted

	if err := db.
		Table("huddles_interacted").
		Where("sender = ? AND huddle_id = ?", username, huddleId).
		First(&interaction).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("huddles_interacted").
		Where("sender = ? AND huddle_id = ?", username, huddleId).
		Delete(&HuddleInteracted{}).
		Error; err != nil {
		return err
	}

	if interaction.Confirmed == 1 {
		update := map[string]interface{}{
			"confirmed": 0,
			"canceled":  1,
		}

		if err := db.Table("huddles").Where("id = ?", huddleId).Updates(update).Error; err != nil {
			return err
		}
	}

	return nil
}

func getUsernamesFromHuddlesInteracted(interactions []Interaction) []string {
	usernames := make([]string, 0)
	for _, interaction := range interactions {
		usernames = append(usernames, interaction.Sender)
	}

	return usernames
}

func getConfirmedUser(interactions []Interaction, usersInteracted []UserInteracted) *UserInteracted {
	for _, interaction := range interactions {
		if interaction.Confirmed == 1 {
			for _, user := range usersInteracted {
				if user.Username == interaction.Sender {
					return &UserInteracted{
						Username:     user.Username,
						Firstname:    user.Firstname,
						ProfilePhoto: user.ProfilePhoto,
					}
				}
			}
		}
	}

	return nil
}

func removeConfirmedUser(usersInteracted []UserInteracted, confirmedUser string) []UserInteracted {
	var array []UserInteracted
	for _, userInteracted := range usersInteracted {
		if userInteracted.Username != confirmedUser {
			array = append(array, userInteracted)
		}
	}

	return array
}
