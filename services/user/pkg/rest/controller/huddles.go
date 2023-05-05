package controller

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/huddles"
)

// AddHuddle POST /huddle
func AddHuddle(c *fiber.Ctx) error {
	t := &huddles.NewHuddle{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.AddHuddle(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle successfully added",
	})
}

// GetUserHuddles GET /huddles/user/:username
func GetUserHuddles(c *fiber.Ctx) error {
	username := c.Params("username")

	huddles, err := huddles.GetUserHuddles(database.DB, username)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetHuddlesResponse{
		Status:  "success",
		Message: "User huddles successfully got",
		Data:    huddles,
	})
}

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

	return c.Status(fiber.StatusOK).JSON(GetHuddlesResponse{
		Status:  "success",
		Message: "Huddles successfully got",
		Data:    huddles,
	})
}

// UpdateHuddle PUT /huddle
func UpdateHuddle(c *fiber.Ctx) error {
	t := &huddles.Update{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.UpdateHuddle(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle successfully updated",
	})
}

// PostHuddleAgain PUT /huddle
func PostHuddleAgain(c *fiber.Ctx) error {
	t := &huddles.PostAgain{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.PostHuddleAgain(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle successfully updated",
	})
}

// GetHuddleById GET /huddle/:id
func GetHuddleById(c *fiber.Ctx) error {
	id := c.Params("id")

	huddleId, parseErr := strconv.Atoi(id)
	if parseErr != nil {
		fmt.Println(parseErr)
	}

	huddle, err := huddles.GetHuddleById(database.DB, uint(huddleId))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetHuddleResponse{
		Status:  "success",
		Message: "Huddle successfully got",
		Data:    huddle,
	})
}

// DeleteHuddle DELETE /huddle/:id
func DeleteHuddle(c *fiber.Ctx) error {
	id := c.Params("id")

	huddleId, parseErr := strconv.Atoi(id)
	if parseErr != nil {
		fmt.Println(parseErr)
	}

	if err := huddles.DeleteHuddle(database.DB, uint(huddleId)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle deleted",
	})
}

// HuddleInteract POST /huddle/interaction
func HuddleInteract(c *fiber.Ctx) error {
	t := &huddles.HuddleNotification{}

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

// GetHuddleInteractions GET /huddle/interactions/:huddleId
func GetHuddleInteractions(c *fiber.Ctx) error {
	huddleId := c.Params("huddleId")

	id, parseErr := strconv.Atoi(huddleId)
	if parseErr != nil {
		fmt.Println(parseErr)
	}

	huddleInteractions, confirmedUser, getErr :=
		huddles.GetHuddleInteractions(database.DB, uint(id))

	if getErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: getErr.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetHuddleInteractionsResponse{
		Status:        "success",
		Message:       "Huddle interactions successfully got",
		Data:          huddleInteractions,
		ConfirmedUser: confirmedUser,
	})
}

// ConfirmHuddle POST /huddle/confirm
func ConfirmHuddle(c *fiber.Ctx) error {
	t := &huddles.HuddleNotification{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.ConfirmHuddle(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle confirmed successfully",
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

// RemoveHuddleConfirm PUT /huddle/confirm
func RemoveHuddleConfirm(c *fiber.Ctx) error {
	t := &huddles.RemoveConfirm{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.RemoveHuddleConfirm(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle confirm removed",
	})
}

// AddHuddleComment POST /huddle/comment
func AddHuddleComment(c *fiber.Ctx) error {
	t := &huddles.HuddleComment{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.AddHuddleComment(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Addded Huddle comment successfully",
	})
}

// AddHuddleMentionComment POST /huddle/comment/mention
func AddHuddleMentionComment(c *fiber.Ctx) error {
	t := &huddles.MentionComment{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.AddHuddleMentionComment(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Addded mention Huddle comment successfully",
	})
}

// GetHuddleComments GET /huddle/comments/:huddleId/:username
func GetHuddleComments(c *fiber.Ctx) error {
	huddleId := c.Params("huddleId")
	username := c.Params("username")

	id, err := strconv.Atoi(huddleId)
	if err != nil {
		fmt.Println(err)
	}

	comments, mentions, err := huddles.GetHuddleComments(database.DB, uint(id), username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetHuddleCommentsResponse{
		Status:   "success",
		Message:  "Huddle comments successfully got",
		Data:     comments,
		Mentions: mentions,
	})
}

// LikeHuddleComment POST /huddle/comment/like
func LikeHuddleComment(c *fiber.Ctx) error {
	t := &huddles.Like{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.LikeHuddleComment(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle comment successfully liked",
	})
}

// RemoveHuddleCommentLike DELETE /huddle/comment/like/:commentId/:sender
func RemoveHuddleCommentLike(c *fiber.Ctx) error {
	commentId := c.Params("commentId")
	sender := c.Params("sender")

	id, err := strconv.Atoi(commentId)
	if err != nil {
		fmt.Println(err)
	}

	if err := huddles.RemoveHuddleCommentLike(database.DB, id, sender); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle comment like removed",
	})
}
