package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
)

// CreatePeopleInvitation POST /create/people/invitation
func CreatePeopleInvitation(c *fiber.Ctx) error {
	t := &people.PeopleInvitationTable{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	message, err := people.CreatePeopleInvitation(database.DB, t)
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

// GetPeople POST /get/people
func GetPeople(c *fiber.Ctx) error {
	t := &people.People{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	people, err := people.GetPeople(database.DB, t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(PeopleResponse{
		Status:  "succes",
		Message: "People succesfully got!",
		Data:    people,
	})
}
