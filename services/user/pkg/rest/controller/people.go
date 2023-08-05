package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/middleware"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
)

// AddPersonInvite POST /person
func AddPersonInvite(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &people.Invite{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	t.Sender = username

	message, err := people.AddPersonInvite(database.DB, t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(PeopleResponse{
		Status:  "success",
		Message: message,
	})
}

// GetPeople GET /people/:lastId?
func GetPeople(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	lastId := c.Params("lastId")

	peopleData, err := people.GetPeople(database.DB, username, lastId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(PeopleResponse{
		Status:  "success",
		Message: "People successfully got",
		Data:    peopleData,
	})
}

// AcceptPersonInvite PUT /person
func AcceptPersonInvite(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &people.Invite{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	t.Sender = username

	if err := people.AcceptPersonInvite(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Invite successfully accepted",
	})
}

// GetUnseenInvites GET /unseen-invites
func GetUnseenInvites(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	number, getErr := people.GetUnseenInvites(database.DB, username)

	if getErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: getErr.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetUnseenInvitesResponse{
		Status:  "success",
		Message: "Unseen invites number successfully got",
		Number:  number,
	})
}

// UpdateSeenInvites PUT /seen-invites
func UpdateSeenInvites(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	if err := people.UpdateSeenInvites(database.DB, username); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Seen invites successfully updated",
	})
}

// BlockUser POST /block-user
func BlockUser(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &people.Blocked{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if t.Blocked == username {
		return c.Status(fiber.StatusOK).JSON(Response{
			Status:  "selfblock",
			Message: "Why are you blocking yourself 😀",
		})
	}

	t.User = username

	if err := people.BlockUser(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Could not block user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "User blocked successfully",
	})
}

// MuteConversation POST /mute-conversation
func MuteConversation(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &people.MutedConversation{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	t.User = username

	if err := people.MuteConversation(database.DB, t); err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{
			Status:  "error",
			Message: "No record found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Conversation mute successfully updated",
	})
}

// IsConversationMuted GET /muted-conversation/:conversationId
func IsConversationMuted(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	conversationId := c.Params("conversationId")

	isMuted, cPeople, getErr := people.IsConversationMuted(database.DB, username, conversationId)

	if getErr != nil {
		return c.Status(fiber.StatusOK).JSON(Response{
			Status:  "error",
			Message: "No record found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetIsConversationMutedResponse{
		Status:  "success",
		Message: "Conversation mute successfully updated",
		Muted:   isMuted,
		People:  cPeople,
	})
}
