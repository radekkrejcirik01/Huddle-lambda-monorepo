package notifications

import (
	"time"

	"gorm.io/gorm"
)

const timeFormat = "2006-01-02 15:04:05"

type Notification struct {
	Sender   string
	Receiver string
	Type     string
	Accepted *int
	Seen     int
	Created  int64
}

type Profile struct {
	Username     string
	Firstname    string
	ProfilePhoto string
}

type NotificationData struct {
	Id           uint   `json:"id"`
	Sender       string `json:"sender"`
	SenderName   string `json:"senderName"`
	ProfilePhoto string `json:"profilePhoto"`
	Type         string `json:"type"`
	Accepted     *int   `json:"accepted,omitempty"`
	Created      string `json:"created"`
}

// Get notifications from notifications_people, notifications_notify and
// notifications_huddles tables
func GetNotifications(db *gorm.DB, username string) ([]NotificationData, error) {
	var notificationsData []NotificationData

	db.Transaction(func(tx *gorm.DB) error {
		tx.Table("notifications_people").Where("receiver = ?", username).Update("seen", 1)
		tx.Table("notifications_notify").Where("receiver = ?", username).Update("seen", 1)
		tx.Table("notifications_huddles").Where("receiver = ?", username).Update("seen", 1)
		return nil
	})

	var notifications []Notification
	// Null as accepted for tables without accepted column
	query := `
			(SELECT sender, receiver, type, accepted, seen, created FROM notifications_people WHERE receiver = ?
			UNION
			SELECT sender, receiver, type, NULL AS accepted, seen, created FROM notifications_notify WHERE receiver = ?
			UNION
			SELECT sender, receiver, type, NULL AS accepted, seen, created FROM notifications_huddles WHERE receiver = ?)
			LIMIT ?
			`

	if err := db.Raw(query, username, username, username, 10).Scan(&notifications).Error; err != nil {
		return notificationsData, err
	}

	usernames := getUsernames(notifications)

	var profiles []Profile
	if err := db.Table("users").Select("username, firstname, profile_photo").Where("username IN ?", usernames).Find(&profiles).Error; err != nil {
		return notificationsData, err
	}

	for i, notification := range notifications {
		for _, profile := range profiles {
			if profile.Username == notification.Sender {
				time := time.Unix(notification.Created, 0).Format(timeFormat)

				notificationsData = append(notificationsData, NotificationData{
					Id:           uint(i),
					Sender:       notification.Sender,
					SenderName:   profile.Firstname,
					ProfilePhoto: profile.ProfilePhoto,
					Created:      time,
					Accepted:     notification.Accepted,
					Type:         notification.Type,
				})

				break
			}
		}
	}

	return notificationsData, nil
}

func getUsernames(notifications []Notification) []string {
	var usernames []string
	for _, notification := range notifications {
		usernames = append(usernames, notification.Sender)
	}

	return usernames
}
