package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/hangouts"
)

// CreateHangout POST /create/hangout
func CreateHangout(c *fiber.Ctx) error {
	t := &hangouts.HangoutInvite{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := hangouts.CreateHangout(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Hangout succesfully created!",
	})
}

// CreateGroupHangout POST /create/hangout/group
func CreateGroupHangout(c *fiber.Ctx) error {
	t := &hangouts.GroupHangoutInvite{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := hangouts.CreateGroupHangout(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Group hangout succesfully created!",
	})
}

// GetHangouts POST /get/hangouts
func GetHangouts(c *fiber.Ctx) error {
	t := &hangouts.GetHangout{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	hangouts, err := hangouts.GetHangouts(database.DB, t)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(HangoutsResponse{
		Status:  "succes",
		Message: "Hangouts succesfully got!",
		Data:    hangouts,
	})
}

// GetHangout POST /get/hangout
func GetHangout(c *fiber.Ctx) error {
	t := &hangouts.HangoutId{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	hangout, err := hangouts.GetHangoutById(database.DB, t)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(HangoutResponse{
		Status:  "succes",
		Message: "Hangout succesfully got!",
		Data:    hangout,
	})
}

// AcceptHangoutInvitation POST /accept/hangout/invitation
func AcceptHangoutInvitation(c *fiber.Ctx) error {
	t := &hangouts.AcceptInvite{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := hangouts.AcceptHangout(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "Hangout invitation succesfully accepted!",
	})
}
