package messages

import (
	"strings"

	"github.com/radekkrejcirik01/PingMe-backend/services/messages/pkg/model/helpers"
	"gorm.io/gorm"
)

type Username struct {
	Username string
}

type MessagedUser struct {
	Sender   string
	Receiver string
	Message  string
	Time     string
	IsRead   uint
}

type User struct {
	Username       string
	Firstname      string
	ProfilePicture string
}

type ConversationList struct {
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	ProfilePicture string `json:"profilePicture"`
	Message        string `json:"message"`
	Time           string `json:"time"`
	IsRead         uint   `json:"isRead"`
}

// GetConversationsList get conversations
func GetConversationsList(db *gorm.DB, t *Username, page string) ([]ConversationList, error) {
	offset := helpers.GetOffset(page)

	messagedUsersQuery := `SELECT
								sender,
								receiver,
								message,
								time,
								is_read
							FROM
								messages
							WHERE
								id IN(
									SELECT
										MAX(id)
										FROM messages
									WHERE
										sender = '` + t.Username + `'
										OR receiver = '` + t.Username + `'
									GROUP BY
										( IF(sender = '` + t.Username + `', receiver, sender)))
							ORDER BY
								id DESC
							LIMIT 10 OFFSET ` + offset

	messagedUsers, err := GetConversationListFromQuery(db, messagedUsersQuery)
	if err != nil {
		return nil, err
	}

	formattedMessagedUsers := getFormattedMessagedUsers(messagedUsers, t.Username)

	usernamesString := getUsernamesString(formattedMessagedUsers)

	usersQuery := `SELECT username, firstname, profile_picture FROM users WHERE username IN (` + usernamesString + `)`

	users, err := GetUsersFromQuery(db, usersQuery)
	if err != nil {
		return nil, err
	}

	var result []ConversationList
	for _, messagedUser := range formattedMessagedUsers {
		for _, user := range users {
			if messagedUser.Sender == user.Username {
				result = append(result, ConversationList{
					Username:       messagedUser.Sender,
					Firstname:      user.Firstname,
					ProfilePicture: user.ProfilePicture,
					Message:        messagedUser.Message,
					Time:           messagedUser.Time,
					IsRead:         messagedUser.IsRead,
				})
			}
		}
	}

	return result, nil
}

func GetConversationListFromQuery(db *gorm.DB, query string) ([]MessagedUser, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var array []MessagedUser
	for rows.Next() {
		db.ScanRows(rows, &array)
	}

	return array, nil
}

func getFormattedMessagedUsers(messagedUsers []MessagedUser, username string) []MessagedUser {
	var result []MessagedUser
	for _, value := range messagedUsers {
		if value.Sender == username {
			result = append(result, MessagedUser{
				Sender:   value.Receiver,
				Receiver: value.Sender,
				Message:  value.Message,
				Time:     value.Time,
				IsRead:   1,
			})
		} else {
			result = append(result, value)
		}
	}
	return result
}

func getUsernamesString(formattedMessagedUsers []MessagedUser) string {
	var usernames []string
	for _, user := range formattedMessagedUsers {
		usernames = append(usernames, "'"+user.Sender+"'")
	}
	usernamesString := strings.Join(usernames, ", ")

	return usernamesString
}

func GetUsersFromQuery(db *gorm.DB, query string) ([]User, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		db.ScanRows(rows, &users)
	}

	return users, nil
}
