package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/messages/pkg/rest/controller"
)

// Create new REST API serveer
func Create() *fiber.App {
	app := fiber.New()

	app.Get("/", controller.Index)

	app.Post("/create/conversation", controller.CreateConversation)
	app.Post("/get/conversations/:page", controller.GetConversations)

	app.Post("/get/messages", controller.GetMessages)

	app.Post("/update/read", controller.UpdatLastRead)
	app.Post("/send/message", controller.SendMessage)

	return app
}
