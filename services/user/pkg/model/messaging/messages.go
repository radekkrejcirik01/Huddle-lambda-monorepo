package messaging

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

type Message struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Sender         string
	ConversationId uint
	Message        string
	Time           int64 `gorm:"autoCreateTime"`
	Url            string
}

func (Message) TableName() string {
	return "messages"
}

type Send struct {
	ConversationId uint
	Message        string
	Time           string
	Buffer         *string
	FileName       *string
}

type Info struct {
	Username              string
	Firstname             string
	ProfilePhoto          string
	MessagesNotifications int
}

type Reaction struct {
	MessageId int
	Value     string
}

type MessageData struct {
	Id        uint     `json:"id"`
	Sender    string   `json:"sender"`
	Message   string   `json:"message"`
	Time      int64    `json:"time"`
	Url       string   `json:"url,omitempty"`
	Reactions []string `json:"reactions,omitempty"`
	ReadBy    []string `json:"readBy,omitempty"`
}

// Add message to messages table
func SendMessage(db *gorm.DB, username string, t *Send) error {
	var receiver string
	var photoUrl string

	if t.Buffer != nil {
		url, err := UploadChatPhoto(username, *t.Buffer, *t.FileName)
		if err != nil {
			return err
		}
		photoUrl = url
	}

	message := Message{
		Sender:         username,
		ConversationId: t.ConversationId,
		Message:        t.Message,
		Url:            photoUrl,
	}

	if err := db.Table("messages").Create(&message).Error; err != nil {
		return err
	}

	if err := db.
		Table("people_in_conversations").
		Select("username").
		Where("conversation_id = ? AND username != ?", t.ConversationId, username).
		Find(&receiver).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("last_read_messages").
		Where("username = ? AND conversation_id = ?", receiver, t.ConversationId).
		Update("seen", 0).
		Error; err != nil {
		return err
	}

	var mutedConversation []string
	if err := db.
		Table("muted_conversations").
		Select("user").
		Where("user = ? AND conversation_id = ?", receiver, t.ConversationId).
		Find(&mutedConversation).
		Error; err != nil {
		return err
	}

	if len(mutedConversation) > 0 {
		return nil
	}

	var info []Info
	if err := db.
		Table("users").
		Select("username, firstname, profile_photo, messages_notifications").
		Where("username IN ?", []string{receiver, username}).
		Find(&info).
		Error; err != nil {
		return err
	}

	if !receiveNotificationsEnabled(info, receiver) {
		return nil
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, receiver); err != nil {
		return nil
	}

	body := t.Message

	if len(t.Message) == 0 && t.Buffer != nil {
		body = "Sends photo"
	}

	senderInfo := getSenderInfo(info, username)

	fcmNotification := service.FcmNotification{
		Data: map[string]interface{}{
			"type":           "message",
			"conversationId": t.ConversationId,
			"name":           senderInfo.Firstname,
			"profilePhoto":   senderInfo.ProfilePhoto,
		},
		Title:   senderInfo.Firstname,
		Body:    body,
		Sound:   "default",
		Devices: *tokens,
	}

	return service.SendNotification(&fcmNotification)
}

// GetConversation messages and reactions from messages table
func GetConversation(db *gorm.DB, conversationId string, lastId string) ([]MessageData, error) {
	var messages []Message
	var messagesReactions []Reaction
	var messagesData []MessageData
	var lastReadMessages []LastReadMessage

	var idCondition string
	if lastId != "" {
		idCondition = fmt.Sprintf("id < %s AND ", lastId)
	}

	if err := db.
		Table("messages").
		Where(idCondition+"conversation_id = ?", conversationId).
		Order("id desc").
		Limit(20).
		Find(&messages).
		Error; err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, nil
	}

	messagesIds := getMessagesIds(messages)

	if err := db.
		Table("messages_reactions").
		Where("conversation_id = ? AND message_id IN ?", conversationId, messagesIds).
		Find(&messagesReactions).
		Error; err != nil {
		return nil, err
	}

	if err := db.
		Table("last_read_messages").
		Where("conversation_id = ?", conversationId).
		Find(&lastReadMessages).
		Error; err != nil {
		return nil, err
	}

	for _, message := range messages {
		reactions := getReactions(message.Id, messagesReactions)
		readBy := getReadBy(lastReadMessages, message.Id, message.Sender)

		messagesData = append(messagesData, MessageData{
			Id:        message.Id,
			Sender:    message.Sender,
			Message:   message.Message,
			Time:      message.Time,
			Url:       message.Url,
			Reactions: reactions,
			ReadBy:    readBy,
		})
	}

	return messagesData, nil
}

