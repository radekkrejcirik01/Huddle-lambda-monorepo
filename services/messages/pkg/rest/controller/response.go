package controller

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/messages/pkg/model/messages"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseCreateConversation struct {
	Status         string `json:"status"`
	Message        string `json:"message"`
	ConversationId uint   `json:"conversationId"`
}

type ResponseConversationList struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []messages.ConversationList `json:"data"`
}

type ResponseMessages struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    []messages.Message `json:"data"`
}
