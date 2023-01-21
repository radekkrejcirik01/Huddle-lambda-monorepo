package users

import (
	"gorm.io/gorm"
)

type User struct {
	Id             uint   `gorm:"primary_key;auto_increment;not_null" json:"id"`
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	ProfilePicture string `json:"profilePicture"`
}

type UserGet struct {
	User           User  `json:"user"`
	People         int64 `json:"people"`
	Hangouts       int64 `json:"hangouts"`
	Notifications  int64 `json:"notifications"`
	UnreadMessages int64 `json:"unreadMessages"`
}

func (User) TableName() string {
	return "users"
}

// Create new User in DB
func CreateUser(db *gorm.DB, t *User) error {
	return db.Create(t).Error
}

// Get User from DB
func GetUser(db *gorm.DB, t *User) (UserGet, error) {
	err := db.Where("username = ?", t.Username).First(&t).Error
	if err != nil {
		return UserGet{}, err
	}

	var peopleCount int64
	if err := db.Table("people_invitations").
		Where("(username = ? AND confirmed = 1) OR (user = ? AND confirmed = 1)", t.Username, t.Username).
		Count(&peopleCount).Error; err != nil {
		return UserGet{}, err
	}

	var hangoutsCount int64
	if err := db.Table("hangouts T1").
		Joins("JOIN hangouts_invitations T2 ON T1.id = T2.hangout_id").
		Where("(T1.created_by = ? OR T2.username = ?) AND T2.confirmed = 1", t.Username, t.Username).
		Distinct("T1.id").
		Count(&hangoutsCount).Error; err != nil {
		return UserGet{}, err
	}

	var notificationsCount int64
	query := `SELECT
					id
					FROM
						people_invitations 
					WHERE
						username = '` + t.Username + `'
					UNION ALL
					SELECT
						id
					FROM
						hangouts_invitations
					WHERE
						username = '` + t.Username + `'`
	if err := db.Raw(query).Count(&notificationsCount).Error; err != nil {
		return UserGet{}, err
	}

	var unreadMessagesCount int64
	queryUnreadMessageCount := `SELECT COUNT(*) FROM (SELECT id AS message_id FROM messages WHERE id IN( SELECT MAX(id) FROM messages WHERE conversation_id IN( SELECT conversation_id FROM people_in_conversations WHERE username = '` + t.Username + `') GROUP BY conversation_id)) T1 WHERE T1.message_id NOT IN( SELECT message_id FROM last_read_messages WHERE username = '` + t.Username + `') GROUP BY message_id`
	if err := db.Raw(queryUnreadMessageCount).Scan(&unreadMessagesCount).Error; err != nil {
		return UserGet{}, err
	}

	result := UserGet{
		User:           *t,
		People:         peopleCount,
		Hangouts:       hangoutsCount,
		Notifications:  notificationsCount,
		UnreadMessages: unreadMessagesCount,
	}

	return result, nil
}
