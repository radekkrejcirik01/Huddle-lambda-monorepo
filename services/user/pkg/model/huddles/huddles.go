package huddles

import (
	"errors"
	"fmt"
	"time"

	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	p "github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

const huddleType = "huddle"

// Huddle is a communication app for suggesting hangouts by adding simple posts called Huddles
type Huddle struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	CreatedBy string
	What      string
	Color     int
	Created   int64 `gorm:"autoCreateTime"`
}

func (Huddle) TableName() string {
	return "huddles"
}

type NewHuddle struct {
	Sender string
	What   string
	Color  int
}

type HuddleData struct {
	Id             int    `json:"id"`
	CreatedBy      string `json:"createdBy"`
	Name           string `json:"name"`
	ProfilePhoto   string `json:"profilePhoto"`
	What           string `json:"what"`
	Color          int    `json:"color"`
	Interacted     int    `json:"interacted,omitempty"`
	CommentsNumber int    `json:"commentsNumber,omitempty"`
}

type Invite struct {
	Sender   string
	Receiver string
}

type Update struct {
	Id   int
	What string
}

// Add Huddle to huddles table
func AddHuddle(db *gorm.DB, t *NewHuddle) error {
	huddle := Huddle{
		CreatedBy: t.Sender,
		What:      t.What,
		Color:     t.Color,
	}
	if err := db.Table("huddles").Create(&huddle).Error; err != nil {
		return err
	}

	var acceptedInvites []Invite
	if err := db.
		Table("invites").
		Where("(sender = ? OR receiver = ?) AND accepted = 1", t.Sender, t.Sender).
		Find(&acceptedInvites).Error; err != nil {
		return err
	}

	var hiddenUsernames []string
	if err := db.
		Table("hides").
		Select("hidden").
		Where("user = ?", t.Sender).
		Find(&hiddenUsernames).
		Error; err != nil {
		return err
	}

	usernames := GetNewHuddleUsernamesFromInvites(acceptedInvites, hiddenUsernames, t.Sender)

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
		Data: map[string]interface{}{
			"type":     huddleType,
			"huddleId": huddle.Id,
		},
		Body:    t.Sender + " added a new Huddle: " + t.What,
		Devices: tokens,
	}

	return service.SendNotification(&hangoutNotification)
}

// Get user huddles from huddles table
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

		var interacted int
		if huddle.CreatedBy != username {
			interacted = 1
		}

		huddlesData = append(huddlesData, HuddleData{
			Id:           int(huddle.Id),
			CreatedBy:    huddle.CreatedBy,
			Name:         profileInfo.Firstname,
			Color:        huddle.Color,
			ProfilePhoto: profileInfo.ProfilePhoto,
			What:         huddle.What,
			Interacted:   interacted,
		})

	}

	return huddlesData, nil
}

// Get Huddles from huddles table
func GetHuddles(db *gorm.DB, username string, lastId string) ([]HuddleData, error) {
	var huddlesData []HuddleData
	var invites []Invite
	var invitesUsernames []string
	var huddles []Huddle
	var huddleIds []int
	var huddlesUsernames []string

	if err := db.
		Table("invites").
		Where("(sender = ? OR receiver = ?) AND accepted = 1", username, username).
		Find(&invites).Error; err != nil {
		return huddlesData, err
	}

	var hiddenUsernames []string
	if err := db.
		Table("hides").
		Select("user").
		Where("hidden = ?", username).
		Find(&hiddenUsernames).
		Error; err != nil {
		return nil, err
	}

	var mutedUsernames []string
	if err := db.
		Table("muted_huddles").
		Select("muted").
		Where("user = ?", username).
		Find(&mutedUsernames).
		Error; err != nil {
		return nil, err
	}

	invitesUsernames = GetUsernamesFromInvites(invites, hiddenUsernames, mutedUsernames, username)

	// Two days ago in unix time
	t := time.Now().AddDate(0, 0, -2).Unix()

	var idCondition string
	if lastId != "" {
		idCondition = fmt.Sprintf("id < %s AND ", lastId)
	}
	if err := db.
		Table("huddles").
		Where(idCondition+"(created_by IN ? OR created_by = ?) AND created > ?",
			invitesUsernames, username, t).
		Order("created DESC").
		Limit(10).
		Find(&huddles).
		Error; err != nil {
		return huddlesData, err
	}

	if len(huddles) < 1 {
		return huddlesData, nil
	}

	huddleIds = GetIdsFromHuddlesArray(huddles)

	var interactedHuddlesIds []int
	if err := db.
		Table("huddles_interacted").
		Select("huddle_id").
		Where("sender = ? AND huddle_id IN ?", username, huddleIds).
		Find(&interactedHuddlesIds).Error; err != nil {
		return huddlesData, err
	}

	huddlesUsernames = GetUsernamesFromHuddles(huddles)

	var profiles []p.Person
	if err := db.
		Table("users").
		Where("username IN ?", huddlesUsernames).
		Find(&profiles).
		Error; err != nil {
		return huddlesData, err
	}

	var comments []HuddleComment
	if err := db.
		Table("huddles_comments").
		Where("huddle_id IN ?", huddleIds).
		Find(&comments).
		Error; err != nil {
		return huddlesData, err
	}

	for _, huddle := range huddles {
		profileInfo := GetProfileInfoFromProfiles(profiles, huddle.CreatedBy)
		interacted := GetInteraction(interactedHuddlesIds, int(huddle.Id))
		commentsNumber := getCommentsNumber(comments, int(huddle.Id))

		huddlesData = append(huddlesData, HuddleData{
			Id:             int(huddle.Id),
			CreatedBy:      huddle.CreatedBy,
			Name:           profileInfo.Firstname,
			ProfilePhoto:   profileInfo.ProfilePhoto,
			What:           huddle.What,
			Color:          huddle.Color,
			Interacted:     interacted,
			CommentsNumber: commentsNumber,
		})

	}

	return huddlesData, nil
}

