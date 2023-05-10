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
	app.Get("/huddles/:username", controller.GetHuddles)
	app.Get("/user-huddles/:username", controller.GetUserHuddles)
	app.Get("/huddle/:id/:username", controller.GetHuddleById)
	app.Get("/interactions/:id", controller.GetHuddleInteractions)
	app.Get("/comments/:huddleId/:username", controller.GetHuddleComments)
	app.Get("/chats/:username", controller.GetChats)
	app.Get("/conversation/:conversationId", controller.GetConversation)
	app.Get("/messages/:user1/:user2", controller.GetMessagesByUsernames)

	app.Post("/user", controller.CreateUser)
	app.Post("/photo", controller.UploadPhoto)
	app.Post("/person", controller.AddPersonInvite)
	app.Post("/huddle", controller.AddHuddle)
	app.Post("/huddle/interaction", controller.HuddleInteract)
	app.Post("/huddle/confirm", controller.ConfirmHuddle)
	app.Post("/huddle/comment", controller.AddHuddleComment)
	app.Post("/huddle/comment/mention", controller.AddHuddleMentionComment)
	app.Post("/huddle/comment/like", controller.LikeHuddleComment)
	app.Post("/message", controller.SendMessage)

	app.Put("/person", controller.AcceptPersonInvite)
	app.Put("/huddle", controller.UpdateHuddle)
	app.Put("/huddle/post", controller.PostHuddleAgain)
	app.Put("/huddle/confirm",
		controller.RemoveHuddleConfirm,
	)

	app.Delete("/person/:user1/:user2", controller.RemovePerson)
	app.Delete("/huddle/:id", controller.DeleteHuddle)
	app.Delete("/interaction/:id/:username",
		controller.RemoveHuddleInteraction,
	)
	app.Delete("/like/:id/:sender",
		controller.RemoveHuddleCommentLike,
	)

	return app
}
