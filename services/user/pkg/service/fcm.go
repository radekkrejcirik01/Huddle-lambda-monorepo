package service

import (
	"log"

	"github.com/appleboy/go-fcm"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"gorm.io/gorm"
)

type FcmNotification struct {
	Sender  string
	Type    string
	Title   string
	Body    string
	Sound   string
	Devices []string
}

func GetTokensByUsername(db *gorm.DB, t *[]string, username string) error {
	return db.Table("devices").Select("device_token").Where("username = ?", username).Find(t).Error
}

func GetTokensByUsernames(db *gorm.DB, t *[]string, usernames string) error {
	return db.Table("devices").Select("device_token").Where(`username IN (` + usernames + `)`).Find(t).Error
}

func SendNotification(t *FcmNotification) error {
	fcmClient := database.GetFcmClient()
	tokens := t.Devices

	for _, token := range tokens {
		msg := &fcm.Message{
			To: token,
			Data: map[string]interface{}{
				"type":   t.Type,
				"sender": t.Sender,
			},
			Notification: &fcm.Notification{
				Title: t.Title,
				Body:  t.Body,
				Sound: t.Sound,
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
