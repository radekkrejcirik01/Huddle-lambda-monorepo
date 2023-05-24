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

type Like struct {
	Sender    string
	Receiver  string
	CommentId int
	HuddleId  int
}

type Liker struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ProfilePhoto string `json:"profilePhoto,omitempty"`
}

// Add Huddle comment like to huddles_comments_likes table
func LikeHuddleComment(db *gorm.DB, t *Like) error {
	var name string

	like := HuddleCommentLike{
		Sender:    t.Sender,
		CommentId: t.CommentId,
		HuddleId:  t.HuddleId,
	}

	if err := db.Table("huddles_comments_likes").Create(&like).Error; err != nil {
		return err
	}

	if err := db.
		Table("users").
		Select("firstname").
		Where("username = ?", t.Sender).
		Find(&name).
		Error; err != nil {
		return err
	}

	tokens := []string{}
	if err := service.GetTokensByUsername(db, &tokens, t.Receiver); err != nil {
		return nil
	}

	fcmNotification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    "comment",
		Body:    name + " liked your comment",
		Sound:   "default",
		Devices: tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// Get Huddle comment likes from huddles_comments_likes table
func GetCommentLikes(db *gorm.DB, commentId int, lastId string) ([]Liker, error) {
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

// Remove Huddle comment like from huddles_comments_likes table
func RemoveHuddleCommentLike(db *gorm.DB, id int, sender string) error {
	return db.
		Table("huddles_comments_likes").
		Where("comment_id = ? AND sender = ?", id, sender).
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
