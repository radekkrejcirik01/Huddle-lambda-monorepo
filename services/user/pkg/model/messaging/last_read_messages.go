package messaging

import (
	"gorm.io/gorm"
)

type LastReadMessage struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Username       string
	ConversationId int
	MessageId      int
	Seen           int
}

func (LastReadMessage) TableName() string {
	return "last_read_messages"
}

// UpdateLastReadMessage in last_read_messages table
func UpdateLastReadMessage(db *gorm.DB, t *LastReadMessage) error {
	if err := db.
		Table("last_read_messages").
		Where("username = ? AND conversation_id = ?", t.Username, t.ConversationId).
		Update("message_id", t.MessageId).
		Error; err != nil {
		return err
	}

	return db.
		Table("last_read_messages").
		Where("username = ? AND conversation_id = ?", t.Username, t.ConversationId).
		Update("seen", 1).
		Error
}

// UpdateLastSeenReadMessage in last_read_messages table
func UpdateLastSeenReadMessage(db *gorm.DB, username string) error {
	return db.
		Table("last_read_messages").
		Where("username = ? AND seen != 1", username).
		Update("seen", 1).
		Error
}
