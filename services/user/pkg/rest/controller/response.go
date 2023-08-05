package controller

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/huddles"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/messaging"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/users"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AuthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
}

type UserResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message,omitempty"`
	Data    users.UserData `json:"data,omitempty"`
}

type UserNotificationsResponse struct {
	Status  string             `json:"status"`
	Message string             `json:"message,omitempty"`
	Data    users.Notification `json:"data,omitempty"`
}

type PeopleResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message,omitempty"`
	Data    []people.PeopleData `json:"data,omitempty"`
}

type GetUnseenInvitesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Number  int64  `json:"number,omitempty"`
}

type GetIsConversationMutedResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message,omitempty"`
	Muted   bool     `json:"muted,omitempty"`
	People  []string `json:"people,omitempty"`
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

type GetHuddleLikesResponse struct {
	Status  string             `json:"status"`
	Message string             `json:"message,omitempty"`
	Data    []huddles.UserLike `json:"data,omitempty"`
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

type GetIsConversationLikedResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	IsLiked int    `json:"isLiked,omitempty"`
}

type UploadPhotoResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message,omitempty"`
	ImageUrl string `json:"imageUrl,omitempty"`
}
