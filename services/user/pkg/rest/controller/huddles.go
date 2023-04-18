package controller

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/huddles"
)

// GetHuddles GET /huddles/:username
func GetHuddles(c *fiber.Ctx) error {
	username := c.Params("username")

	huddles, err := huddles.GetHuddles(database.DB, username)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(HuddlesResponse{
		Status:  "success",
		Message: "Huddles successfully got",
		Data:    huddles,
	})
}

// HuddleInteract POST /huddle/interaction
func HuddleInteract(c *fiber.Ctx) error {
	t := &huddles.HuddleInteracted{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.HuddleInteract(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle successfully interacted",
	})
}

// RemoveHuddleInteraction DELETE /huddle/interaction
func RemoveHuddleInteraction(c *fiber.Ctx) error {
	username := c.Params("username")
	huddleId := c.Params("huddleId")

	id, err := strconv.Atoi(huddleId)
	if err != nil {
		fmt.Println(err)
	}

	if err := huddles.RemoveHuddleInteraction(database.DB, username, uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle interaction removed",
	})
}
