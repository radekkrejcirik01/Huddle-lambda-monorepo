package notifications

import (
	"time"

	"gorm.io/gorm"
)

const timeFormat = "2006-01-02 15:04:05"

// Type can be:
// friend_invite
// friend_accepted
// hangout_notify
// huddle_interacted
// huddle_confirmed
type Notifications struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	Sender    string
	Receiver  string
	Type      string `gorm:"type:enum('friend_invite', 'friend_accepted', 'hangout_notify', 'huddle_interacted', 'huddle_confirmed')"`
	Confirmed *int
	Seen      int   `gorm:"default:0"`
	Created   int64 `gorm:"autoCreateTime"`
}

func (Notifications) TableName() string {
	return "notifications"
}

type AcceptedInvitations struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	EventId  uint
	User     string
	Username string
	Time     string
	Type     string
	Seen     int `gorm:"default:0"`
}

type Data struct {
	Id             uint   `json:"id"`
	Sender         string `json:"sender"`
	SenderName     string `json:"senderName"`
	Created        string `json:"created"`
	ProfilePicture string `json:"profilePicture"`
	Confirmed      *int   `json:"confirmed,omitempty"`
	Type           string `json:"type"`
}

type Profile struct {
	Username       string
	Firstname      string
	ProfilePicture string
}

// Get notifications from DB
func GetNotifications(db *gorm.DB, username string) ([]Data, error) {
	if err := db.Table("notifications").Where("receiver = ?", username).Update("seen", 1).Error; err != nil {
		return []Data{}, err
	}

	var notifications []Notifications
	if err := db.Table("notifications").Where("receiver = ?", username).Find(&notifications).Error; err != nil {
		return []Data{}, err
	}

	usernames := getUsernames(notifications)

	var profiles []Profile
	if err := db.Table("users").Select("username, firstname, profile_picture").Where("username IN ?", usernames).Find(&profiles).Error; err != nil {
		return []Data{}, err
	}

	var data []Data
	for _, notification := range notifications {
		for _, profile := range profiles {
			if profile.Username == notification.Sender {
				time := time.Unix(notification.Created, 0).Format(timeFormat)

				data = append(data, Data{
					Sender:         notification.Sender,
					SenderName:     profile.Firstname,
					ProfilePicture: profile.ProfilePicture,
					Created:        time,
					Confirmed:      notification.Confirmed,
					Type:           notification.Type,
				})

				break
			}
		}
	}

	return data, nil
}

func GetProfilePicturesFromQuery(db *gorm.DB, query string) ([]Profile, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var profilePictures []Profile
	for rows.Next() {
		db.ScanRows(rows, &profilePictures)
	}

	return profilePictures, nil
}

func getUsernames(notifications []Notifications) []string {
	var usersnames []string
	for _, notification := range notifications {
		usersnames = append(usersnames, notification.Sender)
	}

	return usersnames
}
