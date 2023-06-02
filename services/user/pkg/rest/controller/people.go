package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
)

// AddPersonInvite POST /person
func AddPersonInvite(c *fiber.Ctx) error {
	t := &people.Invite{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	message, err := people.AddPersonInvite(database.DB, t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(PeopleResponse{
		Status:  "succes",
		Message: message,
	})
}

// GetPeople GET /people/:username/:lastId?
func GetPeople(c *fiber.Ctx) error {
	username := c.Params("username")
	lastId := c.Params("lastId")

	people, invitesNumber, err := people.GetPeople(database.DB, username, lastId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(PeopleResponse{
		Status:        "succes",
		Message:       "People succesfully got",
		Data:          people,
		InvitesNumber: invitesNumber,
	})
}

// AcceptPersonInvite PUT /people/invite
func AcceptPersonInvite(c *fiber.Ctx) error {
	t := &people.Invite{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := people.AcceptPersonInvite(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Invite succesfully accepted",
	})
}

// GetInvites GET /invites/:username
func GetInvites(c *fiber.Ctx) error {
	username := c.Params("username")

	invites, err := people.GetInvites(database.DB, username)

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{
			Status:  "error",
			Message: "No record found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetInviteResponse{
		Status:  "success",
		Message: "Invites successfully got",
		Data:    invites,
	})
}

// GetHiddenPeople GET /hides/:username/:lastId?
func GetHiddenPeople(c *fiber.Ctx) error {
	username := c.Params("username")
	lastId := c.Params("lastId")

	hiddenPeople, err := people.GetHiddenPeople(database.DB, username, lastId)

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{
			Status:  "error",
			Message: "No record found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetHiddenPeopleResponse{
		Status:  "succes",
		Message: "Hidden people successfully got",
		Data:    hiddenPeople,
	})
}

// UpdateHiddenPeople PUT /hide
func UpdateHiddenPeople(c *fiber.Ctx) error {
	t := &people.HidePeople{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := people.UpdateHiddenPeople(database.DB, t); err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{
			Status:  "error",
			Message: "No record found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Hidden people successfully updated",
	})
}

// MuteHuddles POST /mute-huddles
func MuteHuddles(c *fiber.Ctx) error {
	t := &people.MutedHuddle{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := people.MuteHuddles(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Could not mute huddles",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddles muted successfully",
	})
}

// GetMutedHuddles GET /muted-huddles/:username
func GetMutedHuddles(c *fiber.Ctx) error {
	username := c.Params("username")

	people, err := people.GetMutedHuddles(database.DB, username)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Could not get muted huddles",
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetMutedHuddlesResponse{
		Status:  "success",
		Message: "Muted Huddles successfully got",
		Data:    people,
	})
}

// RemoveMutedHuddles DELETE /muted-huddles/:user1/:user2
func RemoveMutedHuddles(c *fiber.Ctx) error {
	user1 := c.Params("user1")
	user2 := c.Params("user2")

	if err := people.RemoveMutedHuddles(database.DB, user1, user2); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Could not unmute huddles",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddles muted successfully",
	})
}

// MuteConversation POST /mute-conversation
func MuteConversation(c *fiber.Ctx) error {
	t := &people.MutedConversation{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

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

// IsConversationMuted GET /muted/:username/:conversationId
func IsConversationMuted(c *fiber.Ctx) error {
	username := c.Params("username")
	conversationId := c.Params("conversationId")

	isMuted, err := people.IsConversationMuted(database.DB, username, conversationId)

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{
			Status:  "error",
			Message: "No record found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetIsConversationMutedResponse{
		Status:  "success",
		Message: "Conversation mute successfully updated",
		Muted:   isMuted,
	})
}

// RemovePerson DELETE /person
func RemovePerson(c *fiber.Ctx) error {
	user1 := c.Params("user1")
	user2 := c.Params("user2")

	if err := people.RemovePerson(database.DB, user1, user2); err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Connection removed",
	})
}
