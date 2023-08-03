package huddles

import (
	"fmt"

	p "github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

type HuddleCommentLike struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	Sender    string
	CommentId int
	HuddleId  int
}

func (HuddleCommentLike) TableName() string {
	return "huddles_comments_likes"
}

type CommentLike struct {
	Receiver  string
	CommentId int
	HuddleId  int
}

type Liker struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ProfilePhoto string `json:"profilePhoto,omitempty"`
}

// LikeHuddleComment in huddles_comments_likes table
func LikeHuddleComment(db *gorm.DB, username string, t *CommentLike) error {
	var name string

	like := HuddleCommentLike{
		Sender:    username,
		CommentId: t.CommentId,
		HuddleId:  t.HuddleId,
	}

	if err := db.Table("huddles_comments_likes").Create(&like).Error; err != nil {
		return err
	}

	if username == t.Receiver {
		return nil
	}

	if err := db.
		Table("users").
		Select("firstname").
		Where("username = ?", username).
		Find(&name).
		Error; err != nil {
		return err
	}

	tokens := []string{}
	if err := service.GetTokensByUsername(db, &tokens, t.Receiver); err != nil {
		return nil
	}

	fcmNotification := service.FcmNotification{
		Data: map[string]interface{}{
			"type":     huddleType,
			"huddleId": t.HuddleId,
		},
		Body:    name + " liked your comment",
		Sound:   "default",
		Devices: tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// GetCommentLikes from huddles_comments_likes table
func GetCommentLikes(db *gorm.DB, commentId string, lastId string) ([]Liker, error) {
	var likes []HuddleCommentLike
	var likers []Liker
	var profiles []p.Person

	var idCondition string
	if lastId != "" {
		idCondition = fmt.Sprintf("id > %s AND ", lastId)
	}

	if err := db.
		Table("huddles_comments_likes").
		Where(idCondition+"comment_id = ?", commentId).
		Limit(10).
		Find(&likes).
		Error; err != nil {
		return nil, err
	}

	usernames := getUsernamesFromLikes(likes)

	if err := db.
		Table("users").
		Select("username, firstname, profile_photo").
		Where("username in ?", usernames).
		Find(&profiles).
		Error; err != nil {
		return nil, err
	}

	for _, like := range likes {
		for _, profile := range profiles {
			if profile.Username == like.Sender {
				likers = append(likers, Liker{
					Id:           int(like.Id),
					Name:         profile.Firstname,
					ProfilePhoto: profile.ProfilePhoto,
				})
				break
			}
		}
	}

	return likers, nil
}

// RemoveHuddleCommentLike from huddles_comments_likes table
func RemoveHuddleCommentLike(db *gorm.DB, commentId string, username string) error {
	return db.
		Table("huddles_comments_likes").
		Where("comment_id = ? AND sender = ?", commentId, username).
		Delete(&HuddleCommentLike{}).
		Error
}

func getUsernamesFromLikes(likes []HuddleCommentLike) []string {
	var usernames []string
	for _, like := range likes {
		usernames = append(usernames, like.Sender)
	}
	return usernames
}
