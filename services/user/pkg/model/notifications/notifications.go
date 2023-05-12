package notifications

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

const timeFormat = "2006-01-02 15:04:05"

const PersonInviteType = "person_invite"
const PersonInviteAcceptType = "person_invite_accept"
const HuddleInteractType = "huddle_interact"
const HuddleConfirmType = "huddle_confirm"
const CommentType = "comment"
const CommentMentionType = "comment_mention"
const CommentLikeType = "comment_like"

type Notification struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Sender   string
	Receiver string
	EventId  int
	Type     string `gorm:"type:enum('person_invite', 'person_invite_accept', 'huddle_interact', 'huddle_confirm', 'comment', 'comment_mention', 'comment_like')"`
	Seen     int    `gorm:"default:0"`
	Created  int64  `gorm:"autoCreateTime"`
}

func (Notification) TableName() string {
	return "notifications"
}

type NotificationData struct {
	Id           int    `json:"id"`
	EventId      int    `json:"eventId"`
	Sender       string `json:"sender"`
	SenderName   string `json:"senderName"`
	ProfilePhoto string `json:"profilePhoto"`
	Type         string `json:"type"`
	What         string `json:"what,omitempty"`
	Accepted     int    `json:"accepted,omitempty"`
	Confirmed    int    `json:"confirmed,omitempty"`
	Comment      string `json:"comment,omitempty"`
	Created      string `json:"created"`
}

type Person struct {
	Username     string `json:"username"`
	Firstname    string `json:"firstname"`
	ProfilePhoto string `json:"profilePhoto"`
}

type GetHuddles struct {
	Id   int
	What string
}

type GetComments struct {
	Id       int
	HuddleId int
	Message  string
}

// Get notifications from notifications_people and notifications_huddles tables
func GetNotifications(db *gorm.DB, username string) ([]NotificationData, error) {
	var notificationsData []NotificationData

	if err := db.Table("notifications").Where("receiver = ?", username).Update("seen", 1).Error; err != nil {
		return []NotificationData{}, nil
	}

	var notifications []Notification
	if err := db.Table("notifications").Where("receiver = ?", username).Find(&notifications).Error; err != nil {
		return []NotificationData{}, err
	}

	invitesIds := getInvitesIdsFromNotifications(notifications)
	huddlesIds := getHuddlesIdsFromNotifications(notifications)
	commentsIds := getCommentsIdsFromNotifications(notifications)
	usernames := getUsernamesFromNotifications(notifications)

	var acceptedInvites []int
	if len(invitesIds) > 0 {
		if err := db.
			Table("invites").
			Select("id").
			Where("id IN ? AND accepted = 1", invitesIds).
			Find(&acceptedInvites).
			Error; err != nil {
			return []NotificationData{}, err
		}
	}

	var whats []GetHuddles
	var confirmedHuddles []int
	if len(huddlesIds) > 0 {
		if err := db.
			Table("huddles").
			Select("id, what").
			Where("id IN ?", huddlesIds).
			Find(&whats).
			Error; err != nil {
			return []NotificationData{}, err
		}

		if err := db.
			Table("huddles_interacted").
			Select("huddle_id").
			Where("huddle_id IN ? AND confirmed = 1", huddlesIds).
			Find(&confirmedHuddles).
			Error; err != nil {
			return []NotificationData{}, err
		}
	}

	var comments []GetComments
	if len(commentsIds) > 0 {
		if err := db.
			Table("huddles_comments").
			Select("id, huddle_id, message").
			Where("id IN ?", commentsIds).
			Find(&comments).
			Error; err != nil {
			return []NotificationData{}, err
		}
	}

	var profiles []Person
	if err := db.
		Table("users").
		Select("username, firstname, profile_photo").
		Where("username IN ?", usernames).
		Find(&profiles).
		Error; err != nil {
		return []NotificationData{}, err
	}

	for i, notification := range notifications {
		var accepted int
		var what string
		var confirmed int
		var comment string

		eventId := notification.EventId

		name, profilePhoto := getProfileInfo(profiles, notification.Sender)

		if notification.Type == HuddleInteractType {
			what = getWhat(whats, notification.EventId)
			confirmed = getConfirmed(confirmedHuddles, notification.EventId)
		}

		if notification.Type == HuddleConfirmType {
			what = getWhat(whats, notification.EventId)
		}

		if notification.Type == CommentType ||
			notification.Type == CommentMentionType ||
			notification.Type == CommentLikeType {
			// In case of comment type add huddle id as event id to payload
			var huddleId int
			huddleId, comment = getCommentMessage(comments, notification.EventId)

			eventId = huddleId
		}

		if notification.Type == PersonInviteType {
			accepted = getIfAccepted(acceptedInvites, notification.EventId)
		}

		time := time.Unix(notification.Created, 0).Format(timeFormat)

		notificationsData = append(notificationsData, NotificationData{
			Id:           i,
			EventId:      eventId,
			Sender:       notification.Sender,
			SenderName:   name,
			ProfilePhoto: profilePhoto,
			Type:         notification.Type,
			What:         what,
			Accepted:     accepted,
			Confirmed:    confirmed,
			Comment:      comment,
			Created:      time,
		})
	}

	return notificationsData, nil
}

func getInvitesIdsFromNotifications(notifications []Notification) []int {
	var inviteIds []int

	for _, notification := range notifications {
		if strings.Contains(notification.Type, PersonInviteType) {
			inviteIds = append(inviteIds, notification.EventId)
		}
	}

	return inviteIds
}

func getHuddlesIdsFromNotifications(notifications []Notification) []int {
	var huddleIds []int
	for _, notification := range notifications {
		if strings.Contains(notification.Type, "huddle") {
			huddleIds = append(huddleIds, notification.EventId)
		}
	}

	return huddleIds
}

func getCommentsIdsFromNotifications(notifications []Notification) []int {
	var commentsIds []int

	for _, notification := range notifications {
		if strings.Contains(notification.Type, "comment") {
			commentsIds = append(commentsIds, notification.EventId)
		}
	}

	return commentsIds
}

func getUsernamesFromNotifications(notifications []Notification) []string {
	var usernames []string

	for _, notification := range notifications {
		usernames = append(usernames, notification.Sender)
	}

	return usernames
}

func getProfileInfo(profiles []Person, username string) (string, string) {
	var name string
	var profilePhoto string

	for _, profile := range profiles {
		if profile.Username == username {
			name = profile.Firstname
			profilePhoto = profile.ProfilePhoto

			break
		}
	}

	return name, profilePhoto
}

func getWhat(huddles []GetHuddles, huddleId int) string {
	var what string

	for _, huddle := range huddles {
		if huddle.Id == huddleId {
			what = huddle.What

			break
		}
	}

	return what
}

func getCommentMessage(comments []GetComments, commentId int) (int, string) {
	var huddleId int
	var message string

	for _, comment := range comments {
		if comment.Id == commentId {
			huddleId = comment.HuddleId
			message = comment.Message

			break
		}
	}

	return huddleId, message
}

func getConfirmed(confirmedHuddles []int, huddleId int) int {
	var confirmed int

	for _, confirmedHuddle := range confirmedHuddles {
		if confirmedHuddle == huddleId {
			confirmed = 1

			break
		}
	}

	return confirmed
}

func getIfAccepted(acceptedInvites []int, inviteId int) int {
	var accepted int

	for _, acceptedInvite := range acceptedInvites {
		if acceptedInvite == inviteId {
			accepted = 1

			break
		}
	}

	return accepted
}