// GetMessagesByUsernames and reactions from messages table
func GetMessagesByUsernames(db *gorm.DB, username string, user string) ([]MessageData, uint, error) {
	var messages []Message
	var messagesReactions []Reaction
	var messagesData []MessageData
	var conversationId int
	var lastReadMessages []LastReadMessage

	// Get conversation id by 2 usernames
	if err := db.
		Table("people_in_conversations").
		Select("conversation_id").
		Where(`conversation_id IN(
			SELECT
				conversation_id FROM people_in_conversations
			WHERE
				username IN ?
			GROUP BY
				conversation_id
			HAVING
				COUNT(conversation_id) = 2)`, []string{username, user}).
		Find(&conversationId).
		Error; err != nil {
		return nil, 0, err
	}

	if conversationId == 0 {
		id, err := CreateConversation(db, &Create{
			Sender:   username,
			Receiver: user,
		})

		return nil, id, err
	}

	if err := db.
		Table("messages").
		Where("conversation_id = ?", conversationId).
		Order("id desc").
		Limit(20).
		Find(&messages).
		Error; err != nil {
		return nil, 0, err
	}

	if len(messages) == 0 {
		return nil, uint(conversationId), nil
	}

	messagesIds := getMessagesIds(messages)

	if err := db.
		Table("messages_reactions").
		Where("conversation_id = ? AND message_id IN ?", conversationId, messagesIds).
		Find(&messagesReactions).
		Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Table("last_read_messages").
		Where("conversation_id = ?", conversationId).
		Find(&lastReadMessages).
		Error; err != nil {
		return nil, 0, err
	}

	for _, message := range messages {
		reactions := getReactions(message.Id, messagesReactions)
		readBy := getReadBy(lastReadMessages, message.Id, message.Sender)

		messagesData = append(messagesData, MessageData{
			Id:        message.Id,
			Sender:    message.Sender,
			Message:   message.Message,
			Time:      message.Time,
			Url:       message.Url,
			Reactions: reactions,
			ReadBy:    readBy,
		})
	}

	return messagesData, uint(conversationId), nil
}

func UploadChatPhoto(username string, buffer string, fileName string) (string, error) {
	accessKey, secretAccessKey := database.GetCredentials()

	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String("eu-central-1"),
			Credentials: credentials.NewStaticCredentials(
				accessKey,
				secretAccessKey,
				"", // a token will be created when the session it's used.
			),
		}))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	decode, _ := base64.StdEncoding.DecodeString(buffer)
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("notify-bucket-images"),
		Key:         aws.String("messages-images/" + username + "/" + fileName),
		Body:        bytes.NewReader(decode),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}

func receiveNotificationsEnabled(info []Info, receiver string) bool {
	for _, i := range info {
		if i.Username == receiver && i.MessagesNotifications == 1 {
			return true
		}
	}
	return false
}

func getSenderInfo(info []Info, sender string) *Info {
	for _, i := range info {
		if i.Username == sender {
			return &i
		}
	}
	return nil
}

func getReactions(messageId uint, reactions []Reaction) []string {
	var r []string
	for _, reaction := range reactions {
		if reaction.MessageId == int(messageId) {
			r = append(r, reaction.Value)
		}
	}
	return r
}

func getMessagesIds(messages []Message) []uint {
	var ids []uint
	for _, message := range messages {
		ids = append(ids, message.Id)
	}
	return ids
}

func getReadBy(lastReadMessages []LastReadMessage, messageId uint, sender string) []string {
	var users []string
	for _, lastReadMessage := range lastReadMessages {
		if lastReadMessage.MessageId >= int(messageId) && lastReadMessage.Username != sender {
			users = append(users, lastReadMessage.Username)
		}
	}
	return users
}
