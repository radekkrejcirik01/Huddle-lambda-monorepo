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
	ConversationId uint   `json:"conversationId,omitempty"`
}

type ResponseConversationList struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []messages.ConversationList `json:"data,omitempty"`
}

type ResponseConversationDetails struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    messages.ConversationDetails `json:"data,omitempty"`
}

type ResponseMessages struct {
	Status  string                     `json:"status"`
	Message string                     `json:"message"`
	Data    []messages.MessageResponse `json:"data,omitempty"`
}