// Update Huddle in huddles table
func UpdateHuddle(db *gorm.DB, t *Update) error {
	update := map[string]interface{}{
		"what": t.What,
	}

	return db.Table("huddles").Where("id = ?", t.Id).Updates(update).Error
}

// Get Huddle from huddles table by id
func GetHuddleById(db *gorm.DB, id uint, username string) (HuddleData, error) {
	var huddleData HuddleData
	var interacted int

	var huddle Huddle
	if err := db.
		Table("huddles").
		Where("id = ?", id).
		First(&huddle).Error; err != nil {
		return HuddleData{}, err
	}

	var profile p.Person
	if err := db.
		Table("users").
		Where("username = ?", huddle.CreatedBy).
		First(&profile).
		Error; err != nil {
		return HuddleData{}, err
	}

	if huddle.CreatedBy != username {
		var interaction HuddleInteracted

		err := db.
			Table("huddles_interacted").
			Where("sender = ? AND huddle_id = ?", username, id).
			First(&interaction).
			Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return HuddleData{}, err
		}

		if err == nil {
			interacted = 1
		}
	}

	huddleData = HuddleData{
		Id:           int(huddle.Id),
		CreatedBy:    huddle.CreatedBy,
		Name:         profile.Firstname,
		Color:        huddle.Color,
		ProfilePhoto: profile.ProfilePhoto,
		What:         huddle.What,
		Interacted:   interacted,
	}

	return huddleData, nil
}

// Delete huddle from huddles table
func DeleteHuddle(db *gorm.DB, id uint) error {
	if err := db.Table("huddles").Where("id = ?", id).Delete(&Huddle{}).Error; err != nil {
		return err
	}

	if err := db.Table("huddles_interacted").Where("huddle_id = ?", id).Delete(&HuddleInteracted{}).Error; err != nil {
		return err
	}

	return nil
}

// Get new huddles usernames from invites array
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

// Get usernames from invites array
func GetUsernamesFromInvites(
	invites []Invite,
	hiddenUsernames []string,
	mutedUsernames []string,
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

		if !p.IsPersonHidden(hiddenUsernames, user) && !p.IsPersonMuted(mutedUsernames, user) {
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

// Get profile info from profiles array
func GetProfileInfoFromProfiles(profiles []people.Person, username string) p.Person {
	var person p.Person

	for _, profile := range profiles {
		if profile.Username == username {
			person = profile
		}
	}

	return person
}

// Get interaction value from interacted huddles ids and huddle id
func GetInteraction(interactedHuddlesIds []int, huddleId int) int {
	for _, id := range interactedHuddlesIds {
		if id == huddleId {
			return 1
		}
	}

	return 0
}

// Get number of comments for Huddle
func getCommentsNumber(comments []HuddleComment, huddleId int) int {
	commentsNumber := 0

	for _, comment := range comments {
		if comment.HuddleId == huddleId {
			commentsNumber = commentsNumber + 1
		}
	}

	return commentsNumber
}
