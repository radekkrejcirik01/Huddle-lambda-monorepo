package people

import (
	"errors"
	"fmt"

	n "github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/notifications"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

type Invite struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Sender   string
	Receiver string
	Accepted int
	Created  int64 `gorm:"autoCreateTime"`
}

func (Invite) TableName() string {
	return "invites"
}

type Person struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Firstname    string `json:"name"`
	ProfilePhoto string `json:"profilePhoto"`
}

// Add invite to invites table
func AddPersonInvite(db *gorm.DB, t *Invite) (string, error) {
	if err := db.Table("users").Where("username = ?", t.Receiver).First(&Person{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "User doesn't exist yet", nil
		}

		return "", err
	}

	invite := Invite{
		Sender:   t.Sender,
		Receiver: t.Receiver,
	}
	if err := db.
		Table("invites").
		Where("(sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)",
			t.Sender, t.Receiver, t.Receiver, t.Sender).
		FirstOrCreate(&invite).
		Error; err != nil {
		return "", err
	}

	notification := n.Notification{
		Sender:   t.Sender,
		Receiver: t.Receiver,
		EventId:  int(invite.Id),
		Type:     n.PersonInviteType,
	}
	if err := db.Table("notifications").Create(&notification).Error; err != nil {
		return "Sorry we couldn't send an invite", err
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return "", nil
	}

	fcmNotification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    "people",
		Body:    t.Sender + " sends a friend invite",
		Sound:   "default",
		Devices: *tokens,
	}

	service.SendNotification(&fcmNotification)

	return "Invite sent ✅", nil
}

// Get people from invites table
func GetPeople(db *gorm.DB, username string, lastId string) ([]Person, error) {
	var invites []Invite
	var people []Person

	var idCondition string
	if lastId != "" {
		idCondition = fmt.Sprintf("id < %s AND ", lastId)
	}
	if err := db.
		Table("invites").
		Where(idCondition+"(sender = ? OR receiver = ?) AND accepted = 1",
			username, username).
		Order("id DESC").
		Limit(20).
		Find(&invites).Error; err != nil {
		return []Person{}, err
	}

	usernames := GetUsernamesFromInvites(invites, username)

	var profiles []Person
	if err := db.Table("users").Where("username IN ?", usernames).Find(&profiles).Error; err != nil {
		return []Person{}, err
	}

	for _, invite := range invites {
		var inviteUser string

		if invite.Sender == username {
			inviteUser = invite.Receiver
		} else {
			inviteUser = invite.Sender
		}

		person := getPersonByUsername(profiles, inviteUser)
		person.Id = int(invite.Id)

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

	notification := n.Notification{
		Sender:   t.Sender,
		Receiver: t.Receiver,
		EventId:  int(t.Id),
		Type:     n.PersonInviteAcceptType,
	}
	if err := db.Table("notifications").Create(&notification).Error; err != nil {
		return err
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return nil
	}

	fcmNotification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    "people",
		Body:    t.Sender + " accepted friend invite 🙌",
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// Get person invite information from notifications table
func GetPersonInvite(db *gorm.DB, user1 string, user2 string) (Invite, error) {
	var invite Invite
	err := db.
		Table("invites").
		Where("(sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)",
			user1, user2, user2, user1).
		First(&invite).Error

	return invite, err
}

// Update accepted column in invites table to 0
func RemovePerson(db *gorm.DB, user1 string, user2 string) error {
	return db.
		Table("invites").
		Where("(sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)",
			user1, user2, user2, user1).
		Update("accepted", 0).
		Error
}

// Get usernames from accepted invites
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
