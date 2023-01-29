package messages

import (
	"log"
	"strings"
	"time"

	"github.com/appleboy/go-fcm"
	"github.com/radekkrejcirik01/PingMe-backend/services/messages/pkg/database"
	"gorm.io/gorm"
)

const timeFormat = "2006-01-02 15:04:05"

type Message struct {
	Id             uint `gorm:"primary_key;auto_increment;not_null"`
	Sender         string
	ConversationId uint
	Message        string
	Time           string
}

func (Message) TableName() string {
	return "messages"
}

type ConversationId struct {
	ConversationId uint
}

type MessageResponse struct {
	Id             uint     `gorm:"primary_key;auto_increment;not_null" json:"id"`
	Sender         string   `json:"sender"`
	ProfilePicture string   `json:"profilePicture"`
	ConversationId uint     `json:"conversationId"`
	Message        string   `json:"message"`
	Time           string   `json:"time"`
	ReadBy         []ReadBy `json:"readBy"`
}
type SentMessage struct {
	Sender         string
	Name           string
	ConversationId uint
	Message        string
	Time           string
}

type Notification struct {
	Sender  string
	Title   string
	Body    string
	Devices []string
}

type ReadBy struct {
	Username       string `json:"username"`
	ProfilePicture string `json:"profilePicture"`
}

func GetMessages(db *gorm.DB, t *ConversationId) ([]MessageResponse, error) {
	var messages []Message
	if err := db.Where("conversation_id = ?", t.ConversationId).Order("id DESC").Find(&messages).Error; err != nil {
		return []MessageResponse{}, err
	}

	var usernames []string
	if err := db.Table("people_in_conversations").Select("username").Where("conversation_id = ?", t.ConversationId).Find(&usernames).Error; err != nil {
		return []MessageResponse{}, err
	}

	var usernamesArray []string
	for _, username := range usernames {
		usernamesArray = append(usernamesArray, `'`+username+`'`)
	}

	usernamesString := strings.Join(usernamesArray, ", ")

	var users []User
	if err := db.Table("users").Select("username, firstname, profile_picture").Where(`username IN (` + usernamesString + `)`).Find(&users).Error; err != nil {
		return []MessageResponse{}, err
	}

	var lastReads []LastReadMessage
	if err := db.Table("last_read_messages").Where("conversation_id = ?", t.ConversationId).Find(&lastReads).Error; err != nil {
		return []MessageResponse{}, err
	}

	var result []MessageResponse
	for _, message := range messages {
		for _, user := range users {
			if message.Sender == user.Username {
				result = append(result, MessageResponse{
					Id:             message.Id,
					Sender:         message.Sender,
					ProfilePicture: user.ProfilePicture,
					ConversationId: message.ConversationId,
					Message:        message.Message,
					Time:           message.Time,
					ReadBy:         getReadByMessageId(lastReads, message, users),
				})
			}
		}
	}

	return result, nil
}

func getReadByMessageId(lastReads []LastReadMessage, message Message, users []User) []ReadBy {
	var readBy []ReadBy
	for _, lastRead := range lastReads {
		if lastRead.MessageId == message.Id {
			for _, user := range users {
				if user.Username == lastRead.Username && message.Sender != user.Username {
					readBy = append(readBy, ReadBy{
						Username:       user.Username,
						ProfilePicture: user.ProfilePicture,
					})
				}
			}
		}
	}

	return readBy
}

// SendMessage send message
func SendMessage(db *gorm.DB, t *SentMessage) error {
	time := time.Now()
	t.Time = time.Format(timeFormat)
	create := Message{
		Sender:         t.Sender,
		ConversationId: t.ConversationId,
		Message:        t.Message,
		Time:           t.Time,
	}

	err := db.Table("messages").Create(&create).Error
	if err == nil {
		var usernames []string
		if err := db.Table("people_in_conversations").Select("username").Where("conversation_id = ? AND username != ?", t.ConversationId, t.Sender).Find(&usernames).Error; err != nil {
			return err
		}

		var usernamesArray []string
		for _, username := range usernames {
			usernamesArray = append(usernamesArray, `'`+username+`'`)
		}

		usernamesString := strings.Join(usernamesArray, ", ")

		tokens := &[]string{}
		if err := GetTokensByUsernames(db, tokens, usernamesString); err != nil {
			return nil
		}
		notification := Notification{
			Sender:  t.Sender,
			Title:   t.Name,
			Body:    t.Message,
			Devices: *tokens,
		}

		SendNotification(&notification)
		return nil
	}
	return err
}

func GetTokensByUsernames(db *gorm.DB, t *[]string, usernames string) error {
	return db.Table("devices").Select("device_token").Where(`username IN (` + usernames + `)`).Find(t).Error
}

func SendNotification(t *Notification) error {
	fcmClient := database.GetFcmClient()
	tokens := t.Devices

	for _, token := range tokens {
		msg := &fcm.Message{
			To: token,
			Data: map[string]interface{}{
				"type":   "message",
				"sender": t.Sender,
			},
			Notification: &fcm.Notification{
				Title: t.Title,
				Body:  t.Body,
				Badge: "1",
				Sound: "notification.wav",
			},
		}

		client, err := fcm.NewClient(fcmClient)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		response, err := client.Send(msg)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		log.Printf("%#v\n", response)
	}

	return nil
}
