package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/rest/controller"
)

// Create new REST API server
func Create() *fiber.App {
	app := fiber.New()

	app.Get("/", controller.Index)
	app.Get("/user/:username", controller.GetUser)
	app.Get("/people/:username/:lastId?", controller.GetPeople)
	app.Get("/invites/:username", controller.GetInvites)
	app.Get("/hides/:username/:lastId?", controller.GetHiddenPeople)
	app.Get("/huddles/:username/:lastId?", controller.GetHuddles)
	app.Get("/user-huddles/:username/:lastId?", controller.GetUserHuddles)
	app.Get("/huddle/:id/:username", controller.GetHuddleById)
	app.Get("/interactions/:id", controller.GetHuddleInteractions)
	app.Get("/comments/:huddleId/:username/:lastId?", controller.GetHuddleComments)
	app.Get("/likes/:commentId/:lastId?", controller.GetCommentLikes)
	app.Get("/chats/:username/:lastId?", controller.GetChats)
	app.Get("/conversation/:conversationId/:lastId?", controller.GetConversation)
	app.Get("/messages/:user1/:user2", controller.GetMessagesByUsernames)

	app.Post("/user", controller.CreateUser)
	app.Post("/photo", controller.UploadPhoto)
	app.Post("/person", controller.AddPersonInvite)
	app.Post("/huddle", controller.AddHuddle)
	app.Post("/huddle/interaction", controller.HuddleInteract)
	app.Post("/huddle/comment", controller.AddHuddleComment)
	app.Post("/huddle/comment/mention", controller.AddHuddleMentionComment)
	app.Post("/huddle/comment/like", controller.LikeHuddleComment)
	app.Post("/message", controller.SendMessage)

	app.Put("/person", controller.AcceptPersonInvite)
	app.Put("/huddle", controller.UpdateHuddle)
	app.Put("/hide", controller.UpdateHiddenPeople)

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
