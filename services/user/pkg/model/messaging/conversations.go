package messaging

import (
	"fmt"
	"time"

	p "github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"gorm.io/gorm"
)

const timeFormat = "2006-01-02 15:04:05"

type Conversation struct {
	Id uint `gorm:"primary_key;auto_increment;not_null"`
}

func (Conversation) TableName() string {
	return "conversations"
}

type Create struct {
	Sender   string
	Receiver string
}

type Chat struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	ProfilePhoto string `json:"profilePhoto,omitempty"`
	LastMessage  string `json:"lastMessage,omitempty"`
	IsRead       int    `json:"isRead,omitempty"`
	Time         string `json:"time"`
}

type LastMessage struct {
	Id             uint
	Sender         string
	ConversationId uint
	Message        string
	Time           int64
}

// Create conversation in conversations table
func CreateConversation(db *gorm.DB, t *Create) (uint, error) {
	conversaion := Conversation{}

	if err := db.Table("conversations").Create(&conversaion).Error; err != nil {
		return 0, err
	}

	create := []PersonInConversation{
		{ConversationId: conversaion.Id, Username: t.Sender},
		{ConversationId: conversaion.Id, Username: t.Receiver},
	}
	if err := db.Table("people_in_conversations").Create(&create).Error; err != nil {
		return 0, err
	}

	return conversaion.Id, nil
}

// Get chats from conversations table
func GetChats(db *gorm.DB, username string, lastId string) ([]Chat, error) {
	var chats []Chat
	var peopleInConversations []PersonInConversation
	var lastMessages []LastMessage
	var lastReadMessages []LastReadMessage
	var people []p.Person

	var idCondition string
	if lastId != "" {
		idCondition = fmt.Sprintf("id < %s AND ", lastId)
	}
	// Get last messages by username
	if err := db.
		Raw(`
			SELECT
				id,
				sender,
				conversation_id,
				message,
				time
			FROM
				messages
			WHERE
				`+idCondition+`id IN(
					SELECT
						MAX(id)
						FROM messages
					WHERE
						conversation_id IN (
							SELECT
								conversation_id FROM people_in_conversations
							WHERE
								username = ?)
						GROUP BY
							conversation_id)
			ORDER BY
				id DESC
			LIMIT 15`, username).
		Find(&lastMessages).Error; err != nil {
		return []Chat{}, err
	}

	if len(lastMessages) == 0 {
		return []Chat{}, nil
	}

	conversationsIds := getConversationsIds(lastMessages)

	// Get usernames from conversations
	if err := db.
		Raw(`
			SELECT
				conversation_id, username
			FROM
				people_in_conversations
			WHERE
				conversation_id IN ?
				AND username != ?`, conversationsIds, username).
		Find(&peopleInConversations).Error; err != nil {
		return []Chat{}, err
	}

	if err := db.
		Table("last_read_messages").
		Where("username = ? AND conversation_id IN ?", username, conversationsIds).
		Find(&lastReadMessages).
		Error; err != nil {
		return []Chat{}, err
	}

	usernames := getUsernamesFromPeopleInConversations(peopleInConversations)

	if err := db.
		Table("users").
		Select("username, firstname, profile_photo").
		Where("username IN ?", usernames).
		Find(&people).
		Error; err != nil {
		return []Chat{}, err
	}

	for _, lastMessage := range lastMessages {
		name, profilePhoto := getPeopleInfo(
			lastMessage.ConversationId,
			peopleInConversations,
			people,
		)
		isRead := getIsRead(lastReadMessages, lastMessage, username)
		time := time.Unix(lastMessage.Time, 0).Format(timeFormat)

		chats = append(chats, Chat{
			Id:           lastMessage.ConversationId,
			Name:         name,
			ProfilePhoto: profilePhoto,
			LastMessage:  lastMessage.Message,
			IsRead:       isRead,
			Time:         time,
		})
	}

	return chats, nil
}

func containsUint(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getConversationsIds(lastMessages []LastMessage) []uint {
	var ids []uint

	for _, lastMessage := range lastMessages {
		if !containsUint(ids, lastMessage.ConversationId) {
			ids = append(ids, lastMessage.ConversationId)
		}
	}

	return ids
}

func getUsernamesFromPeopleInConversations(peopleInConversations []PersonInConversation) []string {
	var usernames []string

	for _, person := range peopleInConversations {
		if !containsString(usernames, person.Username) {
			usernames = append(usernames, person.Username)
		}
	}

	return usernames
}

func getPeopleInfo(
	conversationId uint,
	peopleInConversations []PersonInConversation,
	people []p.Person,
) (string, string) {
	var name string
	var profilePhoto string

	for _, personInConversation := range peopleInConversations {
		if personInConversation.ConversationId == conversationId {
			for _, person := range people {
				if person.Username == personInConversation.Username {
					name = person.Firstname
					profilePhoto = person.ProfilePhoto

					break
				}
			}
		}
	}

	return name, profilePhoto
}

func getIsRead(lastReadMessages []LastReadMessage, lastMessage LastMessage, username string) int {
	if lastMessage.Sender == username {
		return 1
	}

	for _, lastReadMessage := range lastReadMessages {
		if lastReadMessage.MessageId == lastMessage.Id {
			return 1
		}
	}

	return 0
}
