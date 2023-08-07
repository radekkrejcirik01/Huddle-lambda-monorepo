package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/rest/controller"
)

// Create new REST API server
func Create() *fiber.App {
	app := fiber.New()

	app.Get("/user", controller.GetUser)
	app.Get("/people/:lastId?", controller.GetPeople)
	app.Get("/notifications", controller.GetUserNotifications)
	app.Get("/unseen-invites", controller.GetUnseenInvites)
	app.Get("/user-huddles/:lastId?", controller.GetUserHuddles)
	app.Get("/huddle/:id", controller.GetHuddle)
	app.Get("/huddle-likes/:huddleId", controller.GetHuddleLikes)
	app.Get("/comments/:huddleId/:lastId?", controller.GetHuddleComments)
	app.Get("/comment-likes/:commentId/:lastId?", controller.GetCommentLikes)
	app.Get("/chats/:lastId?", controller.GetChats)
	app.Get("/conversation/:conversationId/:lastId?", controller.GetConversation)
	app.Get("/messages/:user", controller.GetMessagesByUsernames)
	app.Get("/muted-conversation/:conversationId", controller.IsConversationMuted)
	app.Get("/conversation-like/:conversationId", controller.GetConversationLike)

	app.Post("/user", controller.CreateUser)
	app.Post("/login", controller.LoginUser)
	app.Post("/photo", controller.UploadPhoto)
	app.Post("/person", controller.AddPersonInvite)
	app.Post("/huddle", controller.CreateHuddle)
	app.Post("/huddle-photo", controller.UploadHuddlePhoto)
	app.Post("/huddle/like", controller.LikeHuddle)
	app.Post("/comment", controller.AddHuddleComment)
	app.Post("/comment-mention", controller.AddHuddleMentionComment)
	app.Post("/comment-like", controller.LikeHuddleComment)
	app.Post("/message", controller.SendMessage)
	app.Post("/conversation-like", controller.LikeConversation)
	app.Post("/mute-conversation", controller.MuteConversation)
	app.Post("/device", controller.SaveDevice)
	app.Post("/message-react", controller.MessageReact)
	app.Post("/block-user", controller.BlockUser)

	app.Put("/person", controller.AcceptPersonInvite)
	app.Put("/notification", controller.UpdateUserNotification)
	app.Put("/seen-invites", controller.UpdateSeenInvites)
	app.Put("/last-seen", controller.UpdateLastSeen)

	app.Delete("/huddle/:id", controller.DeleteHuddle)
	app.Delete("/comment/:id", controller.DeleteHuddleComment)
	app.Delete("/like/:huddleId", controller.RemoveHuddleLike)
	app.Delete("/comment-like/:commentId", controller.RemoveHuddleCommentLike)
	app.Delete("/device", controller.DeleteDevice)
	app.Delete("/conversation-like/:conversationId", controller.RemoveConversationLike)
	app.Delete("/account", controller.DeleteAccount)

	return app
}
