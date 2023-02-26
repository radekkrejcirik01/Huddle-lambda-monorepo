package controller

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/hangouts"
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
	Data    []people.People `json:"data,omitempty"`
}

type CheckInvitationsResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message,omitempty"`
	Data    people.CheckIfFriend `json:"data,omitempty"`
}

type HangoutsResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message,omitempty"`
	Data    []hangouts.Hangouts `json:"data,omitempty"`
}

type HangoutResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message,omitempty"`
	Data    hangouts.HangoutById `json:"data,omitempty"`
}

type NotificationsResponse struct {
	Status  string                            `json:"status"`
	Message string                            `json:"message,omitempty"`
	Data    []notifications.NotificationsData `json:"data,omitempty"`
}

type UploadPhotoResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message,omitempty"`
	ImageUrl string `json:"imageUrl,omitempty"`
}
