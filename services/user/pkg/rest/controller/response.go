package controller

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/huddles"
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
	Status        string                         `json:"status"`
	Message       string                         `json:"message,omitempty"`
	Data          []huddles.HuddleInteractedData `json:"data,omitempty"`
	ConfirmedUser *string                        `json:"confirmedUser"`
}

type UploadPhotoResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message,omitempty"`
	ImageUrl string `json:"imageUrl,omitempty"`
}
