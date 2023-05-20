package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/notifications"
)

// GetNotifications POST /get/notifications/:username/:lastId?
func GetNotifications(c *fiber.Ctx) error {
	username := c.Params("username")
	lastId := c.Params("lastId")

	notifications, err := notifications.GetNotifications(database.DB, username, lastId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(NotificationsResponse{
		Status:  "succes",
		Message: "Notifications succesfully got!",
		Data:    notifications,
	})
}
