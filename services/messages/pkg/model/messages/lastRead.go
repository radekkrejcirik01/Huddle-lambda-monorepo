package messages

import "gorm.io/gorm"

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
	return db.Where("username = ? and conversation_id = ?", t.Username, t.ConversationId).
		Assign(LastReadMessage{MessageId: t.MessageId}).
		FirstOrCreate(t).Error
}
