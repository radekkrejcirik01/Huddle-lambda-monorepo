package people

import (
	"gorm.io/gorm"
)

type MutedConversation struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	User           string
	ConversationId int
}

func (MutedConversation) TableName() string {
	return "muted_conversations"
}

// MuteConversation in muted_conversations table
func MuteConversation(db *gorm.DB, t *MutedConversation) error {
	if rows := db.Table("muted_conversations").FirstOrCreate(&t).RowsAffected; rows > 0 {
		return nil
	}
	return db.
		Table("muted_conversations").
		Where("user = ? AND conversation_id = ?", t.User, t.ConversationId).
		Delete(&MutedConversation{}).
		Error
}

func IsConversationMuted(db *gorm.DB, username string, conversationId string) (bool, []string, error) {
	var peopleInConversation []string

	if err := db.
		Table("people_in_conversations").
		Select("username").
		Where("conversation_id = ? AND username != ?", conversationId, username).
		Find(&peopleInConversation).
		Error; err != nil {
		return false, peopleInConversation, err
	}

	if err := db.
		Table("muted_conversations").
		Where("user = ? AND conversation_id = ?", username, conversationId).
		First(&MutedConversation{}).
		Error; err != nil {
		return false, peopleInConversation, nil
	}

	return true, peopleInConversation, nil
}
