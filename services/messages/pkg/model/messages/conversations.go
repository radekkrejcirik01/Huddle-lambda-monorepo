package messages

import (
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ConversationsTable struct {
	Id      uint `gorm:"primary_key;auto_increment;not_null"`
	Name    string
	Picture string
}

type ConversationCreate struct {
	Id        uint
	Name      string
	Usernames []string
}

func (ConversationsTable) TableName() string {
	return "conversations"
}

type Username struct {
	Username string
}

type ConversationList struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Message string `json:"message"`
	Time    string `json:"time"`
	IsRead  uint   `json:"isRead"`
}

type Messages struct {
	Id             uint
	Sender         string
	ConversationId uint
	Message        string
	Time           string
}

type Conversation struct {
	Id      uint
	Name    string
	Picture string
}

type User struct {
	Username       string
	Firstname      string
	ProfilePicture string
}

// CreateConversation create conversation
func CreateConversation(db *gorm.DB, t *ConversationCreate) (uint, error) {
	var usernames []string
	for _, user := range t.Usernames {
		usernames = append(usernames, `'`+user+`'`)
	}

	usernamesString := strings.Join(usernames, ", ")
	usernamesCount := len(t.Usernames)

	var conversationIds []uint
	if err := db.Table("people_in_conversations").Select("conversation_id").Where(`conversation_id IN( SELECT conversation_id FROM people_in_conversations WHERE username IN(` + usernamesString + `) GROUP BY conversation_id HAVING COUNT(conversation_id) = ` + strconv.Itoa(usernamesCount) + `) GROUP BY conversation_id HAVING COUNT(conversation_id) = ` + strconv.Itoa(usernamesCount)).Find(&conversationIds).Error; err != nil {
		return 0, err
	}

	if len(conversationIds) > 0 {
		return conversationIds[0], nil
	}

	if len(t.Usernames) > 2 {
		var firstnames []string
		if err := db.Table("users").Select("firstname").Where(`username IN (` + usernamesString + `)`).Find(&firstnames).Error; err != nil {
			return 0, err
		}

		var name string
		for i, firstname := range firstnames {
			if i == 0 {
				name = firstname
			} else {
				name += ` + ` + firstname
			}
		}
		t.Name = name
	}

	if err := db.Table("conversations").Create(t).Error; err != nil {
		return 0, err
	}

	var peopleInConversations []PeopleInConversations
	for _, username := range t.Usernames {
		peopleInConversations = append(peopleInConversations, PeopleInConversations{
			ConversationId: t.Id,
			Username:       username,
		})
	}

	if err := db.Table("people_in_conversations").Create(&peopleInConversations).Error; err != nil {
		return 0, err
	}

	return t.Id, nil
}

// GetConversationsList get conversations
func GetConversationsList(db *gorm.DB, t *Username, page string) ([]ConversationList, error) {
	queryGetMessages := `
			SELECT
				id,
				sender,
				conversation_id,
				message,
				time
			FROM
				messages
			WHERE
				id IN(
					SELECT
						MAX(id)
						FROM messages
					WHERE
						conversation_id IN(
							SELECT
								conversation_id FROM people_in_conversations
							WHERE
								username = '` + t.Username + `')
						GROUP BY
							conversation_id)
			ORDER BY
				id DESC`
	messages, err := GetMessagesFromQuery(db, queryGetMessages)
	if err != nil {
		return []ConversationList{}, err
	}

	// Get non group chat users
	queryGetPeopleInConversations := `
						SELECT
							*
						FROM
							people_in_conversations
						WHERE
							conversation_id IN(
								SELECT
									conversation_id FROM people_in_conversations
								WHERE
									username = '` + t.Username + `')
							AND username != '` + t.Username + `'
						GROUP BY
							conversation_id
						HAVING
							COUNT(conversation_id) = 1`
	peopleInConversations, err := GetPeopleInConversationsFromQuery(db, queryGetPeopleInConversations)
	if err != nil {
		return []ConversationList{}, err
	}

	queryGetConversations := `SELECT * FROM conversations WHERE id IN (SELECT conversation_id FROM people_in_conversations WHERE username = '` + t.Username + `')`
	conversations, err := GetConversationsFromQuery(db, queryGetConversations)
	if err != nil {
		return []ConversationList{}, err
	}

	usernamesString := getUsernamesString(peopleInConversations)

	var users []User
	if len(usernamesString) > 0 {
		queryGetUsers := `SELECT username, firstname, profile_picture FROM users WHERE username IN (` + usernamesString + `)`
		usersFromQuery, err := GetUsersFromQuery(db, queryGetUsers)
		if err != nil {
			return []ConversationList{}, err
		}
		users = usersFromQuery
	}

	var lastReads []LastReadMessage
	if err := db.Table("last_read_messages").Where("username = ?", t.Username).Find(&lastReads).Error; err != nil {
		return []ConversationList{}, err
	}

	var result []ConversationList
	for _, message := range messages {
		name, picture := getNameByConversationId(
			message.ConversationId,
			peopleInConversations,
			users,
			conversations)
		result = append(result, ConversationList{
			Id:      message.ConversationId,
			Name:    name,
			Picture: picture,
			Message: message.Message,
			Time:    message.Time,
			IsRead:  getIsRead(lastReads, message),
		})
	}

	return result, nil
}

func getIsRead(lastReads []LastReadMessage, message Messages) uint {
	for _, lastRead := range lastReads {
		if lastRead.ConversationId == message.ConversationId {
			if lastRead.MessageId == message.Id {
				return 1
			}
		}
	}
	return 0
}

func GetMessagesFromQuery(db *gorm.DB, query string) ([]Messages, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []Messages
	for rows.Next() {
		db.ScanRows(rows, &messages)
	}

	return messages, nil
}

func GetPeopleInConversationsFromQuery(db *gorm.DB, query string) ([]PeopleInConversations, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var peopleInConversations []PeopleInConversations
	for rows.Next() {
		db.ScanRows(rows, &peopleInConversations)
	}

	return peopleInConversations, nil
}

func GetConversationsFromQuery(db *gorm.DB, query string) ([]Conversation, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		db.ScanRows(rows, &conversations)
	}

	return conversations, nil
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

func getUsernamesString(peopleInConversations []PeopleInConversations) string {
	var usernames []string
	for _, user := range peopleInConversations {
		usernames = append(usernames, `'`+user.Username+`'`)
	}
	usernamesString := strings.Join(usernames, ", ")

	return usernamesString
}

func getNameByConversationId(
	conversationId uint,
	peopleInConversations []PeopleInConversations,
	users []User,
	conversations []Conversation,
) (string, string) {
	for _, conversation := range conversations {
		for _, personInConversation := range peopleInConversations {
			if personInConversation.ConversationId == conversationId {
				for _, user := range users {
					if user.Username == personInConversation.Username {
						return user.Firstname, user.ProfilePicture
					}
				}
			}
		}
		if conversation.Id == conversationId {
			return conversation.Name, conversation.Picture
		}
	}
	return "", ""
}
