package huddles

import (
	"fmt"

	p "github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

type HuddleLike struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Sender   string
	HuddleId int
	Created  int64 `gorm:"autoCreateTime"`
}

func (HuddleLike) TableName() string {
	return "huddles_likes"
}

type Like struct {
	HuddleId int
	Message  string
	Receiver string
}

type UserLike struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ProfilePhoto string `json:"profilePhoto"`
}

// LikeHuddle to huddles_likes table
func LikeHuddle(db *gorm.DB, username string, t *Like) error {
	like := HuddleLike{
		Sender:   username,
		HuddleId: t.HuddleId,
	}
	if err := db.Table("huddles_likes").Create(&like).Error; err != nil {
		return err
	}

	if t.Receiver == username {
		return nil
	}

	var likesNotifications int
	if err := db.
		Table("users").
		Select("huddle_likes_notifications").
		Where("username = ?", t.Receiver).
		Find(&likesNotifications).
		Error; err != nil {
		return err
	}

	if likesNotifications != 1 {
		return nil
	}

	tokens, err := service.GetTokensByUsername(db, t.Receiver)
	if err != nil {
		return err
	}

	fcmNotification := service.FcmNotification{
		Title:   username + " liked your post",
		Body:    t.Message,
		Sound:   "default",
		Devices: tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// GetHuddleLikes from huddles_likes table
func GetHuddleLikes(db *gorm.DB, huddleId string, lastId string) ([]UserLike, error) {
	var usersLikes []UserLike
	var huddleLikes []HuddleLike
	var profiles []p.Person

	var idCondition string
	if lastId != "" {
		idCondition = fmt.Sprintf("id > %s AND ", lastId)
	}

	if err := db.
		Table("huddles_likes").
		Where(idCondition+"huddle_id = ?", huddleId).
		Limit(10).
		Find(&huddleLikes).
		Error; err != nil {
		return usersLikes, err
	}

	usernames := getUsernamesFromHuddleLikes(huddleLikes)

	if err := db.
		Table("users").
		Select("username, firstname, profile_photo").
		Where("username IN ?", usernames).
		Find(&profiles).
		Error; err != nil {
		return usersLikes, err
	}

	for _, huddleLike := range huddleLikes {
		for _, profile := range profiles {
			if profile.Username == huddleLike.Sender {
				usersLikes = append(usersLikes, UserLike{
					Id:           int(huddleLike.Id),
					Name:         profile.Firstname,
					ProfilePhoto: profile.ProfilePhoto,
				})
				break
			}
		}
	}

	return usersLikes, nil
}

func getUsernamesFromHuddleLikes(huddleLikes []HuddleLike) []string {
	var usernames []string
	for _, huddleLike := range huddleLikes {
		usernames = append(usernames, huddleLike.Sender)
	}

	return usernames
}

// RemoveHuddleLike from huddles_likes table
func RemoveHuddleLike(db *gorm.DB, username string, huddleId string) error {
	return db.
		Table("huddles_likes").
		Where("sender = ? AND huddle_id = ?", username, huddleId).
		Delete(&HuddleLike{}).
		Error
}
