package huddles

type LastSeenHuddle struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Username       string
	ConversationId int
	HuddleId       int
}

func (LastSeenHuddle) TableName() string {
	return "last_seen_huddles"
}
