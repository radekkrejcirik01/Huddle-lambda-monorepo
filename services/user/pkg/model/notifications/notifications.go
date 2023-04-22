package notifications

import (
	"time"

	"gorm.io/gorm"
)

const timeFormat = "2006-01-02 15:04:05"

const huddleInteractedType = "huddle_interacted"
const huddleConfirmedType = "huddle_confirmed"

type Notification struct {
	HuddleId *uint
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
	Id           uint    `json:"id"`
	HuddleId     *uint   `json:"huddleId,omitempty"`
	Sender       string  `json:"sender"`
	SenderName   string  `json:"senderName"`
	ProfilePhoto string  `json:"profilePhoto"`
	Type         string  `json:"type"`
	What         *string `json:"what,omitempty"`
	Accepted     *int    `json:"accepted,omitempty"`
	Confirmed    *int    `json:"confirmed,omitempty"`
	Created      string  `json:"created"`
}

type GetHuddles struct {
	Id   uint
	What string
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
			(SELECT NULL AS huddle_id, sender, receiver, type, accepted, seen, created FROM notifications_people WHERE receiver = ?
			UNION
			SELECT NULL AS huddle_id, sender, receiver, type, NULL AS accepted, seen, created FROM notifications_notify WHERE receiver = ?
			UNION
			SELECT huddle_id, sender, receiver, type, NULL AS accepted, seen, created FROM notifications_huddles WHERE receiver = ?)
			LIMIT ?
			`

	if err := db.Raw(query, username, username, username, 10).Scan(&notifications).Error; err != nil {
		return notificationsData, err
	}

	huddleIds := getHuddleIds(notifications)

	var whats []GetHuddles
	var confirmedHuddles []uint
	if len(huddleIds) > 0 {
		if err := db.Table("huddles").Select("id, what").Where("id IN ?", huddleIds).Find(&whats).Error; err != nil {
			return notificationsData, err
		}

		if err := db.
			Table("huddles_interacted").
			Select("huddle_id").
			Where("huddle_id IN ? AND confirmed = 1", huddleIds).
			Find(&confirmedHuddles).
			Error; err != nil {
			return notificationsData, err
		}
	}

	usernames := getUsernames(notifications)

	var profiles []Profile
	if err := db.
		Table("users").
		Select("username, firstname, profile_photo").
		Where("username IN ?", usernames).
		Find(&profiles).
		Error; err != nil {
		return notificationsData, err
	}

	for i, notification := range notifications {
		for _, profile := range profiles {
			if profile.Username == notification.Sender {
				time := time.Unix(notification.Created, 0).Format(timeFormat)

				var what *string
				var confirmed *int
				if notification.Type == huddleInteractedType {
					what = getWhat(whats, *notification.HuddleId)
					confirmed = getConfirmed(confirmedHuddles, *notification.HuddleId)
				}
				if notification.Type == huddleConfirmedType {
					what = getWhat(whats, *notification.HuddleId)
				}

				notificationsData = append(notificationsData, NotificationData{
					Id:           uint(i),
					HuddleId:     notification.HuddleId,
					Sender:       notification.Sender,
					SenderName:   profile.Firstname,
					ProfilePhoto: profile.ProfilePhoto,
					Type:         notification.Type,
					What:         what,
					Accepted:     notification.Accepted,
					Confirmed:    confirmed,
					Created:      time,
				})

				break
			}
		}
	}

	return notificationsData, nil
}

func getHuddleIds(notifications []Notification) []*uint {
	var huddleIds []*uint
	for _, notification := range notifications {
		huddleIds = append(huddleIds, notification.HuddleId)
	}

	return huddleIds
}

func getUsernames(notifications []Notification) []string {
	var usernames []string
	for _, notification := range notifications {
		usernames = append(usernames, notification.Sender)
	}

	return usernames
}

func getWhat(huddles []GetHuddles, huddleId uint) *string {
	for _, huddle := range huddles {
		if huddle.Id == huddleId {
			return &huddle.What
		}
	}

	return nil
}

func getConfirmed(confirmedHuddles []uint, huddleId uint) *int {
	confirmed := 1 // Pass value as pointer
	for _, confirmedHuddle := range confirmedHuddles {
		if confirmedHuddle == huddleId {
			return &confirmed
		}
	}

	return nil
}
