package messages

import (
	"log"

	"github.com/appleboy/go-fcm"
	"github.com/radekkrejcirik01/PingMe-backend/services/messages/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/messages/pkg/model/helpers"
	"gorm.io/gorm"
)

type MessagesBody struct {
	Email string
	User  string
}

type Message struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Sender   string
	Receiver string
	Message  string
	Time     string
	IsRead   uint `gorm:"default:0"`
}

func (Message) TableName() string {
	return "messages"
}

type SentMessage struct {
	Sender          string
	SenderFirstname string
	Receiver        string
	Message         string
	Time            string
	IsRead          uint
}

type Messages struct {
	Id       uint   `json:"id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
	Time     string `json:"time"`
}

type Notification struct {
	Sender  string
	Title   string
	Body    string
	Devices []string
}

// GetMessages get messages
func GetMessages(db *gorm.DB, t *MessagesBody, page string) ([]Messages, error) {
	offset := helpers.GetOffset(page)

	messagesQuery := `SELECT
							id,
							sender,
							receiver,
							message,
							time
						FROM
							messages
						WHERE (sender = '` + t.Email + `'
							AND receiver = '` + t.User + `')
							OR(sender = '` + t.User + `'
								AND receiver = '` + t.Email + `')
						ORDER BY
							id DESC
						LIMIT 10 OFFSET ` + offset

	messages, err := GetMessagesFromQuery(db, messagesQuery)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// SendMessage send message
func SendMessage(db *gorm.DB, t *SentMessage) error {
	create := Message{
		Sender:   t.Sender,
		Receiver: t.Receiver,
		Message:  t.Message,
		Time:     t.Time,
	}

	err := db.Select("sender", "receiver", "message", "time").Create(&create).Error
	if err == nil {
		tokens := &[]string{}
		if err := GetUserTokensByUser(db, tokens, t.Receiver); err != nil {
			return err
		}
		notification := Notification{
			Sender:  t.Sender,
			Title:   t.SenderFirstname,
			Body:    t.Message,
			Devices: *tokens,
		}

		SendNotification(&notification)
		return err
	}

	return err
}

func GetMessagesFromQuery(db *gorm.DB, query string) ([]Messages, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var array []Messages
	for rows.Next() {
		db.ScanRows(rows, &array)
	}

	return array, nil
}

func GetUserTokensByUser(db *gorm.DB, t *[]string, user string) error {
	return db.Table("devices").Select("device_token").Where("email = ?", user).Find(t).Error
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
	return db.Table("messages").Where("sender = ? AND receiver = ?", t.User, t.Email).Update("is_read", 1).Error
}
