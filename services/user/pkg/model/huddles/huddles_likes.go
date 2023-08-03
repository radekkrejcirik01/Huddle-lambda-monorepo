package huddles

import (
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
	Username     string `json:"username"`
	Firstname    string `json:"name"`
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

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return nil
	}

	fcmNotification := service.FcmNotification{
		Title:   username + " liked your huddle",
		Body:    t.Message,
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// GetHuddleLikes from huddles_likes table
func GetHuddleLikes(db *gorm.DB, huddleId string) ([]UserLike, error) {
	var usersLiked []UserLike

	var likes []string
	if err := db.
		Table("huddles_likes").
		Select("sender").
		Where("huddle_id = ?", huddleId).
		Find(&likes).
		Error; err != nil {
		return usersLiked, err
	}

	if err := db.
		Table("users").
		Where("username IN ?", likes).
		Find(&usersLiked).
		Error; err != nil {
		return usersLiked, err
	}

	return usersLiked, nil
}

// RemoveHuddleLike from huddles_likes table
func RemoveHuddleLike(db *gorm.DB, username string, huddleId string) error {
	return db.
		Table("huddles_likes").
		Where("sender = ? AND huddle_id = ?", username, huddleId).
		Delete(&HuddleLike{}).
		Error
}
