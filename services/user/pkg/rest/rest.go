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
	app.Get("/invites/:lastId?", controller.GetInvites)
	app.Get("/unseen-invites", controller.GetUnseenInvites)
	app.Get("/hidden-people/:lastId?", controller.GetHiddenPeople)
	app.Get("/huddles/:lastId?", controller.GetHuddles)
	app.Get("/user-huddles/:lastId?", controller.GetUserHuddles)
	app.Get("/huddle/:id", controller.GetHuddleById)
	app.Get("/interactions/:huddleId", controller.GetHuddleInteractions)
	app.Get("/comments/:huddleId/:lastId?", controller.GetHuddleComments)
	app.Get("/comment-likes/:commentId/:lastId?", controller.GetCommentLikes)
	app.Get("/chats/:lastId?", controller.GetChats)
	app.Get("/conversation/:conversationId/:lastId?", controller.GetConversation)
	app.Get("/messages/:user", controller.GetMessagesByUsernames)
	app.Get("/muted-conversation/:conversationId", controller.IsConversationMuted)
	app.Get("/muted-huddles", controller.GetMutedHuddles)
	app.Get("/conversation-like/:conversationId", controller.GetConversationLike)
	app.Get("/unread-messages", controller.GetUnreadMessagesNumber)

	app.Post("/user", controller.CreateUser)
	app.Post("/login", controller.LoginUser)
	app.Post("/photo", controller.UploadPhoto)
	app.Post("/person", controller.AddPersonInvite)
	app.Post("/huddle", controller.CreateHuddle)
	app.Post("/huddle/interaction", controller.HuddleInteract)
	app.Post("/comment", controller.AddHuddleComment)
	app.Post("/comment-mention", controller.AddHuddleMentionComment)
	app.Post("/comment-like", controller.LikeHuddleComment)
	app.Post("/message", controller.SendMessage)
	app.Post("/conversation-like", controller.LikeConversation)
	app.Post("/mute-conversation", controller.MuteConversation)
	app.Post("/mute-huddles", controller.MuteHuddles)
	app.Post("/device", controller.SaveDevice)
	app.Post("/message-react", controller.MessageReact)

	app.Put("/person", controller.AcceptPersonInvite)
	app.Put("/huddle", controller.UpdateHuddle)
	app.Put("/hide", controller.UpdateHiddenPeople)
	app.Put("/notification", controller.UpdateUserNotification)
	app.Put("/seen-invites", controller.UpdateSeenInvites)
	app.Put("/last-read-message", controller.UpdateLastReadMessage)
	app.Put("/last-seen-read-message", controller.UpdateLastSeenReadMessage)

	app.Delete("/huddle/:id", controller.DeleteHuddle)
	app.Delete("/comment/:id", controller.DeleteHuddleComment)
	app.Delete("/interaction/:huddleId", controller.RemoveHuddleInteraction)
	app.Delete("/comment-like/:commentId", controller.RemoveHuddleCommentLike)
	app.Delete("/muted-huddles/:user", controller.RemoveMutedHuddles)
	app.Delete("/device", controller.DeleteDevice)
	app.Delete("/conversation-like/:conversationId", controller.RemoveConversationLike)

	return app
}
