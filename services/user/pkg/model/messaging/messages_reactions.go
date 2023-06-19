package messaging

import (
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

type MessageReaction struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Sender         string
	ConversationId int
	MessageId      int
	Value          string
}

func (MessageReaction) TableName() string {
	return "messages_reactions"
}

type SendReaction struct {
	Sender         string
	Receiver       string
	Message        string
	ConversationId int
	MessageId      int
	Value          string
}

func MessageReact(db *gorm.DB, t *SendReaction) error {
	interaction := MessageReaction{
		Sender:         t.Sender,
		ConversationId: t.ConversationId,
		MessageId:      t.MessageId,
		Value:          t.Value,
	}
	if err := db.Table("messages_reactions").Create(&interaction).Error; err != nil {
		return err
	}

	if t.Sender == t.Receiver {
		return nil
	}

	var mutedConversation []string
	if err := db.
		Table("muted_conversations").
		Select("user").
		Where("user = ? AND conversation_id = ?", t.Receiver, t.ConversationId).
		Find(&mutedConversation).
		Error; err != nil {
		return err
	}

	if len(mutedConversation) > 0 {
		return nil
	}

	var info []Info
	if err := db.
		Table("users").
		Select("username, firstname, profile_photo, messages_notifications").
		Where("username IN ?", []string{t.Receiver, t.Sender}).
		Find(&info).
		Error; err != nil {
		return err
	}

	if !receiveNotificationsEnabled(info, t.Receiver) {
		return nil
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Receiver); err != nil {
		return nil
	}

	body := t.Message

	senderInfo := getSenderInfo(info, t.Sender)

	fcmNotification := service.FcmNotification{
		Data: map[string]interface{}{
			"type":           "message",
			"conversationId": t.ConversationId,
			"name":           senderInfo.Firstname,
			"profilePhoto":   senderInfo.ProfilePhoto,
		},
		Title:   senderInfo.Firstname + " reacted " + t.Value,
		Body:    body,
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&fcmNotification)
}
