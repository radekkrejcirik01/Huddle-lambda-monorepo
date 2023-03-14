package messages

import (
	"gorm.io/gorm"
)

type PeopleInConversations struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	ConversationId uint
	Username       string
	Deleted        uint `gorm:"default:0"`
}

func (PeopleInConversations) TableName() string {
	return "people_in_conversations"
}

// Get conversation usernames from DB
func GetConversationUsernames(db *gorm.DB, t *ConversationId) ([]string, error) {
	var usernames []string
	err := db.Table("people_in_conversations").Select("username").Where("conversation_id = ?", t.ConversationId).Find(&usernames).Error

	return usernames, err
}
