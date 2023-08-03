package huddles

import (
	"fmt"

	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	p "github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

// Huddle is a communication app for suggesting hangouts by adding simple posts called Huddles
type Huddle struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	CreatedBy string
	Message   string
	Created   int64 `gorm:"autoCreateTime"`
}

func (Huddle) TableName() string {
	return "huddles"
}

type NewHuddle struct {
	Name    string
	Message string
}

type HuddleData struct {
	Id             int    `json:"id"`
	CreatedBy      string `json:"createdBy"`
	Name           string `json:"name"`
	ProfilePhoto   string `json:"profilePhoto"`
	Message        string `json:"message"`
	Liked          int    `json:"liked,omitempty"`
	CommentsNumber int    `json:"commentsNumber"`
	LikesNumber    int    `json:"likesNumber"`
}

type Invite struct {
	Sender   string
	Receiver string
}

type Update struct {
	Id      int
	Message string
}

// CreateHuddle in huddles table
func CreateHuddle(db *gorm.DB, username string, t *NewHuddle) error {
	huddle := Huddle{
		CreatedBy: username,
		Message:   t.Message,
	}
	if err := db.Table("huddles").Create(&huddle).Error; err != nil {
		return err
	}

	var acceptedInvites []Invite
	if err := db.
		Table("invites").
		Where("(sender = ? OR receiver = ?) AND accepted = 1", username, username).
		Find(&acceptedInvites).Error; err != nil {
		return err
	}

	var hiddenUsernames []string
	if err := db.
		Table("hides").
		Select("hidden").
		Where("user = ?", username).
		Find(&hiddenUsernames).
		Error; err != nil {
		return err
	}

	usernames := GetNewHuddleUsernamesFromInvites(acceptedInvites, hiddenUsernames, username)

	var notifyUsernames []string
	if err := db.
		Table("users").
		Select("username").
		Where("username IN ? AND new_huddles_notifications = 1", usernames).
		Find(&notifyUsernames).
		Error; err != nil {
		return err
	}

	tokens, getErr := service.GetTokensByUsernames(db, notifyUsernames)
	if getErr != nil {
		return nil
	}

	hangoutNotification := service.FcmNotification{
		Title:   t.Name + " posted",
		Body:    t.Message,
		Devices: tokens,
	}

	return service.SendNotification(&hangoutNotification)
}

// GetHuddle from huddles table
func GetHuddle(db *gorm.DB, huddleId string, username string) (HuddleData, error) {
	var huddle Huddle
	var user p.Person
	var huddleComments []HuddleComment
	var huddleLikes []HuddleLike

	if err := db.Table("huddles").Where("id = ?", huddleId).First(&huddle).Error; err != nil {
		return HuddleData{}, err
	}

	if err := db.
		Table("users").
		Select("firstname, profile_photo").
		Where("username = ?", huddle.CreatedBy).
		First(&user).
		Error; err != nil {
		return HuddleData{}, err
	}

	if err := db.
		Table("huddles_likes").
		Where("huddle_id = ?", huddle.Id).
		Find(&huddleLikes).Error; err != nil {
		return HuddleData{}, err
	}

	if err := db.
		Table("huddles_comments").
		Where("huddle_id = ?", huddle.Id).
		Find(&huddleComments).Error; err != nil {
		return HuddleData{}, err
	}

	return HuddleData{
		Id:             int(huddle.Id),
		CreatedBy:      huddle.CreatedBy,
		Message:        huddle.Message,
		Name:           user.Firstname,
		ProfilePhoto:   user.ProfilePhoto,
		CommentsNumber: len(huddleComments),
		LikesNumber:    len(huddleLikes),
		Liked:          isHuddleLiked(huddleLikes, username, int(huddle.Id)),
	}, nil
}

// GetUserHuddles from huddles table
func GetUserHuddles(db *gorm.DB, username string, lastId string) ([]HuddleData, error) {
	var huddlesData []HuddleData
	var huddles []Huddle
	var profiles []p.Person

	var idCondition string
	if lastId != "" {
		idCondition = fmt.Sprintf("id < %s AND ", lastId)
	}

	if err := db.
		Table("huddles").
		Where(idCondition+"created_by = ?", username).
		Order("created DESC").
		Limit(20).
		Find(&huddles).Error; err != nil {
		return huddlesData, err
	}

	users := GetUsernamesFromHuddles(huddles)
	if err := db.Table("users").Where("username IN ?", users).Find(&profiles).Error; err != nil {
		return huddlesData, err
	}

	for _, huddle := range huddles {
		profileInfo := GetProfileInfoFromProfiles(profiles, huddle.CreatedBy)

		var liked int
		if huddle.CreatedBy != username {
			liked = 1
		}

		huddlesData = append(huddlesData, HuddleData{
			Id:           int(huddle.Id),
			CreatedBy:    huddle.CreatedBy,
			Name:         profileInfo.Firstname,
			ProfilePhoto: profileInfo.ProfilePhoto,
			Message:      huddle.Message,
			Liked:        liked,
		})

	}

	return huddlesData, nil
}

// DeleteHuddle from huddles table
func DeleteHuddle(db *gorm.DB, id string) error {
	if err := db.Table("huddles").Where("id = ?", id).Delete(&Huddle{}).Error; err != nil {
		return err
	}

	if err := db.Table("huddles_likes").Where("huddle_id = ?", id).Delete(&HuddleLike{}).Error; err != nil {
		return err
	}

	return nil
}

// GetNewHuddleUsernamesFromInvites from invites array
func GetNewHuddleUsernamesFromInvites(
	invites []Invite,
	hiddenUsernames []string,
	username string,
) []string {
	var usernames []string
	for _, invite := range invites {
		var user string

		if invite.Sender == username {
			user = invite.Receiver
		} else {
			user = invite.Sender
		}

		if !p.IsPersonHidden(hiddenUsernames, user) {
			usernames = append(usernames, user)
		}
	}
	return usernames
}

// Get ids from huddles array
func GetIdsFromHuddlesArray(huddles []Huddle) []int {
	var ids []int

	for _, huddle := range huddles {
		ids = append(ids, int(huddle.Id))
	}

	return ids
}

// Get usernames from huddles array
func GetUsernamesFromHuddles(huddles []Huddle) []string {
	var usernames []string

	for _, h := range huddles {
		if !contains(usernames, h.CreatedBy) {
			usernames = append(usernames, h.CreatedBy)
		}
	}

	return usernames
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

// GetProfileInfoFromProfiles returns profile
func GetProfileInfoFromProfiles(profiles []people.Person, username string) p.Person {
	var person p.Person

	for _, profile := range profiles {
		if profile.Username == username {
			person = profile
		}
	}

	return person
}

func isHuddleLiked(likedHuddles []HuddleLike, username string, huddleId int) int {
	for _, likedHuddle := range likedHuddles {
		if likedHuddle.HuddleId == huddleId && likedHuddle.Sender == username {
			return 1
		}
	}

	return 0
}
