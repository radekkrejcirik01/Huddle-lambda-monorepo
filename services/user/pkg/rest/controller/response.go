package controller

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/huddles"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/messaging"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/notifications"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/users"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Data    users.UserGet `json:"data,omitempty"`
}

type PeopleNumberResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	PeopleNumber int64  `json:"peopleNumber,omitempty"`
}

type PeopleResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message,omitempty"`
	Data    []people.Person `json:"data,omitempty"`
}

type GetInviteResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Data    people.Invite `json:"data,omitempty"`
}

type NotificationsResponse struct {
	Status  string                           `json:"status"`
	Message string                           `json:"message,omitempty"`
	Data    []notifications.NotificationData `json:"data,omitempty"`
}

type GetHuddlesResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message,omitempty"`
	Data    []huddles.HuddleData `json:"data,omitempty"`
}

type GetHuddleResponse struct {
	Status  string             `json:"status"`
	Message string             `json:"message,omitempty"`
	Data    huddles.HuddleData `json:"data,omitempty"`
}

type GetHuddleInteractionsResponse struct {
	Status        string                   `json:"status"`
	Message       string                   `json:"message,omitempty"`
	Data          []huddles.UserInteracted `json:"data,omitempty"`
	ConfirmedUser *huddles.UserInteracted  `json:"confirmedUser"`
}

type GetHuddleCommentsResponse struct {
	Status   string                      `json:"status"`
	Message  string                      `json:"message,omitempty"`
	Data     []huddles.HuddleCommentData `json:"data,omitempty"`
	Mentions []huddles.Mention           `json:"mentions,omitempty"`
}

type GetHuddleCommentsLikesResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message,omitempty"`
	Data    []huddles.Liker `json:"data,omitempty"`
}

type GetChatsResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message,omitempty"`
	Data    []messaging.Chat `json:"data,omitempty"`
}

type GetMessagesResponse struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message,omitempty"`
	Data    []messaging.MessageData `json:"data,omitempty"`
}

type GetMessagesByUsernamesResponse struct {
	Status         string                  `json:"status"`
	Message        string                  `json:"message,omitempty"`
	Data           []messaging.MessageData `json:"data,omitempty"`
	ConversationId uint                    `json:"conversationId,omitempty"`
}

type UploadPhotoResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message,omitempty"`
	ImageUrl string `json:"imageUrl,omitempty"`
}
