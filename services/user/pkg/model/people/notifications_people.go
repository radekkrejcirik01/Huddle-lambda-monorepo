package people

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

const personInviteType = "person_invite"
const personInviteAcceptedType = "person_invite_accepted"

type PeopleNotification struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Sender   string
	Receiver string
	Type     string `gorm:"type:enum('person_invite', 'person_invite_accepted')"`
	Accepted *int
	Seen     int   `gorm:"default:0"`
	Created  int64 `gorm:"autoCreateTime"`
}

func (PeopleNotification) TableName() string {
	return "notifications_people"
}

type Person struct {
	Username     string `json:"username"`
	Firstname    string `json:"firstname"`
	ProfilePhoto string `json:"profilePhoto"`
}

type Invite struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Type     string `json:"type"`
	Accepted int    `json:"accepted"`
}

// Add person invite to notifications_people table
func AddPersonInvite(db *gorm.DB, t *PeopleNotification) (string, error) {
	errMessage := "Sorry we couldn't send a invite"

	var userExists bool
	if err := db.Table("users").Select("count(*) > 0").Where("username = ?", t.Receiver).Find(&userExists).Error; err != nil {
		return errMessage, err
	}

	if userExists {
		var user string
		if err := db.Table("notifications_people").Select("sender").Where("((sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)) AND type = 'person_invite'", t.Sender, t.Receiver, t.Receiver, t.Sender).Find(&user).Error; err != nil {
			return errMessage, err
		}

		if user == t.Sender {
			return "This user is already invited", nil
		}
		if user == t.Receiver {
			return "This user already invited you", nil
		}

		notification := PeopleNotification{
			Sender:   t.Sender,
			Receiver: t.Receiver,
			Type:     personInviteType,
		}
		if err := db.Table("notifications_people").Create(&notification).Error; err != nil {
			return errMessage, err
		}

		tokens := &[]string{}
		if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
			return "", nil
		}
		personInviteNotification := service.FcmNotification{
			Sender:  t.Sender,
			Type:    "people",
			Body:    t.Sender + " sends a friend invite",
			Sound:   "default",
			Devices: *tokens,
		}
		service.SendNotification(&personInviteNotification)

		return "Invite sent ✅", nil
	}
	return "Sorry this user doesn't exist", nil
}

// Get people list from notifications_people table
func GetPeople(db *gorm.DB, username string) ([]Person, error) {
	var acceptedInvites []PeopleNotification
	var people []Person

	if err := db.
		Table("notifications_people").
		Where("(sender = ? OR receiver = ?) AND type = 'person_invite' AND accepted = 1", username, username).
		Find(&acceptedInvites).Error; err != nil {
		return people, err
	}

	usernames := GetUsernamesFromAcceptedInvites(acceptedInvites, username)
	if err := db.Table("users").Where("username IN ?", usernames).Find(&people).Error; err != nil {
		return people, err
	}

	return people, nil
}

// Update accepted column in notifications_people table and
// create notification with accepted invite type
func AcceptPersonInvite(db *gorm.DB, t *PeopleNotification) error {
	if err := db.
		Table("notifications_people").
		Where("sender = ? AND receiver = ? AND type = 'person_invite'", t.Receiver, t.Sender).
		Update("accepted", 1).Error; err != nil {
		return err
	}

	notification := PeopleNotification{
		Sender:   t.Sender,
		Receiver: t.Receiver,
		Type:     personInviteAcceptedType,
	}
	if err := db.Table("notifications_people").Create(&notification).Error; err != nil {
		return err
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return nil
	}
	acceptFriendInviteNotification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    "people",
		Body:    t.Sender + " accepted friend invite 🙌",
		Sound:   "default",
		Devices: *tokens,
	}
	service.SendNotification(&acceptFriendInviteNotification)

	return nil
}

// Get person invite information from notifications_people table
func GetPersonInvite(db *gorm.DB, user1 string, user2 string) (Invite, error) {
	var invite Invite
	err := db.
		Table("notifications_people").
		Where("((sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)) AND type = 'person_invite'",
			user1, user2, user2, user1).
		First(&invite).Error

	return invite, err
}

// Remove person connection from notifications_people table
func RemovePerson(db *gorm.DB, user1 string, user2 string) error {
	return db.
		Table("notifications_people").
		Where("((sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)) AND type = 'person_invite'",
			user1, user2, user2, user1).
		Delete(&PeopleNotification{}).
		Error
}

// Get usernames from accepted invites array
func GetUsernamesFromAcceptedInvites(acceptedInvites []PeopleNotification, username string) []string {
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
