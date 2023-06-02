package people

import (
	"fmt"
	"gorm.io/gorm"
)

type Hide struct {
	Id     uint `gorm:"primary_key;auto_increment;not_null"`
	User   string
	Hidden string
}

func (Hide) TableName() string {
	return "hides"
}

type HiddenPeopleData struct {
	Id     int    `json:"id"`
	User   Person `json:"user"`
	Hidden bool   `json:"hidden"`
}

type HidePeople struct {
	User      string
	Usernames []string
}

// GetHiddenPeople from hides table
func GetHiddenPeople(db *gorm.DB, username string, lastId string) ([]HiddenPeopleData, error) {
	var hiddenPeople []HiddenPeopleData
	var invites []Invite
	var invitesUsernames []string
	var profiles []Person
	var people []Person
	var hiddenUsernames []string

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
		return nil, err
	}

	invitesUsernames = GetUsernamesFromInvites(invites, username)

	if err := db.
		Table("users").
		Where("username IN ?", invitesUsernames).
		Find(&profiles).
		Error; err != nil {
		return nil, err
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
		Table("hides").
		Select("hidden").
		Where("user = ? AND hidden IN ?",
			username, invitesUsernames).
		Find(&hiddenUsernames).
		Error; err != nil {
		return nil, err
	}

	for _, person := range people {
		hidden := IsPersonHidden(hiddenUsernames, person.Username)

		hiddenPeople = append(hiddenPeople, HiddenPeopleData{
			Id:     person.Id,
			User:   person,
			Hidden: hidden,
		})
	}

	return hiddenPeople, nil
}

// UpdateHiddenPeople in hide table
func UpdateHiddenPeople(db *gorm.DB, t *HidePeople) error {
	var hidden []string

	if err := db.
		Table("hides").
		Select("hidden").
		Where("user = ?", t.User).
		Find(&hidden).
		Error; err != nil {
		return err
	}

	hiddenUsernames := getHiddenUsernames(t.Usernames, hidden)
	hideUsernames := getHideUsernames(t.Usernames, hidden)

	if err := db.
		Table("hides").
		Where("user = ? AND hidden IN ?", t.User, hiddenUsernames).
		Delete(&Hide{}).
		Error; err != nil {
		return err
	}

	var hides []Hide
	for _, hide := range hideUsernames {
		hides = append(hides, Hide{
			User:   t.User,
			Hidden: hide,
		})
	}

	if err := db.
		Table("hides").
		Create(&hides).
		Error; err != nil {
		return err
	}

	return nil
}

func IsPersonHidden(hiddenUsernames []string, username string) bool {
	for _, hiddenUsername := range hiddenUsernames {
		if hiddenUsername == username {
			return true
		}
	}
	return false
}

func IsPersonMuted(mutedUsernames []string, username string) bool {
	for _, mutedUsername := range mutedUsernames {
		if mutedUsername == username {
			return true
		}
	}
	return false
}

func getHiddenUsernames(users []string, hiddenUsers []string) []string {
	var hiddenUsernames []string
	for _, user := range users {
		for _, hiddenUser := range hiddenUsers {
			if user == hiddenUser {
				hiddenUsernames = append(hiddenUsernames, user)
				break
			}
		}
	}
	return hiddenUsernames
}

func getHideUsernames(users []string, hiddenUsers []string) []string {
	var usernames []string
	for _, user := range users {
		isInHidden := false
		for _, hiddenUser := range hiddenUsers {
			if user == hiddenUser {
				isInHidden = true
				break
			}
		}
		if !isInHidden {
			usernames = append(usernames, user)
		}
	}
	return usernames
}
