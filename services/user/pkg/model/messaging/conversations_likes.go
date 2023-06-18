package messaging

import (
	"gorm.io/gorm"
)

type ConversationLike struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Sender         string
	ConversationId int
}

func (ConversationLike) TableName() string {
	return "conversations_likes"
}

// LikeConversation like conversation
func LikeConversation(db *gorm.DB, t *ConversationLike) error {
	return db.Table("conversations_likes").Create(&t).Error
}

// GetConversationLike // get if conversation is liked
func GetConversationLike(db *gorm.DB, sender string, conversationId string) (int, error) {
	var like ConversationLike

	err := db.
		Table("conversations_likes").
		Where("sender = ? AND conversation_id = ?", sender, conversationId).
		Find(&like).
		Error

	if err != nil {
		return 0, err
	}

	if like == (ConversationLike{}) {
		return 0, nil
	}

	return 1, nil
}

// LikeConversation like conversation
func RemoveConversationLike(db *gorm.DB, sender string, conversationId string) error {
	return db.
		Table("conversations_likes").
		Where("sender = ? AND conversation_id = ?", sender, conversationId).
		Delete(&ConversationLike{}).
		Error
}
