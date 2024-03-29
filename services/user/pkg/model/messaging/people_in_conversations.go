package messaging

type PersonInConversation struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	ConversationId int
	Username       string
}

func (PersonInConversation) TableName() string {
	return "people_in_conversations"
}
