package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/rest/controller"
)

// Create new REST API serveer
func Create() *fiber.App {
	app := fiber.New()

	app.Get("/", controller.Index)
	app.Get("/notifications/:username", controller.GetNotifications)

	app.Post("/create", controller.CreateUser)
	app.Post("/get", controller.GetUser)

	app.Post("/upload/photo", controller.UploadPhoto)

	app.Post("/create/people/invitation", controller.CreatePeopleInvitation)
	app.Post("/get/people", controller.GetPeople)
	app.Post("/accept/people/invitation", controller.AcceptPeopleInvitation)
	app.Post("/check/people/invitations", controller.CheckInvitations)
	app.Post("/remove/friend", controller.RemoveFriend)

	app.Post("/notify", controller.Notify)

	return app
}
