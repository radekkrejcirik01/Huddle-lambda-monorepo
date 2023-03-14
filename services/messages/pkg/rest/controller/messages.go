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

	details, err := messages.CreateConversation(database.DB, t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponseConversationDetails{
		Status:  "succes",
		Message: "Conversation succesfully created",
		Data:    details,
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

// GetConversationDetails POST /get/conversation/details
func GetConversationDetails(c *fiber.Ctx) error {
	t := &messages.GetConversation{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	details, err := messages.GetConversationDetails(database.DB, t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponseConversationDetails{
		Status:  "succes",
		Message: "Conversation details succesfully get!",
		Data:    details,
	})
}

// GetConversationUsernames POST /get/conversation/usernames
func GetConversationUsernames(c *fiber.Ctx) error {
	t := &messages.ConversationId{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	users, err := messages.GetConversationUsernames(database.DB, t)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponseGetConversationUsernames{
		Status:  "succes",
		Message: "Conversation usernames succesfully got",
		Data:    users,
	})
}

// UpdateConversation POST /update/conversation
func UpdateConversation(c *fiber.Ctx) error {
	t := &messages.UpdateConversation{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messages.UpdateConversationById(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Conversation succesfully updated!",
	})
}

// AddConversationUsers POST /add/conversation/users
func AddConversationUsers(c *fiber.Ctx) error {
	t := &messages.Add{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messages.AddConversationUsers(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "User succesfully added to conversation!",
	})
}

// RemoveConversation POST /remove/conversation
func RemoveConversation(c *fiber.Ctx) error {
	t := &messages.Remove{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messages.RemoveConversation(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Conversation succesfully removed",
	})
}

// RemoveUserFromConversation POST /remove/conversation/user
func RemoveUserFromConversation(c *fiber.Ctx) error {
	t := &messages.Remove{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messages.RemoveUserFromConversation(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "User succesfully removed from conversation",
	})
}

// DeleteConversation POST /delete/conversation
func DeleteConversation(c *fiber.Ctx) error {
	t := &messages.Delete{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messages.DeleteConversation(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Conversation succesfully deleted",
	})
}

// GetMessages POST /get/messages
func GetMessages(c *fiber.Ctx) error {
	t := &messages.ConversationId{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	messages, err := messages.GetMessages(database.DB, t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponseMessages{
		Status:  "succes",
		Message: "Messages succesfully got!",
		Data:    messages,
	})
}

// UpdateRead POST /update/read
func UpdatLastRead(c *fiber.Ctx) error {
	t := &messages.LastReadMessage{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messages.UpdateLastRead(database.DB, t); err != nil {
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

// MessageReacted POST /react/message
func MessageReacted(c *fiber.Ctx) error {
	t := &messages.Reacted{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messages.MessageReacted(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Succesfully reacted on the message",
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

// SendTyping POST /send/typing
func SendTyping(c *fiber.Ctx) error {
	t := &messages.Typing{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messages.SendTyping(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Typing succesfully sent",
	})
}
