package controller

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/messages/pkg/model/matches"
	"github.com/radekkrejcirik01/PingMe-backend/services/messages/pkg/model/messages"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseMatches struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    []matches.Matched `json:"data"`
}

type ResponseConversationList struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []messages.ConversationList `json:"data"`
}

type ResponseMessages struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    []messages.Messages `json:"data"`
}

type ResponseUser struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    messages.ChatUser `json:"data"`
}
