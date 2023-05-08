package messaging

type LastReadMessage struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Username       string
	ConversationId uint
	MessageId      uint
}

func (LastReadMessage) TableName() string {
	return "last_read_messages"
}
