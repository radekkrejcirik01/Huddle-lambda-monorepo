package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/middleware"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/messaging"
)

// GetChats GET /chats/:lastId?
func GetChats(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
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

// GetUnreadMessagesNumber GET /unread-message
func GetUnreadMessagesNumber(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	unread, err := messaging.GetUnreadMessagesNumber(database.DB, username)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetUnreadMessagesNumberResponse{
		Status:  "success",
		Message: "Unread messages number successfully got",
		Unread:  unread,
	})
}

// SendMessage POST /message
func SendMessage(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &messaging.Send{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messaging.SendMessage(database.DB, username, t); err != nil {
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

// GetConversation GET /conversation/:conversationId/:lastId?
func GetConversation(c *fiber.Ctx) error {
	_, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	conversationId := c.Params("conversationId")
	lastId := c.Params("lastId")

	messages, getErr := messaging.GetConversation(database.DB, conversationId, lastId)

	if getErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: getErr.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetMessagesResponse{
		Status:  "success",
		Message: "Messages successfully got",
		Data:    messages,
	})
}

// GetMessagesByUsernames GET /messages/:user
func GetMessagesByUsernames(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	user := c.Params("user")

	messages, conversationId, getErr := messaging.GetMessagesByUsernames(
		database.DB,
		username,
		user,
	)

	if getErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: getErr.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetMessagesByUsernamesResponse{
		Status:         "success",
		Message:        "Messages successfully got",
		Data:           messages,
		ConversationId: conversationId,
	})
}

// LikeConversation POST /conversation-like
func LikeConversation(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &messaging.ConversationLike{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	t.Sender = username

	if err := messaging.LikeConversation(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Conversation successfully likes",
	})
}

// GetConversationLike GET /conversation-like/:conversationId
func GetConversationLike(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	conversationId := c.Params("conversationId")

	isLiked, getErr := messaging.GetConversationLike(database.DB, username, conversationId)

	if getErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: getErr.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetIsConversationLikedResponse{
		Status:  "success",
		Message: "Got conversation like",
		IsLiked: isLiked,
	})
}

// RemoveConversationLike DELETE /conversation-like/:conversationId
func RemoveConversationLike(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	conversationId := c.Params("conversationId")

	if err := messaging.RemoveConversationLike(database.DB, username, conversationId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Conversation like removed",
	})
}

// MessageReact POST /message-react
func MessageReact(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &messaging.SendReaction{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := messaging.MessageReact(database.DB, username, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Successfully reacted to message",
	})
}

// UpdateLastReadMessage POST /last-read-message
func UpdateLastReadMessage(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &messaging.LastReadMessage{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	t.Username = username

	if err := messaging.UpdateLastReadMessage(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Last read message successfully updated",
	})
}

// UpdateLastSeenReadMessage PUT /last-seen-read-message
func UpdateLastSeenReadMessage(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	if err := messaging.UpdateLastSeenReadMessage(database.DB, username); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Last seen read message successfully updated",
	})
}
