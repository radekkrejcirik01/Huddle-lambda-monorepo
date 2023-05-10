package huddles

import (
	"errors"

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
	Where     string
	When      string
	Confirmed int   `gorm:"default:0"`
	Created   int64 `gorm:"autoCreateTime"`
	Canceled  int   `gorm:"default:0"`
}

func (Huddle) TableName() string {
	return "huddles"
}

type NewHuddle struct {
	Sender string
	What   string
	Where  string
	When   string
}

type HuddleData struct {
	Id             uint   `json:"id"`
	CreatedBy      string `json:"createdBy"`
	Name           string `json:"name"`
	ProfilePhoto   string `json:"profilePhoto"`
	What           string `json:"what"`
	Where          string `json:"where"`
	When           string `json:"when"`
	Interacted     int    `json:"interacted,omitempty"`
	Confirmed      int    `json:"confirmed,omitempty"`
	Canceled       int    `json:"canceled,omitempty"`
	CommentsNumber int    `json:"commentsNumber,omitempty"`
}

type Invite struct {
	Sender   string
	Receiver string
}

type Update struct {
	Id    uint
	What  string
	Where string
	When  string
}

type PostAgain struct {
	Id uint
}

// Add Huddle to huddles table
func AddHuddle(db *gorm.DB, t *NewHuddle) error {
	huddle := Huddle{
		CreatedBy: t.Sender,
		What:      t.What,
		Where:     t.Where,
		When:      t.When,
	}
	if err := db.Table("huddles").Create(&huddle).Error; err != nil {
		return err
	}

	var acceptedInvites []p.PeopleNotification
	if err := db.
		Table("notifications_people").
		Where("(sender = ? OR receiver = ?) AND type = 'person_invite' AND accepted = 1", t.Sender, t.Sender).
		Find(&acceptedInvites).Error; err != nil {
		return err
	}

	usernames := p.GetUsernamesFromAcceptedInvites(acceptedInvites, t.Sender)

	tokens, getErr := service.GetTokensByUsernames(db, usernames)
	if getErr != nil {
		return nil
	}

	hangoutNotification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    huddleType,
		Body:    t.Sender + " added a new Huddle: " + t.What,
		Devices: tokens,
	}

	return service.SendNotification(&hangoutNotification)
}

// Get user huddles from huddles table
func GetUserHuddles(db *gorm.DB, username string) ([]HuddleData, error) {
	var huddlesData []HuddleData
	var huddles []Huddle
	var profiles []p.Person

	query :=
		`
		SELECT
			*
		FROM
			huddles
		WHERE
			created_by = ?
			OR id IN(
				SELECT
					huddle_id FROM huddles_interacted
				WHERE
					sender = ?
					AND confirmed = 1)
		ORDER BY
			created DESC
		`
	if err := db.Raw(query, username, username).Find(&huddles).Error; err != nil {
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
			Id:           huddle.Id,
			CreatedBy:    huddle.CreatedBy,
			Name:         profileInfo.Firstname,
			ProfilePhoto: profileInfo.ProfilePhoto,
			What:         huddle.What,
			Where:        huddle.Where,
			When:         huddle.When,
			Interacted:   interacted,
			Confirmed:    huddle.Confirmed,
			Canceled:     huddle.Canceled,
		})

	}

	return huddlesData, nil
}

// Get Huddles from huddles table
func GetHuddles(db *gorm.DB, username string) ([]HuddleData, error) {
	var huddlesData []HuddleData

	var invites []Invite
	if err := db.
		Table("notifications_people").
		Where(
			"(sender = ? OR receiver = ?) AND type = 'person_invite' AND accepted = 1",
			username, username,
		).
		Find(&invites).Error; err != nil {
		return huddlesData, err
	}

	people := GetUsernamesFromInvites(invites, username)

	var huddles []Huddle
	if err := db.
		Table("huddles").
		Where("(created_by IN ? OR created_by = ?) AND confirmed = 0 AND canceled = 0", people, username).
		Order("created DESC").
		Find(&huddles).
		Error; err != nil {
		return huddlesData, err
	}

	if len(huddles) < 1 {
		return huddlesData, nil
	}

	huddlesIds := GetIdsFromHuddlesArray(huddles)
	var interactedHuddlesIds []uint
	if err := db.
		Table("huddles_interacted").
		Select("huddle_id").
		Where("sender = ? AND huddle_id IN ?", username, huddlesIds).
		Find(&interactedHuddlesIds).Error; err != nil {
		return huddlesData, err
	}

	users := GetUsernamesFromHuddles(huddles)
	var profiles []p.Person
	if err := db.Table("users").Where("username IN ?", users).Find(&profiles).Error; err != nil {
		return huddlesData, err
	}

	var comments []HuddleComment
	if err := db.
		Table("huddles_comments").
		Where("huddle_id IN ?", huddlesIds).
		Find(&comments).
		Error; err != nil {
		return huddlesData, err
	}

	for _, huddle := range huddles {
		profileInfo := GetProfileInfoFromProfiles(profiles, huddle.CreatedBy)
		interacted := GetInteraction(interactedHuddlesIds, huddle.Id)
		commentsNumber := getCommentsNumber(comments, huddle.Id)

		huddlesData = append(huddlesData, HuddleData{
			Id:             huddle.Id,
			CreatedBy:      huddle.CreatedBy,
			Name:           profileInfo.Firstname,
			ProfilePhoto:   profileInfo.ProfilePhoto,
			What:           huddle.What,
			Where:          huddle.Where,
			When:           huddle.When,
			Interacted:     interacted,
			CommentsNumber: commentsNumber,
		})

	}

	return huddlesData, nil
}

// Update Huddle in huddles table
func UpdateHuddle(db *gorm.DB, t *Update) error {
	update := map[string]interface{}{
		"what":  t.What,
		"where": t.Where,
		"when":  t.When,
	}

	return db.Table("huddles").Where("id = ?", t.Id).Updates(update).Error
}

// Update Huddle canceled column in huddles table
func PostHuddleAgain(db *gorm.DB, t *PostAgain) error {
	return db.Table("huddles").Where("id = ?", t.Id).Update("canceled", 0).Error
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
		Id:           huddle.Id,
		CreatedBy:    huddle.CreatedBy,
		Name:         profile.Firstname,
		ProfilePhoto: profile.ProfilePhoto,
		What:         huddle.What,
		Where:        huddle.Where,
		When:         huddle.When,
		Interacted:   interacted,
		Confirmed:    huddle.Confirmed,
		Canceled:     huddle.Canceled,
	}

	return huddleData, nil
}

func DeleteHuddle(db *gorm.DB, id uint) error {
	if err := db.Table("huddles").Where("id = ?", id).Delete(&Huddle{}).Error; err != nil {
		return err
	}

	if err := db.Table("huddles_interacted").Where("huddle_id = ?", id).Delete(&HuddleInteracted{}).Error; err != nil {
		return err
	}

	return nil
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

func getCommentsNumber(comments []HuddleComment, huddlesId uint) int {
	commentsNumber := 0

	for _, comment := range comments {
		if comment.HuddleId == huddlesId {
			commentsNumber = commentsNumber + 1
		}
	}

	return commentsNumber
}
