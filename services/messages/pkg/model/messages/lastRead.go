package messages

import (
	"strings"

	"gorm.io/gorm"
)

type LastReadMessage struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Username       string
	ConversationId uint
	MessageId      uint
}

func (LastReadMessage) TableName() string {
	return "last_read_messages"
}

// UpdateLastRead update message as read
func UpdateLastRead(db *gorm.DB, t *LastReadMessage) error {
	if err := db.Table("last_read_messages").Where("username = ? and conversation_id = ?", t.Username, t.ConversationId).
		Assign(LastReadMessage{MessageId: t.MessageId}).
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
		Devices:        *tokens,
		Type:           "conversationRead",
	}

	return SendNotification(&notification)
}
