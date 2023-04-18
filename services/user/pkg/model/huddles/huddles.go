package huddles

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	p "github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"gorm.io/gorm"
)

// Huddle is a communication app for creating hang outs with people by sharing simple posts called Huddles
type Huddle struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	CreatedBy string
	What      string
	Where     string
	When      string
	Type      string `gorm:"type:enum('huddle', 'group_huddle')"`
	Created   int64  `gorm:"autoCreateTime"`
}

func (Huddle) TableName() string {
	return "huddles"
}

type HuddlesData struct {
	Id           uint   `json:"id"`
	CreatedBy    string `json:"createdBy"`
	Name         string `json:"name"`
	ProfilePhoto string `json:"profilePhoto"`
	What         string `json:"what"`
	Where        string `json:"where"`
	When         string `json:"when"`
	Type         string `json:"type"`
	Interacted   int    `json:"interacted"`
}

type Invite struct {
	Sender   string
	Receiver string
}

// Get Huddles from huddles table
func GetHuddles(db *gorm.DB, username string) ([]HuddlesData, error) {
	var huddlesData []HuddlesData

	var invites []Invite
	if err := db.
		Table("notifications_people").
		Where("(sender = ? OR receiver = ?) AND type = 'person_invite' AND accepted = 1", username, username).
		Find(&invites).Error; err != nil {
		return huddlesData, err
	}

	people := GetUsernamesFromInvites(invites, username)
	// Get also user's Huddles
	people = append(people, username)

	var huddles []Huddle
	if err := db.
		Table("huddles").
		Where("created_by IN ?", people).
		Find(&huddles).Error; err != nil {
		return huddlesData, err
	}

	if len(huddles) < 1 {
		return huddlesData, nil
	}

	huddlesIds := GetIdsFromHuddlesArray(huddles)
	var interactedHuddlesIds []uint
	if err := db.
		Table("huddles_interacted").
		Select("huddle_id").Where("username = ? AND huddle_id IN ?", username, huddlesIds).
		Find(&interactedHuddlesIds).Error; err != nil {
		return huddlesData, err
	}

	users := GetUsernamesFromHuddles(huddles)
	var profiles []p.Person
	if err := db.Table("users").Where("username IN ?", users).Find(&profiles).Error; err != nil {
		return huddlesData, err
	}

	for _, huddle := range huddles {
		profileInfo := GetProfileInfoFromProfiles(profiles, huddle.CreatedBy)
		interacted := GetInteraction(interactedHuddlesIds, huddle.Id)

		huddlesData = append(huddlesData, HuddlesData{
			Id:           huddle.Id,
			CreatedBy:    huddle.CreatedBy,
			Name:         profileInfo.Firstname,
			ProfilePhoto: profileInfo.ProfilePicture,
			What:         huddle.What,
			Where:        huddle.Where,
			When:         huddle.When,
			Type:         huddle.Type,
			Interacted:   interacted,
		})

	}

	return huddlesData, nil
}

// Get usernames from invites array
func GetUsernamesFromInvites(invites []Invite, username string) []string {
	usernames := make([]string, 0)
	for _, invite := range invites {
		if invite.Sender == username {
			usernames = append(usernames, invite.Receiver)
		} else {
			usernames = append(usernames, invite.Sender)
		}
	}

	return usernames
}

// Get ids from huddles array
func GetIdsFromHuddlesArray(huddles []Huddle) []uint {
	var ids []uint
	for _, h := range huddles {
		ids = append(ids, h.Id)
	}

	return ids
}

// Get usernames from huddles array
func GetUsernamesFromHuddles(huddles []Huddle) []string {
	usernames := make([]string, 0)
	for _, h := range huddles {
		usernames = append(usernames, h.CreatedBy)
	}

	return usernames
}

// Get profile info from profiles array
func GetProfileInfoFromProfiles(profiles []people.Person, username string) p.Person {
	for _, profile := range profiles {
		if profile.Username == username {
			return profile
		}
	}

	return p.Person{}
}

// Return interaction value from interacted huddles ids and huddle id
func GetInteraction(interactedHuddlesIds []uint, huddleId uint) int {
	for _, id := range interactedHuddlesIds {
		if id == huddleId {
			return 1
		}
	}

	return 0
}
