package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/rest/controller"
)

// Create new REST API serveer
func Create() *fiber.App {
	app := fiber.New()

	app.Get("/", controller.Index)

	app.Post("/create", controller.CreateUser)
	app.Post("/get", controller.GetUser)

	app.Post("/upload/photo", controller.UploadPhoto)

	app.Post("/create/people/invitation", controller.CreatePeopleInvitation)
	app.Post("/get/people", controller.GetPeople)
	app.Post("/accept/people/invitation", controller.AcceptPeopleInvitation)

	app.Post("/create/hangout/group", controller.CreateGroupHangout)
	app.Post("/create/hangout", controller.CreateHangout)
	app.Post("/get/hangouts", controller.GetHangouts)
	app.Post("/get/hangout", controller.GetHangout)
	app.Post("/accept/hangout/invitation", controller.AcceptHangoutInvitation)
	app.Post("/send/hangout/invitation", controller.SendHangoutInvitation)

	app.Post("/get/notifications", controller.GetNotifications)

	return app
}
