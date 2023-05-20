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
	Url            *string
}

func (Message) TableName() string {
	return "messages"
}

type Send struct {
	Sender         string
	Name           string
	ConversationId uint
	Message        string
	Time           string
	Buffer         *string
	FileName       *string
}

type MessageData struct {
	Id      uint   `json:"id"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
	Time    string `json:"time"`
	Url     string `json:"url,omitempty"`
}

// Add message to messages table
func SendMessage(db *gorm.DB, t *Send) error {
	var receiver string
	var photoUrl string

	if t.Buffer != nil {
		url, err := UplaodChatPhoto(db, t.Sender, *t.Buffer, *t.FileName)
		if err != nil {
			return err
		}
		photoUrl = url
	}

	message := Message{
		Sender:         t.Sender,
		ConversationId: t.ConversationId,
		Message:        t.Message,
		Url:            &photoUrl,
	}

	if err := db.Table("messages").Create(&message).Error; err != nil {
		return err
	}

	if err := db.
		Table("people_in_conversations").
		Select("username").
		Where("conversation_id = ? AND username != ?", t.ConversationId, t.Sender).
		Find(&receiver).
		Error; err != nil {
		return err
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, receiver); err != nil {
		return nil
	}

	body := t.Message

	if len(t.Message) == 0 && t.Buffer != nil {
		body = "Sends a photo"
	}

	notification := service.FcmNotification{
		Sender:  t.Sender,
		Type:    "message",
		Title:   t.Name,
		Body:    body,
		Sound:   "notification.wav",
		Devices: *tokens,
	}

	return service.SendNotification(&notification)
}

// Get conversation messages from messages table
func GetConversation(db *gorm.DB, conversaionId int, lastId string) ([]MessageData, error) {
	var message []MessageData

	var idCondition string
	if lastId != "" {
		idCondition = fmt.Sprintf("id < %s AND ", lastId)
	}

	if err := db.
		Table("messages").
		Where(idCondition+"conversation_id = ?", conversaionId).
		Order("id desc").
		Limit(20).
		Find(&message).
		Error; err != nil {
		return nil, err
	}

	return message, nil
}

// Get messages by usernames from messages table
func GetMessagesByUsernames(db *gorm.DB, user1 string, user2 string) ([]MessageData, uint, error) {
	var message []MessageData
	var conversaionId int

	// Get conversation id by 2 usernames
	if err := db.
		Raw(`
			SELECT
				conversation_id
			FROM
				people_in_conversations
			WHERE
				conversation_id IN(
					SELECT
						conversation_id FROM people_in_conversations
					WHERE
						username IN ?
					GROUP BY
						conversation_id
					HAVING
						COUNT(conversation_id) = 2)
			GROUP BY
				conversation_id
			HAVING
				COUNT(conversation_id) = 2`, []string{user1, user2}).
		Find(&conversaionId).
		Error; err != nil {
		return nil, 0, err
	}

	if conversaionId == 0 {
		id, err := CreateConversation(db, &Create{
			Sender:   user1,
			Receiver: user2,
		})

		return nil, id, err
	}

	if err := db.
		Table("messages").
		Where("conversation_id = ?", conversaionId).
		Order("id desc").
		Limit(20).
		Find(&message).
		Error; err != nil {
		return nil, 0, err
	}

	return message, uint(conversaionId), nil
}

func UplaodChatPhoto(db *gorm.DB, username string, buffer string, fileName string) (string, error) {
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
