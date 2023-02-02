package messages

type PeopleInConversations struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	ConversationId uint
	Username       string
	Deleted        uint `gorm:"default:0"`
}

func (PeopleInConversations) TableName() string {
	return "people_in_conversations"
}
