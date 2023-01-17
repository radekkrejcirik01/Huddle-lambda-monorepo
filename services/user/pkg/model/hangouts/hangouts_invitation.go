package hangouts

type HangoutsInvitationTable struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	HangoutId uint
	User      string
	Username  string
	Time      string
	Confirmed uint
}

func (HangoutsInvitationTable) TableName() string {
	return "hangouts_invitations"
}
