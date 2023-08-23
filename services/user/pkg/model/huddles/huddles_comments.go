package huddles

import (
	"fmt"

	p "github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

type HuddleComment struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Sender   string
	HuddleId int
	Message  string
	Mention  *string
	Created  int64 `gorm:"autoCreateTime"`
}

func (HuddleComment) TableName() string {
	return "huddles_comments"
}

type MentionComment struct {
	Receiver string
	HuddleId int
	Message  string
}

type Mention struct {
	Username     string `json:"username,omitempty"`
	Name         string `json:"name"`
	ProfilePhoto string `json:"profilePhoto,omitempty"`
}

type HuddleCommentData struct {
	Id           int      `json:"id"`
	Sender       string   `json:"sender"`
	Name         string   `json:"name"`
	ProfilePhoto string   `json:"profilePhoto,omitempty"`
	Message      string   `json:"message"`
	Mention      *Mention `json:"mention,omitempty"`
	LikesNumber  int      `json:"likesNumber,omitempty"`
	Liked        int      `json:"liked,omitempty"`
	Time         int64    `json:"time"`
}

// AddHuddleComment to huddles_comments table
func AddHuddleComment(db *gorm.DB, t *HuddleComment) error {
	var createdBy string

	if err := db.Table("huddles_comments").Create(&t).Error; err != nil {
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

	if t.Sender == createdBy {
		return nil
	}

	var commentsNotification int
	if err := db.
		Table("users").
		Select("comments_notifications").
		Where("username = ?", createdBy).
		Find(&commentsNotification).
		Error; err != nil {
		return err
	}

	if commentsNotification != 1 {
		return nil
	}

	tokens, err := service.GetTokensByUsername(db, createdBy)
	if err != nil {
		return err
	}

	fcmNotification := service.FcmNotification{
		Title:   t.Sender + " commented your post",
		Body:    t.Message,
		Sound:   "default",
		Devices: tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// AddHuddleMentionComment to huddles_comments table
func AddHuddleMentionComment(db *gorm.DB, username string, t *MentionComment) error {
	var name string

	comment := HuddleComment{
		Sender:   username,
		HuddleId: t.HuddleId,
		Message:  t.Message,
		Mention:  &t.Receiver,
	}

	if err := db.Table("huddles_comments").Create(&comment).Error; err != nil {
		return err
	}

	if username == t.Receiver {
		return nil
	}

	var commentsNotification int
	if err := db.
		Table("users").
		Select("mentions_notifications").
		Where("username = ?", t.Receiver).
		Find(&commentsNotification).
		Error; err != nil {
		return err
	}

	if commentsNotification != 1 {
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

	tokens, err := service.GetTokensByUsername(db, t.Receiver)
	if err != nil {
		return err
	}

	fcmNotification := service.FcmNotification{
		Title:   name + " mentioned you in comment",
		Body:    t.Message,
		Sound:   "default",
		Devices: tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// GetHuddleComments from huddles_comments table
func GetHuddleComments(
	db *gorm.DB,
	huddleId string,
	username string,
	lastId string) ([]HuddleCommentData, error) {
	var comments []HuddleCommentData
	var huddleComments []HuddleComment
	var people []p.Person
	var likes []HuddleCommentLike

	var idCondition string
	if lastId != "" {
		idCondition = fmt.Sprintf("id > %s AND ", lastId)
	}

	if err := db.
		Table("huddles_comments").
		Where(idCondition+"huddle_id = ?", huddleId).
		Limit(15).
		Find(&huddleComments).
		Error; err != nil {
		return comments, err
	}

	commentsUsernames := getCommentsUsernames(huddleComments)
	if err := db.
		Table("users").
		Where("username IN ?", commentsUsernames).
		Find(&people).
		Error; err != nil {
		return comments, err
	}

	if err := db.
		Table("huddles_comments_likes").
		Where("huddle_id = ?", huddleId).
		Find(&likes).
		Error; err != nil {
		return comments, err
	}

	for _, comment := range huddleComments {
		var mention *Mention

		user := getCommentUser(comment.Sender, people)
		usersLiked := getLikesNumberPerComment(likes, int(comment.Id))
		liked := liked(usersLiked, username)

		if comment.Mention != nil {
			mention = getMention(*comment.Mention, people)
		}

		comments = append(comments, HuddleCommentData{
			Id:           int(comment.Id),
			Sender:       comment.Sender,
			Name:         user.Firstname,
			ProfilePhoto: user.ProfilePhoto,
			Message:      comment.Message,
			LikesNumber:  len(usersLiked),
			Liked:        liked,
			Mention:      mention,
			Time:         comment.Created,
		})
	}

	return comments, nil
}

// GetMentions from users and invites tables
func GetMentions(db *gorm.DB, username string) ([]p.Person, error) {
	var mentions []p.Person
	var invites []p.Invite

	if err := db.
		Table("invites").
		Where("(sender = ? OR receiver = ?) AND accepted = 1",
			username, username).
		Order("id DESC").
		Find(&invites).Error; err != nil {
		return nil, err
	}

	usernames := p.GetUsernamesFromInvites(invites, username)

	if err := db.
		Table("users").
		Where("username IN ?", usernames).
		Find(&mentions).
		Error; err != nil {
		return nil, err
	}

	return mentions, nil
}

// DeleteHuddleComment from huddles comments table
func DeleteHuddleComment(db *gorm.DB, id string) error {
	return db.Table("huddles_comments").Where("id = ?", id).Delete(&HuddleComment{}).Error
}

func getCommentsUsernames(huddleComments []HuddleComment) []string {
	var usernames []string

	for _, comment := range huddleComments {
		if !contains(usernames, comment.Sender) {
			usernames = append(usernames, comment.Sender)
		}
		if comment.Mention != nil {
			if !contains(usernames, *comment.Mention) {
				usernames = append(usernames, *comment.Mention)
			}
		}
	}

	return usernames
}

func getCommentUser(username string, people []p.Person) p.Person {
	for _, user := range people {
		if user.Username == username {
			return user
		}
	}

	return p.Person{}
}

func getMention(mention string, people []p.Person) *Mention {
	for _, user := range people {
		if user.Username == mention {
			return &Mention{
				Name:         user.Firstname,
				ProfilePhoto: user.ProfilePhoto,
			}
		}
	}

	return nil
}

func getLikesNumberPerComment(likes []HuddleCommentLike, commentId int) []string {
	var usersLiked []string

	for _, like := range likes {
		if like.CommentId == commentId {
			usersLiked = append(usersLiked, like.Sender)
		}
	}

	return usersLiked
}

func liked(usersLiked []string, username string) int {
	for _, userLiked := range usersLiked {
		if userLiked == username {
			return 1
		}
	}

	return 0
}
