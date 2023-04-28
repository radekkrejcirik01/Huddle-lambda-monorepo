package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/rest/controller"
)

// Create new REST API serveer
func Create() *fiber.App {
	app := fiber.New()

	app.Get("/", controller.Index)
	app.Get("/user/:username", controller.GetUser)
	app.Get("/notifications/:username", controller.GetNotifications)
	app.Get("/people/:username", controller.GetPeople)
	app.Get("/person/:user1/:user2", controller.GetPersonInvite)
	app.Get("/huddles/user/:username", controller.GetUserHuddles)
	app.Get("/huddles/:username", controller.GetHuddles)
	app.Get("/huddle/:id", controller.GetHuddleById)
	app.Get("/huddle/interactions/:huddleId",
		controller.GetHuddleInteractions,
	)

	app.Post("/user", controller.CreateUser)
	app.Post("/photo", controller.UploadPhoto)
	app.Post("/person", controller.AddPersonInvite)
	app.Post("/notify", controller.SendNotify)
	app.Post("/huddle", controller.AddHuddle)
	app.Post("/huddle/interaction", controller.HuddleInteract)
	app.Post("/huddle/confirm", controller.ConfirmHuddle)

	app.Put("/person", controller.AcceptPersonInvite)
	app.Put("/huddle", controller.UpdateHuddle)
	app.Put("/huddle/post", controller.PostHuddleAgain)

	app.Delete("/person/:user1/:user2", controller.RemovePerson)
	app.Delete("/huddle/interaction/:username/:huddleId",
		controller.RemoveHuddleInteraction,
	)

	return app
}
