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
	IsRead         uint
}

func (Message) TableName() string {
	return "messages"
}

type ConversationId struct {
	ConversationId uint
}

type MessagesBody struct {
	Username string
	User     string
}

type MessageResponse struct {
	Id             uint   `gorm:"primary_key;auto_increment;not_null" json:"id"`
	Sender         string `json:"sender"`
	ProfilePicture string `json:"profilePicture"`
	ConversationId uint   `json:"conversationId"`
	Message        string `json:"message"`
	Time           string `json:"time"`
	IsRead         uint   `gorm:"default:0" json:"isRead"`
}
type SentMessage struct {
	Sender         string
	ConversationId uint
	Message        string
	Time           string
	IsRead         uint
}

type Notification struct {
	Sender  string
	Title   string
	Body    string
	Devices []string
}

func GetMessages(db *gorm.DB, t *ConversationId) ([]MessageResponse, error) {
	var messages []Message
	if err := db.Where("conversation_id = ?", t.ConversationId).Order("id DESC").Find(&messages).Error; err != nil {
		return []MessageResponse{}, err
	}

	var usernames []string
	for _, user := range messages {
		if !contains(usernames, user.Sender) {
			usernames = append(usernames, `'`+user.Sender+`'`)
		}
	}

	usernamesString := strings.Join(usernames, ", ")

	var users []User
	if err := db.Table("users").Select("username, firstname, profile_picture").Where(`username IN (` + usernamesString + `)`).Find(&users).Error; err != nil {
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
					IsRead:         message.IsRead,
				})
			}
		}
	}

	return result, nil
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

	if err := db.Table("messages").Create(&create).Error; err == nil {
		tokens := &[]string{}
		if err := GetUserTokensByUser(db, tokens, t.Sender); err != nil {
			return nil
		}
		notification := Notification{
			Sender:  t.Sender,
			Title:   t.Sender,
			Body:    t.Message,
			Devices: *tokens,
		}

		SendNotification(&notification)
		return err
	}

	return nil
}

func GetUserTokensByUser(db *gorm.DB, t *[]string, user string) error {
	return db.Table("devices").Select("device_token").Where("username = ?", user).Find(t).Error
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
				Sound: "default",
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

// UpdateRead update message as read
func UpdateRead(db *gorm.DB, t *MessagesBody) error {
	return db.Table("messages").Where("sender = ? AND receiver = ?", t.User, t.Username).Update("is_read", 1).Error
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
