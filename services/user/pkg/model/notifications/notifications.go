package notifications

import (
	"strings"

	"gorm.io/gorm"
)

type AcceptedInvitations struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	EventId  uint
	User     string
	Username string
	Time     string
	Type     string
	Seen     int `gorm:"default:0"`
}

func (AcceptedInvitations) TableName() string {
	return "accepted_invitations"
}

type Notification struct {
	Id        uint
	Username  string
	Time      string
	Confirmed int
	Type      string
}

type NotificationsData struct {
	Id             uint   `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	Time           string `json:"time"`
	ProfilePicture string `json:"profilePicture"`
	Confirmed      int    `json:"confirmed"`
	Type           string `json:"type"`
}

type Profile struct {
	Username       string
	Firstname      string
	ProfilePicture string
}

// Get notifications from DB
func GetNotifications(db *gorm.DB, t *Notification) ([]NotificationsData, error) {
	db.Transaction(func(tx *gorm.DB) error {
		tx.Table("people_invitations").Where("username = ?", t.Username).Update("seen", 1)
		tx.Table("hangouts_invitations").Where("username = ?", t.Username).Update("seen", 1)
		tx.Table("accepted_invitations").Where("username = ?", t.Username).Update("seen", 1)
		return nil
	})

	query := `
				SELECT
				id,
				user AS username,
				time,
				confirmed,
				'people' AS type
			FROM
				people_invitations
			WHERE
				username = '` + t.Username + `'
			UNION ALL
			SELECT
				hangout_id,
				user AS username,
				time,
				confirmed,
				'hangout' AS type
			FROM
				hangouts_invitations
			WHERE
				username = '` + t.Username + `'
			UNION ALL
			SELECT
				event_id,
				USER AS username,
				time,
				1 AS confirmed,
				TYPE
			FROM
				accepted_invitations
			WHERE
				username = '` + t.Username + `'
			ORDER BY
				time DESC`

	notifications, err := GetNotificationsFromQuery(db, query)
	if err != nil {
		return []NotificationsData{}, err
	}

	usernames := getUsernames(notifications)

	queryProfiles := `SELECT firstname, username, profile_picture FROM users WHERE username IN (` + usernames + `)`
	profiles, err := GetProfilePicturesFromQuery(db, queryProfiles)
	if err != nil {
		return []NotificationsData{}, err
	}

	var result []NotificationsData
	for _, notification := range notifications {
		for _, profile := range profiles {
			if notification.Username == profile.Username {
				result = append(result, NotificationsData{
					Id:             notification.Id,
					Username:       notification.Username,
					Name:           profile.Firstname,
					ProfilePicture: profile.ProfilePicture,
					Time:           notification.Time,
					Confirmed:      notification.Confirmed,
					Type:           notification.Type,
				})
			}
		}
	}

	return result, nil
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

func getUsernames(notifications []Notification) string {
	var usersnames []string
	for _, notification := range notifications {
		usersnames = append(usersnames, `'`+notification.Username+`'`)
	}

	return strings.Join(usersnames, ", ")
}

func GetNotificationsFromQuery(db *gorm.DB, query string) ([]Notification, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		db.ScanRows(rows, &notifications)
	}

	return notifications, nil
}
