package messaging

import (
	"gorm.io/gorm"
)

type LastSeenMessage struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Username       string
	ConversationId int
	MessageId      int
}

func (LastSeenMessage) TableName() string {
	return "last_seen_messages"
}

type LastSeen struct {
	ConversationId int
}

// UpdateLastSeen in last_seen_messages table
func UpdateLastSeen(db *gorm.DB, t *LastSeen, username string) error {
	var lastMessageId int
	var personInConversation string
	var lastHuddleId int

	if err := db.
		Table("messages").
		Select("id").
		Where("conversation_id = ?", t.ConversationId).
		Find(&lastMessageId).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("people_in_conversations").
		Select("username").
		Where("conversation_id = ? AND username != ?", t.ConversationId, username).
		Find(&personInConversation).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("huddles").
		Select("id").
		Where("created_by = ?", personInConversation).
		Find(&lastHuddleId).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("last_seen_messages").
		Where("username = ? AND conversation_id = ?", username, t.ConversationId).
		Update("message_id", lastMessageId).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("last_seen_huddles").
		Where("username = ? AND conversation_id = ?", username, t.ConversationId).
		Update("huddle_id", lastHuddleId).
		Error; err != nil {
		return err
	}

	return nil
}
