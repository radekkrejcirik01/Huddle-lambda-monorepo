package messaging

import (
	"errors"
	"fmt"

	p "github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"gorm.io/gorm"
)

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
	Id           int    `json:"id"`
	Sender       string `json:"sender"`
	Name         string `json:"name"`
	ProfilePhoto string `json:"profilePhoto,omitempty"`
	LastMessage  string `json:"lastMessage,omitempty"`
	IsNewMessage int    `json:"isNewMessage,omitempty"`
	IsRead       int    `json:"isRead,omitempty"`
	IsLiked      int    `json:"isLiked,omitempty"`
	Time         int64  `json:"time"`
}

type LastMessage struct {
	Id             uint
	Sender         string
	ConversationId int
	Message        string
	Time           int64
	Url            string
}

// Create conversation in conversations table
func CreateConversation(db *gorm.DB, t *Create) (uint, error) {
	conversation := Conversation{}

	if err := db.Table("conversations").Create(&conversation).Error; err != nil {
		return 0, err
	}

	create := []PersonInConversation{
		{ConversationId: int(conversation.Id), Username: t.Sender},
		{ConversationId: int(conversation.Id), Username: t.Receiver},
	}
	if err := db.Table("people_in_conversations").Create(&create).Error; err != nil {
		return 0, err
	}

	if err := db.Table("last_read_messages").Create(
		[]LastReadMessage{
			{Username: t.Sender, ConversationId: int(conversation.Id)},
			{Username: t.Receiver, ConversationId: int(conversation.Id)},
		}).Error; err != nil {
		return 0, err
	}

	return conversation.Id, nil
}

// GetUnreadMessagesNumber from messages and last_read_messages table
func GetUnreadMessagesNumber(db *gorm.DB, username string) (int64, error) {
	var number int64
	var lastMessagesIds []int64

	// Get last messages by username
	if err := db.
		Table("messages").
		Select("id").
		Where(`
					id IN(
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
							conversation_id)`, username).
		Find(&lastMessagesIds).Error; err != nil {
		return 0, err
	}

	if len(lastMessagesIds) == 0 {
		return 0, nil
	}

	if err := db.
		Table("last_read_messages").
		Where("username = ? AND message_id NOT IN ? AND seen != 1", username, lastMessagesIds).
		Count(&number).
		Error; err != nil {
		return 0, err
	}

	return number, nil
}

// Get chats from conversations table
func GetChats(db *gorm.DB, username string, lastId string) ([]Chat, error) {
	var chats []Chat
	var peopleInConversations []PersonInConversation
	var lastMessages []LastMessage
	var lastReadMessages []LastReadMessage
	var people []p.Person
	var likedConversations []ConversationLike

	var idCondition string
	if lastId != "" {
		idCondition = fmt.Sprintf("id < %s AND ", lastId)
	}
	// Get last messages by username
	if err := db.
		Table("messages").
		Select("id, sender, conversation_id, message, time, url").
		Where(idCondition+`
					id IN(
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
		return nil, err
	}

	if len(lastMessages) == 0 {
		return nil, nil
	}

	conversationsIds := getConversationsIds(lastMessages)

	// Get usernames from conversations
	if err := db.
		Table("people_in_conversations").
		Where("conversation_id IN ? AND username != ?", conversationsIds, username).
		Find(&peopleInConversations).
		Error; err != nil {
		return nil, err
	}

	if err := db.
		Table("last_read_messages").
		Where("conversation_id IN ?", conversationsIds).
		Find(&lastReadMessages).
		Error; err != nil {
		return nil, err
	}

	usernamesInConversations := getUsernamesFromPeopleInConversations(peopleInConversations)

	var invites []p.Invite
	if err := db.
		Table("invites").
		Where("((sender = ? AND receiver IN ?) OR (sender IN ? AND receiver = ?)) AND accepted = 1",
			username, usernamesInConversations, usernamesInConversations, username).
		Find(&invites).Error; err != nil {
		return nil, err
	}

	usernames := p.GetUsernamesFromInvites(invites, username)

	if err := db.
		Table("users").
		Select("username, firstname, profile_photo").
		Where("username IN ?", usernames).
		Find(&people).
		Error; err != nil {
		return nil, err
	}

	if err := db.
		Table("conversations_likes").
		Where("sender = ?", username).
		Find(&likedConversations).
		Error; err != nil {
		return nil, err
	}

	for _, lastMessage := range lastMessages {
		message := lastMessage.Message
		if len(lastMessage.Url) > 0 {
			message = "Photo shared"
		}
		name, profilePhoto, err := getPeopleInfo(
			lastMessage.ConversationId,
			peopleInConversations,
			people,
		)
		if err != nil {
			continue
		}

		isNewMessage := getIsNewMessage(lastReadMessages, lastMessage, username)
		isRead := getIsRead(lastReadMessages, lastMessage, username)
		isLiked := getIsLiked(lastMessage, likedConversations)

		chats = append(chats, Chat{
			Id:           lastMessage.ConversationId,
			Sender:       lastMessage.Sender,
			Name:         name,
			ProfilePhoto: profilePhoto,
			LastMessage:  message,
			IsNewMessage: isNewMessage,
			IsRead:       isRead,
			IsLiked:      isLiked,
			Time:         lastMessage.Time,
		})
	}

	return rearrangeChats(chats), nil
}

func containsInt(s []int, e int) bool {
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

func getConversationsIds(lastMessages []LastMessage) []int {
	var ids []int

	for _, lastMessage := range lastMessages {
		if !containsInt(ids, lastMessage.ConversationId) {
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
	conversationId int,
	peopleInConversations []PersonInConversation,
	people []p.Person,
) (string, string, error) {
	var name string
	var profilePhoto string
	err := errors.New("no info")

	for _, personInConversation := range peopleInConversations {
		if personInConversation.ConversationId == conversationId {
			for _, person := range people {
				if person.Username == personInConversation.Username {
					name = person.Firstname
					profilePhoto = person.ProfilePhoto
					err = nil

					break
				}
			}
		}
	}

	return name, profilePhoto, err
}

func getIsNewMessage(lastReadMessages []LastReadMessage, lastMessage LastMessage, username string) int {
	if lastMessage.Sender == username {
		return 0
	}

	for _, lastReadMessage := range lastReadMessages {
		if lastReadMessage.MessageId == int(lastMessage.Id) && lastReadMessage.Username == username {
			return 0
		}
	}
	return 1
}

func getIsRead(lastReadMessages []LastReadMessage, lastMessage LastMessage, username string) int {
	if lastMessage.Sender != username {
		return 0
	}

	for _, lastReadMessage := range lastReadMessages {
		if lastReadMessage.MessageId == int(lastMessage.Id) && lastReadMessage.Username != username {
			return 1
		}
	}
	return 0
}

func getIsLiked(lastMessage LastMessage, likedConversations []ConversationLike) int {
	for _, likedC := range likedConversations {
		if likedC.ConversationId == lastMessage.ConversationId {
			return 1
		}
	}
	return 0
}

func rearrangeChats(chats []Chat) []Chat {
	var liked []Chat
	var notLiked []Chat

	for _, chat := range chats {
		if chat.IsLiked == 1 {
			liked = append(liked, chat)
		} else {
			notLiked = append(notLiked, chat)
		}
	}

	return append(liked, notLiked...)
}
