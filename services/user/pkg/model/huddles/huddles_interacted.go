package huddles

import (
	n "github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/notifications"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

type HuddleInteracted struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	Sender    string
	HuddleId  int
	Confirmed int   `gorm:"default:0"`
	Created   int64 `gorm:"autoCreateTime"`
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
	notification := n.Notification{
		Sender:   t.Sender,
		Receiver: t.Receiver,
		EventId:  t.HuddleId,
		Type:     n.HuddleInteractType,
	}

	if err := db.Table("notifications").Create(&notification).Error; err != nil {
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
func GetHuddleInteractions(db *gorm.DB, huddleId int) ([]UserInteracted, *UserInteracted, error) {
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

// Update confirmed value in huddles_interacted table
func ConfirmHuddle(db *gorm.DB, t *Interact) error {
	notification := n.Notification{
		Sender:   t.Sender,
		Receiver: t.Receiver,
		EventId:  t.HuddleId,
		Type:     n.HuddleConfirmType,
	}

	if err := db.Table("notifications").Create(&notification).Error; err != nil {
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

	fcmNotification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    huddleType,
		Body:    t.Sender + " confirmed a Huddle ✅",
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&fcmNotification)
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

// Remove Huddle confirm in huddles and huddles_interacted tables
func RemoveHuddleConfirm(db *gorm.DB, t *RemoveConfirm) error {
	if err := db.Table("huddles").Where("id = ?", t.Id).Update("confirmed", 0).Error; err != nil {
		return err
	}

	if err := db.
		Table("huddles_interacted").
		Where("huddle_id = ? AND confirmed = 1", t.Id).
		Update("confirmed", 0).
		Error; err != nil {
		return err
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
