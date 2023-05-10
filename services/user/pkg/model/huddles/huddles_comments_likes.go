package huddles

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

const huddleCommentLikedType = "comment_liked"

type HuddleCommentLike struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	Sender    string
	CommentId uint
	HuddleId  uint
}

func (HuddleCommentLike) TableName() string {
	return "huddles_comments_likes"
}

type Like struct {
	Sender    string
	Receiver  string
	CommentId uint
	HuddleId  uint
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

	commentNotification := HuddleNotification{
		HuddleId: t.HuddleId,
		Sender:   t.Sender,
		Receiver: t.Receiver,
		Type:     huddleCommentLikedType,
	}

	if err := db.Table("notifications_huddles").Create(&commentNotification).Error; err != nil {
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

	notification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    "comment",
		Body:    name + " liked your comment",
		Sound:   "default",
		Devices: tokens,
	}

	return service.SendNotification(&notification)
}

// Remove Huddle comment like from huddles_comments_likes table
func RemoveHuddleCommentLike(db *gorm.DB, id int, sender string) error {
	return db.
		Table("huddles_comments_likes").
		Where("comment_id = ? AND sender = ?", id, sender).
		Delete(&HuddleCommentLike{}).
		Error
}
