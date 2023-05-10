package huddles

type HuddleNotification struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	HuddleId uint
	Sender   string
	Receiver string
	Type     string `gorm:"type:enum('huddle_interacted', 'huddle_confirmed', 'huddle_commented', 'huddle_mention_commented', 'comment_liked')"`
	Seen     int    `gorm:"default:0"`
	Created  int64  `gorm:"autoCreateTime"`
}

func (HuddleNotification) TableName() string {
	return "notifications_huddles"
}
