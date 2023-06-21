package messaging

import (
	"gorm.io/gorm"
)

type LastReadMessage struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Username       string
	ConversationId int
	MessageId      int
}

func (LastReadMessage) TableName() string {
	return "last_read_messages"
}

// UpdateLastReadMessage in messages table
func UpdateLastReadMessage(db *gorm.DB, t *LastReadMessage) error {
	return db.
		Table("last_read_messages").
		Where("username = ? AND conversation_id = ?", t.Username, t.ConversationId).
		Update("message_id", t.MessageId).
		Error
}
