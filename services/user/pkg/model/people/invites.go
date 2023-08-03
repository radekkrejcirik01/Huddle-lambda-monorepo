package people

import (
	"errors"
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
		return "Cannot invite yourself 😀", nil
	}

	if err := db.Table("users").Where("username = ?", t.Receiver).First(&Person{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "User doesn't exist yet", nil
		}

		return "", err
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

	invite := Invite{
		Sender:   t.Sender,
		Receiver: t.Receiver,
	}
	if rows := db.
		Table("invites").
		Where("(sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)",
			t.Sender, t.Receiver, t.Receiver, t.Sender).
		FirstOrCreate(&invite).
		RowsAffected; rows == 0 {
		var message string
		if invite.Sender == t.Sender {
			message = "Invite already sent"
		} else {
			message = t.Receiver + " already invited you"
		}
		return message, nil
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

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return "", nil
	}

	fcmNotification := service.FcmNotification{
		Data: map[string]interface{}{
			"type": "invite",
		},
		Body:    t.Sender + " sends friend invite",
		Sound:   "default",
		Devices: *tokens,
	}

	err := service.SendNotification(&fcmNotification)
	if err != nil {
		fmt.Println(err)
	}

	return "Invite sent ✅", nil
}

// Get people from invites table
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

		person.Id = getInviteId(invites, username)
		person.User = getPersonByUsername(profiles, username)
		person.Accepted = isInviteAccepted(invites, username)

		people = append(people, person)
	}

	return people, nil
}

// Update accepted column in invites table to 1
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

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return nil
	}

	fcmNotification := service.FcmNotification{
		Data: map[string]interface{}{
			"type": "friends",
		},
		Body:    t.Sender + " accepted friend invite 🙌",
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&fcmNotification)
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
func getPersonByUsername(profiles []Person, username string) Person {
	var person Person
	for _, profile := range profiles {
		if profile.Username == username {
			person = profile
			break
		}
	}
	return person
}

func getInviteId(invites []Invite, username string) int {
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
