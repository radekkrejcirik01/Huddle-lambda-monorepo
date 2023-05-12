package huddles

import (
	n "github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/notifications"
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

	notification := n.Notification{
		EventId:  int(like.Id),
		Sender:   t.Sender,
		Receiver: t.Receiver,
		Type:     n.CommentLikeType,
	}

	if err := db.Table("notifications").Create(&notification).Error; err != nil {
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

// Remove Huddle comment like from huddles_comments_likes table
func RemoveHuddleCommentLike(db *gorm.DB, id int, sender string) error {
	return db.
		Table("huddles_comments_likes").
		Where("comment_id = ? AND sender = ?", id, sender).
		Delete(&HuddleCommentLike{}).
		Error
}
