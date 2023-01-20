package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/messages/pkg/database"
	messages "github.com/radekkrejcirik01/PingMe-backend/services/messages/pkg/model/messages"
)

// CreateConversation POST /get/conversations/:page
func CreateConversation(c *fiber.Ctx) error {
	t := &messages.ConversationCreate{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	conversationId, err := messages.CreateConversation(database.DB, t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponseCreateConversation{
		Status:         "succes",
		Message:        "Conversation succesfully created",
		ConversationId: conversationId,
	})
}

// GetConversations POST /get/conversations/:page
func GetConversations(c *fiber.Ctx) error {
	page := c.Params("page")

	t := &messages.Username{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	conversationList, err := messages.GetConversationsList(database.DB, t, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponseConversationList{
		Status:  "succes",
		Message: "Conversation list succesfully get",
		Data:    conversationList,
	})
}

// UpdateRead POST /update/read
func UpdateRead(c *fiber.Ctx) error {
	t := &messages.MessagesBody{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messages.UpdateRead(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Read succesfully updated",
	})
}

// SendMessage POST /send/messages
func SendMessage(c *fiber.Ctx) error {
	t := &messages.SentMessage{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messages.SendMessage(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Message succesfully sent",
	})
}
