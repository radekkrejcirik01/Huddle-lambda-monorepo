package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/middleware"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/huddles"
)

// CreateHuddle POST /huddle
func CreateHuddle(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &huddles.NewHuddle{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.CreateHuddle(database.DB, username, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle successfully created",
	})
}

// GetUserHuddles GET /user-huddles/:lastId?
func GetUserHuddles(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	lastId := c.Params("lastId")

	huddleData, err := huddles.GetUserHuddles(database.DB, username, lastId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetHuddlesResponse{
		Status:  "success",
		Message: "User huddles successfully got",
		Data:    huddleData,
	})
}

// GetHuddles GET /huddles/:lastId?
func GetHuddles(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	lastId := c.Params("lastId")

	huddleData, getErr := huddles.GetHuddles(database.DB, username, lastId)

	if getErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: getErr.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetHuddlesResponse{
		Status:  "success",
		Message: "Huddles successfully got",
		Data:    huddleData,
	})
}

// UpdateHuddle PUT /huddle
func UpdateHuddle(c *fiber.Ctx) error {
	_, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

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

// GetHuddleById GET /huddle/:id
func GetHuddleById(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	id := c.Params("id")

	huddle, err := huddles.GetHuddleById(database.DB, id, username)

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
	_, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	id := c.Params("id")

	if err := huddles.DeleteHuddle(database.DB, id); err != nil {
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
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &huddles.Interact{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.HuddleInteract(database.DB, username, t); err != nil {
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

// GetHuddleInteractions GET /interactions/:huddleId
func GetHuddleInteractions(c *fiber.Ctx) error {
	_, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	huddleId := c.Params("huddleId")

	huddleInteractions, getErr := huddles.GetHuddleInteractions(database.DB, huddleId)

	if getErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: getErr.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetHuddleInteractionsResponse{
		Status:  "success",
		Message: "Huddle interactions successfully got",
		Data:    huddleInteractions,
	})
}

// RemoveHuddleInteraction DELETE /interaction/:huddleId
func RemoveHuddleInteraction(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	huddleId := c.Params("huddleId")

	if err := huddles.RemoveHuddleInteraction(database.DB, username, huddleId); err != nil {
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

// AddHuddleComment POST /comment
func AddHuddleComment(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &huddles.HuddleComment{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	t.Sender = username

	if err := huddles.AddHuddleComment(database.DB, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Comment successfully added",
	})
}

// AddHuddleMentionComment POST /comment-mention
func AddHuddleMentionComment(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &huddles.MentionComment{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.AddHuddleMentionComment(database.DB, username, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Mention comment successfully added",
	})
}

// GetHuddleComments GET /comments/:huddleId/:lastId?
func GetHuddleComments(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	huddleId := c.Params("huddleId")
	lastId := c.Params("lastId")

	comments, mentions, getErr := huddles.GetHuddleComments(database.DB, huddleId, username, lastId)
	if getErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: getErr.Error(),
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
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}

	t := &huddles.Like{}

	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := huddles.LikeHuddleComment(database.DB, username, t); err != nil {
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

// GetCommentLikes GET /comment-likes/:commentId/lastId?
func GetCommentLikes(c *fiber.Ctx) error {
	_, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	commentId := c.Params("commentId")
	lastId := c.Params("lastId")

	profiles, getErr := huddles.GetCommentLikes(database.DB, commentId, lastId)

	if getErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: getErr.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetHuddleCommentsLikesResponse{
		Status:  "success",
		Message: "Huddle comment likes successfully got",
		Data:    profiles,
	})
}

// DeleteHuddleComment DELETE /comment/:id
func DeleteHuddleComment(c *fiber.Ctx) error {
	_, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	id := c.Params("id")

	if err := huddles.DeleteHuddleComment(database.DB, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Huddle comment successfully deleted",
	})
}

// RemoveHuddleCommentLike DELETE /comment-like/:commentId
func RemoveHuddleCommentLike(c *fiber.Ctx) error {
	username, err := middleware.Authorize(c)
	if err != nil {
		return err
	}
	commentId := c.Params("commentId")

	if err := huddles.RemoveHuddleCommentLike(database.DB, commentId, username); err != nil {
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
