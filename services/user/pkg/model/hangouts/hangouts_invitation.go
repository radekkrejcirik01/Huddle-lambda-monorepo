package hangouts

import (
	"strings"
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
	Type      string
	Confirmed int `gorm:"default:0"`
	Seen      int `gorm:"default:0"`
}

type AcceptInvite struct {
	Id       uint
	Value    int
	User     string
	Username string
	Name     string
	Type     string
}

type HangoutInvitations struct {
	HangoutId uint
	User      string
	Name      string
	Usernames []string
}

type Get struct {
	HangoutId uint
}

func (HangoutsInvitationTable) TableName() string {
	return "hangouts_invitations"
}

// Accept hangout invitation in DB
func AcceptHangout(db *gorm.DB, t *AcceptInvite) error {
	if err := db.Table("hangouts_invitations").Where("hangout_id = ? AND username = ?", t.Id, t.User).Update("confirmed", t.Value).Error; err != nil {
		return err
	}

	now := time.Now().Format(timeFormat)
	acceptedType := "accepted_people"
	if t.Type == hangoutType {
		acceptedType = "accepted_hangout"
	}
	if t.Type == groupHangoutType {
		acceptedType = "accepted_group_hangout"
	}
	acceptedInvitation := notifications.AcceptedInvitations{
		EventId:  t.Id,
		User:     t.User,
		Username: t.Username,
		Time:     now,
		Type:     acceptedType,
	}

	if rowsAffected := db.Table("accepted_invitations").Where(notifications.AcceptedInvitations{EventId: t.Id, User: t.User}).FirstOrCreate(&acceptedInvitation).RowsAffected; rowsAffected == 0 {
		return nil
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Username); err != nil {
		return nil
	}

	body := t.Name + " accepted hangout invite!"
	if t.Type == groupHangoutType {
		body = t.Name + " accepted group hangout invite!"
	}
	acceptHangoutInviteNotification := service.FcmNotification{
		Sender:  t.User,
		Type:    acceptedType,
		Body:    body,
		Sound:   "notification.wav",
		Devices: *tokens,
	}
	service.SendNotification(&acceptHangoutInviteNotification)

	return nil
}

// Get hangout usernames from DB
func GetHangoutUsernames(db *gorm.DB, t *Get) ([]string, error) {
	var usernames []string
	err := db.Table("hangouts_invitations").Where("hangout_id = ?", t.HangoutId).Select("username").Find(&usernames).Error

	return usernames, err
}

// Add hangout invitation in DB
func SendHangoutInvitation(db *gorm.DB, t *HangoutInvitations) error {
	now := time.Now().Format(timeFormat)

	var hangoutInvitations []HangoutsInvitationTable
	for _, username := range t.Usernames {
		hangoutInvitations = append(hangoutInvitations, HangoutsInvitationTable{
			HangoutId: t.HangoutId,
			User:      t.User,
			Username:  username,
			Time:      now,
			Type:      groupHangoutType,
		})
	}

	if err := db.Table("hangouts_invitations").Create(&hangoutInvitations).Error; err != nil {
		return err
	}

	if err := db.Table("hangouts").Where("id = ?", t.HangoutId).Updates(map[string]interface{}{"title": "Group hangout", "type": groupHangoutType}).Error; err != nil {
		return err
	}

	var usernamesArray []string
	for _, username := range t.Usernames {
		usernamesArray = append(usernamesArray, `'`+username+`'`)
	}

	usernamesString := strings.Join(usernamesArray, ", ")

	tokens := &[]string{}
	if err := service.GetTokensByUsernames(db, tokens, usernamesString); err != nil {
		return nil
	}

	hangoutNotification := service.FcmNotification{
		Sender:  t.User,
		Type:    groupHangoutType,
		Title:   t.Name + " sends a group hangout!",
		Sound:   "notification.wav",
		Devices: *tokens,
	}
	service.SendNotification(&hangoutNotification)
	return nil
}
