package service

import (
	"log"

	"github.com/appleboy/go-fcm"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"gorm.io/gorm"
)

type FcmNotification struct {
	Data    map[string]interface{}
	Title   string
	Body    string
	Devices []string
}

func GetTokensByUsername(db *gorm.DB, t *[]string, username string) error {
	return db.Table("devices").Select("device_token").Where("username = ?", username).Find(t).Error
}

func GetTokensByUsernames(db *gorm.DB, usernames []string) ([]string, error) {
	var tokens []string
	if err := db.
		Table("devices").
		Select("device_token").
		Where("username IN ?", usernames).
		Find(&tokens).Error; err != nil {
		return tokens, err
	}
	return tokens, nil
}

func SendNotification(t *FcmNotification) error {
	fcmClient := database.GetFcmClient()
	tokens := t.Devices

	for _, token := range tokens {
		msg := &fcm.Message{
			To:   token,
			Data: t.Data,
			Notification: &fcm.Notification{
				Title: t.Title,
				Body:  t.Body,
				Badge: "1",
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
