package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/users"
)

// CreateUser POST /create
func CreateUser(c *fiber.Ctx) error {
	t := &users.User{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := users.CreateUser(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "succes",
		Message: "User succesfully created!",
	})
}

// GetUser GET /user/:username
func GetUser(c *fiber.Ctx) error {
	username := c.Params("username")

	user, err := users.GetUser(database.DB, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(UserResponse{
		Status:  "succes",
		Message: "User succesfully got",
		Data:    user,
	})
}

// GetUserNotifications GET /notifications/:username
func GetUserNotifications(c *fiber.Ctx) error {
	username := c.Params("username")

	notifications, err := users.GetUserNotifications(database.DB, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(UserNotificationsResponse{
		Status:  "success",
		Message: "User notifications successfully got",
		Data:    notifications,
	})
}

// UpdateUserNotification PUT /notification
func UpdateUserNotification(c *fiber.Ctx) error {
	t := &users.UpdateNotification{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := users.UpdateUserNotification(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "User successfully got",
	})
}

// UploadPhoto POST /upload/photo
func UploadPhoto(c *fiber.Ctx) error {
	t := &users.UploadProfilePhotoBody{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	imageUrl, err := users.UploadProfilePhoto(database.DB, t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(UploadPhotoResponse{
		Status:   "succes",
		Message:  "Photo succesfully uploaded!",
		ImageUrl: imageUrl,
	})
}
