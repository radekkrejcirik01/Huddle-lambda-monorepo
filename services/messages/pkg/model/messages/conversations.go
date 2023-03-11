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
	Usernames []string
	Username  string
}

func (ConversationsTable) TableName() string {
	return "conversations"
}

type Username struct {
	Username string
}

type ConversationList struct {
	Id        uint   `json:"id"`
	Usernames []User `json:"usernames"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	Message   string `json:"message"`
	Time      string `json:"time"`
	IsRead    uint   `json:"isRead"`
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

type GetConversation struct {
	ConversationId uint
	Username       string
}

type User struct {
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	ProfilePicture string `json:"profilePicture"`
}

type ConversationDetails struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Users   []User `json:"users,omitempty"`
}

type UserInfo struct {
	Firstname      string
	ProfilePicture string
}

type UpdateConversation struct {
	Id       uint
	Buffer   string
	FileName string
	Name     string
}

type Delete struct {
	ConversationId uint
	Username       string
}

// CreateConversation create conversation
func CreateConversation(db *gorm.DB, t *ConversationCreate) (ConversationDetails, error) {
	var usernames []string
	for _, user := range t.Usernames {
		usernames = append(usernames, `'`+user+`'`)
	}

	usernamesString := strings.Join(usernames, ", ")
	usernamesCount := len(t.Usernames)

	var conversationIds []uint
	if err := db.Table("people_in_conversations").Select("conversation_id").Where(`conversation_id IN( SELECT conversation_id FROM people_in_conversations WHERE username IN(` + usernamesString + `) GROUP BY conversation_id HAVING COUNT(conversation_id) = ` + strconv.Itoa(usernamesCount) + `) GROUP BY conversation_id HAVING COUNT(conversation_id) = ` + strconv.Itoa(usernamesCount)).Find(&conversationIds).Error; err != nil {
		return ConversationDetails{}, err
	}

	if len(conversationIds) > 0 {
		return GetDetails(db, conversationIds[0], t.Usernames, t.Username)
	}

	var newConversation ConversationsTable
	if err := db.Table("conversations").Create(&newConversation).Error; err != nil {
		return ConversationDetails{}, err
	}

	var peopleInConversations []PeopleInConversations
	for _, username := range t.Usernames {
		peopleInConversations = append(peopleInConversations, PeopleInConversations{
			ConversationId: newConversation.Id,
			Username:       username,
		})
	}

	if err := db.Table("people_in_conversations").Create(&peopleInConversations).Error; err != nil {
		return ConversationDetails{}, err
	}

	return GetDetails(db, newConversation.Id, t.Usernames, t.Username) // -1 to not include user
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
									username = '` + t.Username + `')`
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

	var deletedConversations []uint
	if err := db.Table("people_in_conversations").Select("conversation_id").Where("username = ? AND deleted = ?", t.Username, 1).Find(&deletedConversations).Error; err != nil {
		return []ConversationList{}, err
	}

	var result []ConversationList
	for _, message := range messages {
		if contains(deletedConversations, message.ConversationId) {
			continue
		}
		name, picture := getTitleAndPictureByConversationId(
			message.ConversationId,
			peopleInConversations,
			users,
			conversations,
			t.Username,
		)

		usernames := getUsersByConversationId(
			message.ConversationId,
			peopleInConversations,
			users,
		)

		result = append(result, ConversationList{
			Id:        message.ConversationId,
			Usernames: usernames,
			Name:      name,
			Picture:   picture,
			Message:   message.Message,
			Time:      message.Time,
			IsRead:    getIsRead(lastReads, message),
		})
	}

	return result, nil
}

// GetConversationDetails get conversation details from DB
func GetConversationDetails(db *gorm.DB, t *GetConversation) (ConversationDetails, error) {
	var peopleInConversation []PeopleInConversations
	if err := db.Table("people_in_conversations").Where("conversation_id = ?", t.ConversationId).Find(&peopleInConversation).Error; err != nil {
		return ConversationDetails{}, err
	}

	var usernames []string
	for _, user := range peopleInConversation {
		usernames = append(usernames, user.Username)
	}

	return GetDetails(db, t.ConversationId, usernames, t.Username)
}

