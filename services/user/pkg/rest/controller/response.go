package controller

import (
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
