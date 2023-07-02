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

// Add Huddle comment to huddles_comments table
func AddHuddleComment(db *gorm.DB, t *HuddleComment) error {
	var createdBy string
	var name string

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

	if err := db.
		Table("users").
		Select("firstname").
		Where("username = ?", t.Sender).
		Find(&name).
		Error; err != nil {
		return err
	}

	tokens := []string{}
	if err := service.GetTokensByUsername(db, &tokens, createdBy); err != nil {
		return nil
	}

	fcmNotification := service.FcmNotification{
		Data: map[string]interface{}{
			"type":     huddleType,
			"huddleId": t.HuddleId,
		},
		Title:   name + " commented",
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

	tokens := []string{}
	if err := service.GetTokensByUsername(db, &tokens, t.Receiver); err != nil {
		return nil
	}

	fcmNotification := service.FcmNotification{
		Data: map[string]interface{}{
			"type": huddleType,
		},
		Title:   name + " mentioned you",
		Body:    t.Message,
		Sound:   "default",
		Devices: tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// Get Huddle comments from huddles_comments table
func GetHuddleComments(
	db *gorm.DB,
	huddleId string,
	username string,
	lastId string) ([]HuddleCommentData, []Mention, error) {
	var comments []HuddleCommentData
	var mentions []Mention
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
		return comments, mentions, err
	}

	var mentionsUsernames []string
	if err := db.
		Table("huddles_comments").
		Distinct().
		Select("sender").
		Where("huddle_id = ?", huddleId).
		Find(&mentionsUsernames).
		Error; err != nil {
		return comments, mentions, err
	}

	commentsUsernames := getCommentsUsernames(huddleComments)
	if err := db.
		Table("users").
		Where("username IN ?", commentsUsernames).
		Find(&people).
		Error; err != nil {
		return comments, mentions, err
	}

	if err := db.
		Table("huddles_comments_likes").
		Where("huddle_id = ?", huddleId).
		Find(&likes).
		Error; err != nil {
		return comments, mentions, err
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

	mentions = getMentions(mentionsUsernames, people)

	return comments, mentions, nil
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

func getMentions(usernames []string, people []p.Person) []Mention {
	var mentions []Mention

	for _, username := range usernames {
		for _, person := range people {
			if username == person.Username {
				mentions = append(mentions, Mention{
					Username:     person.Username,
					Name:         person.Firstname,
					ProfilePhoto: person.ProfilePhoto,
				})

				break
			}
		}
	}

	return mentions
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
