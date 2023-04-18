package huddles

type HuddleInteracted struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Username string
	HuddleId uint
	Created  int64 `gorm:"autoCreateTime"`
}

func (HuddleInteracted) TableName() string {
	return "huddles_interacted"
}
