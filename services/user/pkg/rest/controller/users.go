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

// GetPeopleNumber GET /people-number/:username
func GetPeopleNumber(c *fiber.Ctx) error {
	username := c.Params("username")

	number, err := users.GetPeopleNumber(database.DB, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(PeopleNumberResponse{
		Status:       "succes",
		Message:      "People number succesfully got",
		PeopleNumber: number,
	})
}

// UploadPhoto POST /upload/photo
func UploadPhoto(c *fiber.Ctx) error {
	t := &users.UplaodProfilePhotoBody{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	imageUrl, err := users.UplaodProfilePhoto(database.DB, t)
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