func GetDetails(db *gorm.DB, conversationId uint, usernames []string, user string) (ConversationDetails, error) {
	var users []User
	var conversation Conversation
	var conversationDetails ConversationDetails
	if len(usernames) > 2 {
		if err := db.Table("conversations").Where(`id = ?`, conversationId).First(&conversation).Error; err != nil {
			return ConversationDetails{}, err
		}

		var usernamesArray []string
		for _, value := range usernames {
			usernamesArray = append(usernamesArray, `'`+value+`'`)
		}
		usernamesString := strings.Join(usernamesArray, ", ")
		if err := db.Table("users").Select("username, firstname, profile_picture").Where(`username IN (` + usernamesString + `)`).Find(&users).Error; err != nil {
			return ConversationDetails{}, err
		}
		conversationDetails = ConversationDetails{
			Id:      conversationId,
			Name:    conversation.Name,
			Picture: conversation.Picture,
			Users:   users,
		}
	} else {
		username := usernames[0]
		if username == user {
			username = usernames[1]
		}
		if err := db.Table("users").Select("username, firstname, profile_picture").Where("username = ?", username).First(&users).Error; err != nil {
			return ConversationDetails{}, err
		}
		conversation.Name = users[0].Firstname
		conversation.Picture = users[0].ProfilePicture
		conversationDetails = ConversationDetails{
			Id:      conversationId,
			Name:    conversation.Name,
			Picture: conversation.Picture,
		}
	}
	return conversationDetails, nil
}

// Update hangout in DB
func UpdateConversationById(db *gorm.DB, t *UpdateConversation) error {
	update := map[string]interface{}{}

	var photoUrl string
	if len(t.Buffer) > 0 {
		url, err := UplaodPhoto(db, strconv.Itoa(int(t.Id)), t.Buffer, t.FileName)
		if err != nil {
			return err
		}
		photoUrl = url
	}

	update["name"] = t.Name

	if len(photoUrl) > 0 {
		update["picture"] = photoUrl
	}

	return db.Table("conversations").Where("id = ?", t.Id).Updates(update).Error
}

// DeleteConversation delete conversation
func DeleteConversation(db *gorm.DB, t *Delete) error {
	return db.Table("people_in_conversations").Where("conversation_id = ? AND username = ?", t.ConversationId, t.Username).Update("deleted", 1).Error
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
		if !containsString(usernames, `'`+user.Username+`'`) {
			usernames = append(usernames, `'`+user.Username+`'`)
		}
	}
	usernamesString := strings.Join(usernames, ", ")

	return usernamesString
}

func getTitleAndPictureByConversationId(
	conversationId uint,
	peopleInConversations []PeopleInConversations,
	users []User,
	conversations []Conversation,
	username string,
) (string, string) {
	var people []string
	for _, person := range peopleInConversations {
		if person.ConversationId == conversationId {
			people = append(people, person.Username)
		}
	}

	for _, conversation := range conversations {
		if conversation.Id == conversationId {
			if len(people) > 2 {
				return conversation.Name, conversation.Picture
			} else {
				for _, user := range users {
					if user.Username != username {
						return user.Firstname, user.ProfilePicture
					}
				}
			}
		}
	}
	return "", ""
}

func getUsersByConversationId(
	conversationId uint,
	peopleInConversations []PeopleInConversations,
	users []User,
) []User {
	var usersInConversation []User
	for _, personInConversation := range peopleInConversations {
		if personInConversation.ConversationId == conversationId {
			for _, user := range users {
				if user.Username == personInConversation.Username {
					usersInConversation = append(usersInConversation, user)
				}
			}
		}
	}
	return usersInConversation
}

func contains(array []uint, value uint) bool {
	for _, a := range array {
		if a == value {
			return true
		}
	}

	return false
}

func containsString(array []string, value string) bool {
	for _, a := range array {
		if a == value {
			return true
		}
	}

	return false
}
