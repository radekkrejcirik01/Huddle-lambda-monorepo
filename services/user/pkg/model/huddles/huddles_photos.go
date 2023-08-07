package huddles

type HuddlePhoto struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	HuddleId int
	Url      string
}

func (HuddlePhoto) TableName() string {
	return "huddles_photos"
}
