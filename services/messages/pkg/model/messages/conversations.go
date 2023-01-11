package messages

import (
	"strings"

	"github.com/radekkrejcirik01/Casblanca-backend/services/messages/pkg/model/helpers"
	"gorm.io/gorm"
)

type Email struct {
	Email string
}

type MessagedUser struct {
	Sender   string
	Receiver string
	Message  string
	Time     string
	IsRead   uint
}

type User struct {
	Email          string
	Firstname      string
	ProfilePicture string
}

type ConversationList struct {
	Email          string `json:"email"`
	Firstname      string `json:"firstname"`
	ProfilePicture string `json:"profilePicture"`
	Message        string `json:"message"`
	Time           string `json:"time"`
	IsRead         uint   `json:"isRead"`
}

// GetConversationsList get conversations
func GetConversationsList(db *gorm.DB, t *Email, page string) ([]ConversationList, error) {
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
										sender = '` + t.Email + `'
										OR receiver = '` + t.Email + `'
									GROUP BY
										( IF(sender = '` + t.Email + `', receiver, sender)))
							ORDER BY
								id DESC
							LIMIT 10 OFFSET ` + offset

	messagedUsers, err := GetConversationListFromQuery(db, messagedUsersQuery)
	if err != nil {
		return nil, err
	}

	formattedMessagedUsers := getFormattedMessagedUsers(messagedUsers, t.Email)

	userEmails := getUserEmailsString(formattedMessagedUsers)

	usersQuery := `SELECT email, firstname, profile_picture FROM users WHERE email IN (` + userEmails + `)`

	users, err := GetUsersFromQuery(db, usersQuery)
	if err != nil {
		return nil, err
	}

	var result []ConversationList
	for _, messagedUser := range formattedMessagedUsers {
		for _, user := range users {
			if messagedUser.Sender == user.Email {
				result = append(result, ConversationList{
					Email:          messagedUser.Sender,
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

func getFormattedMessagedUsers(messagedUsers []MessagedUser, email string) []MessagedUser {
	var result []MessagedUser
	for _, value := range messagedUsers {
		if value.Sender == email {
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

func getUserEmailsString(formattedMessagedUsers []MessagedUser) string {
	var emails []string
	for _, user := range formattedMessagedUsers {
		emails = append(emails, "'"+user.Sender+"'")
	}
	userEmails := strings.Join(emails, ", ")

	return userEmails
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
