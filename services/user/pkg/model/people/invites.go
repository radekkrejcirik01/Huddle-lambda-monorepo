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
	Created  int64 `gorm:"autoCreateTime"`
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

type InviteResponseData struct {
	Id   int    `json:"id"`
	User Person `json:"user"`
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
		Sender:  t.Sender,
		Type:    "people",
		Body:    t.Sender + " sends a friend invite",
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
func GetPeople(db *gorm.DB, username string, lastId string) ([]Person, int64, error) {
	var invites []Invite
	var people []Person
	var invitesNumber int64

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
		return nil, 0, err
	}

	usernames := GetUsernamesFromInvites(invites, username)

	var profiles []Person
	if err := db.Table("users").Where("username IN ?", usernames).Find(&profiles).Error; err != nil {
		return nil, 0, err
	}

	// Reorder people by invites
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

	if err := db.
		Table("invites").
		Where("receiver = ? AND accepted = 0", username).
		Count(&invitesNumber).
		Error; err != nil {
		return nil, 0, err
	}

	return people, invitesNumber, nil
}

// Get invites from invites table
func GetInvites(db *gorm.DB, username string) ([]InviteResponseData, error) {
	var invites []Invite
	var profiles []Person
	var invitesResponse []InviteResponseData

	if err := db.
		Table("invites").
		Where("receiver = ? AND accepted = 0", username).
		Order("id DESC").
		Limit(20).
		Find(&invites).
		Error; err != nil {
		return nil, err
	}

	senders := getSenders(invites)

	if err := db.
		Table("users").
		Select("username, profile_photo").
		Where("username IN ?", senders).
		Find(&profiles).
		Error; err != nil {
		return nil, err
	}

	for _, invite := range invites {
		user := getUser(profiles, invite.Sender)
		invitesResponse = append(invitesResponse, InviteResponseData{
			Id:   int(invite.Id),
			User: user,
		})
	}

	return invitesResponse, nil
}

// Update accepted column in invites table to 1
func AcceptPersonInvite(db *gorm.DB, t *Invite) error {
	if err := db.
		Table("invites").
		Where("id = ?", t.Id).
		Update("accepted", 1).Error; err != nil {
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

// Get senders from invites
func getSenders(invites []Invite) []string {
	var senders []string
	for _, invite := range invites {
		senders = append(senders, invite.Sender)
	}
	return senders
}

// Get senders from invites
func getUser(profiles []Person, username string) Person {
	var profile Person
	for _, p := range profiles {
		if p.Username == username {
			profile = p
			break
		}
	}
	return profile
}
