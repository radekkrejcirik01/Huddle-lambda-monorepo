package people

import (
	"fmt"

	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

type Invite struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Sender   string
	Receiver string
	Accepted int
	Seen     int
}

func (Invite) TableName() string {
	return "invites"
}

type Person struct {
	Id           int    `json:"id,omitempty"`
	Username     string `json:"username"`
	Firstname    string `json:"name,omitempty"`
	ProfilePhoto string `json:"profilePhoto"`
}

type PeopleData struct {
	Id       int    `json:"id"`
	User     Person `json:"user"`
	Accepted int    `json:"accepted"`
}

// AddPersonInvite in invites table
func AddPersonInvite(db *gorm.DB, t *Invite) (string, error) {
	if t.Sender == t.Receiver {
		return "Why are you inviting yourself? 😀", nil
	}

	var user Person
	if err := db.
		Table("users").
		Where("username = ?", t.Receiver).
		Find(&user).
		Error; err != nil {
		return "", err
	}

	if len(user.Username) == 0 {
		return "Username doesn't exist", nil
	}

	var blockedUsernames []string
	if err := db.
		Table("blocked").
		Select("blocked").
		Where("user = ? AND blocked = ?", t.Receiver, t.Sender).
		Find(&blockedUsernames).
		Error; err != nil {
		return "", err
	}

	if len(blockedUsernames) > 0 {
		return "Could not send invite to this user", nil
	}

	var invite Invite
	if err := db.
		Table("invites").
		Where("(sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)",
			t.Sender, t.Receiver, t.Receiver, t.Sender).
		Find(&invite).
		Error; err != nil {
		return "", err
	}

	if invite.Sender == t.Sender {
		return "Invite already sent", nil
	}
	if invite.Sender == t.Receiver {
		return "User already invited you", nil
	}

	newInvite := Invite{
		Sender:   t.Sender,
		Receiver: t.Receiver,
	}
	if err := db.Table("invites").Create(&newInvite).Error; err != nil {
		return "", err
	}

	var friendInviteNotification int
	if err := db.
		Table("users").
		Select("friends_invites_notifications").
		Where("username = ?", t.Receiver).
		Find(&friendInviteNotification).
		Error; err != nil {
		return "", err
	}

	if friendInviteNotification != 1 {
		return "Invite sent ✅", nil
	}

	tokens, err := service.GetTokensByUsername(db, t.Receiver)
	if err != nil {
		return "", err
	}

	fcmNotification := service.FcmNotification{
		Data: map[string]interface{}{
			"type": "contacts",
		},
		Body:    t.Sender + " sends a friend invite",
		Sound:   "default",
		Devices: tokens,
	}

	return "Invite sent ✅", service.SendNotification(&fcmNotification)
}

// GetPeople from invites table
func GetPeople(db *gorm.DB, username string, lastId string) ([]PeopleData, error) {
	var invites []Invite
	var people []PeopleData

	var idCondition string
	if lastId != "" {
		idCondition = fmt.Sprintf("id < %s AND ", lastId)
	}
	if err := db.
		Table("invites").
		Where(idCondition+"((sender = ? OR receiver = ?) AND accepted = 1) OR (receiver = ? AND accepted = 0)",
			username, username, username).
		Order("id DESC").
		Limit(20).
		Find(&invites).Error; err != nil {
		return nil, err
	}

	usernames := GetUsernamesFromInvites(invites, username)

	var profiles []Person
	if err := db.
		Table("users").
		Where("username IN ?", usernames).
		Find(&profiles).
		Error; err != nil {
		return nil, err
	}

	// Reorder people by invites
	for _, username := range usernames {
		var person PeopleData

		person.Id = GetInviteId(invites, username)
		person.User = GetPersonByUsername(profiles, username)
		person.Accepted = isInviteAccepted(invites, username)

		people = append(people, person)
	}

	return people, nil
}

// AcceptPersonInvite in invites table
func AcceptPersonInvite(db *gorm.DB, t *Invite) error {
	if err := db.
		Table("invites").
		Where("id = ?", t.Id).
		Update("accepted", 1).Error; err != nil {
		return err
	}

	if err := db.
		Table("blocked").
		Where("user = ? AND blocked = ?", t.Receiver, t.Sender).
		Delete(&Blocked{}).
		Error; err != nil {
		return err
	}

	tokens, err := service.GetTokensByUsername(db, t.Receiver)
	if err != nil {
		return err
	}

	fcmNotification := service.FcmNotification{
		Data: map[string]interface{}{
			"type": "contacts",
		},
		Body:    t.Sender + " accepted friend invite 🙌",
		Sound:   "default",
		Devices: tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// DeleteInvite in invites table
func DeleteInvite(db *gorm.DB, id string) error {
	return db.
		Table("invites").
		Where("id = ?", id).
		Delete(&Invite{}).
		Error
}

// GetUnseenInvites from invites table
func GetUnseenInvites(db *gorm.DB, username string) (int64, error) {
	var number int64
	if err := db.
		Table("invites").
		Where("receiver = ? AND seen != 1", username).
		Count(&number).
		Error; err != nil {
		return 0, err
	}
	return number, nil
}

// UpdateSeenInvites in invites table
func UpdateSeenInvites(db *gorm.DB, username string) error {
	return db.
		Table("invites").
		Where("receiver = ? AND seen != 1", username).
		Update("seen", 1).
		Error
}

// GetUsernamesFromInvites
func GetUsernamesFromInvites(acceptedInvites []Invite, username string) []string {
	usernames := make([]string, 0)
	for _, invite := range acceptedInvites {
		if invite.Sender == username {
			usernames = append(usernames, invite.Receiver)
		} else {
			usernames = append(usernames, invite.Sender)
		}
	}

	return usernames
}

// Get person by username from profiles
func GetPersonByUsername(profiles []Person, username string) Person {
	var person Person
	for _, profile := range profiles {
		if profile.Username == username {
			person = profile
			break
		}
	}
	return person
}

func GetInviteId(invites []Invite, username string) int {
	for _, invite := range invites {
		if invite.Sender == username || invite.Receiver == username {
			return int(invite.Id)
		}
	}

	return 0
}

func isInviteAccepted(invites []Invite, username string) int {
	for _, invite := range invites {
		if invite.Sender == username || invite.Receiver == username {
			return invite.Accepted
		}
	}

	return 0
}
