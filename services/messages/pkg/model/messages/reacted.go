package messages

import (
	"strings"

	"gorm.io/gorm"
)

type Reacted struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Username       string
	ConversationId uint
	MessageId      uint
	Reaction       string
}

func (Reacted) TableName() string {
	return "reacted"
}

// MessageReacted reacts on the message
func MessageReacted(db *gorm.DB, t *Reacted) error {
	if err := db.Table("reacted").Where("username = ? and conversation_id = ? and message_id = ? and reaction = ?", t.Username, t.ConversationId, t.MessageId, t.Reaction).
		FirstOrCreate(t).Error; err != nil {
		return err
	}

	var usernames []string
	if err := db.Table("people_in_conversations").Select("username").Where("conversation_id = ? AND username != ?", t.ConversationId, t.Username).Find(&usernames).Error; err != nil {
		return err
	}

	var usernamesArray []string
	for _, username := range usernames {
		usernamesArray = append(usernamesArray, `'`+username+`'`)
	}

	usernamesString := strings.Join(usernamesArray, ", ")

	tokens := &[]string{}
	if err := GetTokensByUsernames(db, tokens, usernamesString); err != nil {
		return nil
	}
	notification := Notification{
		ConversationId: t.ConversationId,
		Sender:         t.Username,
		Body:           t.Username + " reacted with " + t.Reaction,
		Devices:        *tokens,
		Type:           "messageReacted",
	}

	return SendNotification(&notification)
}
