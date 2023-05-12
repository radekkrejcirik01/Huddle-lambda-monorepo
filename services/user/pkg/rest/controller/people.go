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

// GetPeople GET /people
func GetPeople(c *fiber.Ctx) error {
	username := c.Params("username")

	people, err := people.GetPeople(database.DB, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(PeopleResponse{
		Status:  "succes",
		Message: "People succesfully got",
		Data:    people,
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

// GetPersonInvite GET /person/invite/:user1/:user2
func GetPersonInvite(c *fiber.Ctx) error {
	user1 := c.Params("user1")
	user2 := c.Params("user2")

	invite, err := people.GetPersonInvite(database.DB, user1, user2)

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{
			Status:  "error",
			Message: "No record found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetInviteResponse{
		Status:  "succes",
		Message: "Invite succesfully got",
		Data:    invite,
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
