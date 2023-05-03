package huddles

import (
	"time"

	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

const timeFormat = "2006-01-02 15:04:05"

type HuddleComment struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Sender   string
	HuddleId uint
	Message  string
	Created  int64 `gorm:"autoCreateTime"`
}

type AddComment struct {
	Sender   string
	HuddleId uint
	Message  string
	Mentions []string
}

type HuddleCommentData struct {
	Id           uint   `json:"id"`
	Sender       string `json:"sender"`
	Name         string `json:"name"`
	ProfilePhoto string `json:"profilePhoto"`
	Message      string `json:"message"`
	Time         string `json:"time"`
}

func (HuddleComment) TableName() string {
	return "huddles_comments"
}

// Add Huddle comment to huddles_comments table
func AddHuddleComment(db *gorm.DB, t *AddComment) error {
	var createdBy string
	var name string

	huddleComment := HuddleComment{
		Sender:   t.Sender,
		HuddleId: t.HuddleId,
		Message:  t.Message,
	}

	if err := db.Table("huddles_comments").Create(&huddleComment).Error; err != nil {
		return err
	}

	if err := db.
		Table("huddles").
		Select("created_by").
		Where("id = ?", t.HuddleId).
		Find(&createdBy).
		Error; err != nil {
		return err
	}

	if t.Sender == createdBy && len(t.Mentions) == 0 {
		return nil
	}

	if err := db.
		Table("users").
		Select("firstname").
		Where("username = ?", t.Sender).
		Find(&name).
		Error; err != nil {
		return err
	}

	usernames := []string{createdBy}
	if len(t.Mentions) > 0 {
		usernames = t.Mentions
	}

	tokens, err := service.GetTokensByUsernames(db, usernames)
	if err != nil {
		return nil
	}

	title := name + " added a comment"
	if len(t.Mentions) > 0 {
		title = name + " mentioned you"
	}

	notification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    "comment",
		Title:   title,
		Body:    t.Message,
		Sound:   "default",
		Devices: tokens,
	}

	return service.SendNotification(&notification)
}

// Get Huddle comments from huddles_comments table
func GetHuddleComments(db *gorm.DB, huddleId uint) ([]HuddleCommentData, error) {
	var comments []HuddleCommentData
	var huddleComments []HuddleComment
	var people []people.Person

	if err := db.
		Table("huddles_comments").
		Where("huddle_id = ?", huddleId).
		Find(&huddleComments).
		Error; err != nil {
		return comments, err
	}

	commentsUsernames := getUsernamesFromComments(huddleComments)
	if err := db.
		Table("users").
		Where("username IN ?", commentsUsernames).
		Find(&people).
		Error; err != nil {
		return comments, err
	}

	for _, comment := range huddleComments {
		for _, user := range people {
			if comment.Sender == user.Username {
				time := time.Unix(comment.Created, 0).Format(timeFormat)

				comments = append(comments, HuddleCommentData{
					Id:           comment.Id,
					Sender:       comment.Sender,
					Name:         user.Firstname,
					ProfilePhoto: user.ProfilePhoto,
					Message:      comment.Message,
					Time:         time,
				})

				break
			}
		}
	}

	return comments, nil
}

func getUsernamesFromComments(huddleComments []HuddleComment) []string {
	var usernames []string

	for _, comment := range huddleComments {
		if !contains(usernames, comment.Sender) {
			usernames = append(usernames, comment.Sender)
		}
	}

	return usernames
}
