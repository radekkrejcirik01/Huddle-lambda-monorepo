package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Index GET /
func Index(c *fiber.Ctx) error {
	return c.JSON(struct {
		Status bool
		TS     time.Time
	}{
		Status: true,
		TS:     time.Now(),
	})
}
