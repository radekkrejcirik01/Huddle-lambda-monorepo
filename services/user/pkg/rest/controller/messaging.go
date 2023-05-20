package controller

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/messaging"
)

// GetChats GET /chats/:username/:lastId?
func GetChats(c *fiber.Ctx) error {
	username := c.Params("username")
	lastId := c.Params("lastId")

	chats, err := messaging.GetChats(database.DB, username, lastId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetChatsResponse{
		Status:  "success",
		Message: "Chats successfully got",
		Data:    chats,
	})
}

// SendMessage POST /message
func SendMessage(c *fiber.Ctx) error {
	t := &messaging.Send{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messaging.SendMessage(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Message successfully sent",
	})
}

// GetMessages GET /conversation/:conversationId/:lastId?
func GetConversation(c *fiber.Ctx) error {
	id, parseErr := strconv.Atoi(c.Params("conversationId"))
	if parseErr != nil {
		fmt.Println(parseErr)
	}
	lastId := c.Params("lastId")

	messages, err := messaging.GetConversation(database.DB, id, lastId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetMessagesResponse{
		Status:  "success",
		Message: "Messages successfully got",
		Data:    messages,
	})
}

// GetMessagesByUsernames GET /messages/:user1/:user2
func GetMessagesByUsernames(c *fiber.Ctx) error {
	user1 := c.Params("user1")
	user2 := c.Params("user2")

	messages, conversationId, err := messaging.GetMessagesByUsernames(
		database.DB,
		user1,
		user2,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetMessagesByUsernamesResponse{
		Status:         "success",
		Message:        "Messages successfully got",
		Data:           messages,
		ConversationId: conversationId,
	})
}
