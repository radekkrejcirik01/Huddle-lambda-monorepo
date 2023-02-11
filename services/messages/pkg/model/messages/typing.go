package messages

import (
	"log"
	"strconv"
	"strings"

	"github.com/appleboy/go-fcm"
	"github.com/radekkrejcirik01/PingMe-backend/services/messages/pkg/database"
	"gorm.io/gorm"
)

type Typing struct {
	ConversationId string
	Username       string
	Value          int
}

type TypingNotification struct {
	Username       string
	Value          string
	ConversationId string
	Tokens         []string
}

func SendTyping(db *gorm.DB, t *Typing) error {
	var usernames []string
	if err := db.Table("people_in_conversations").Select("username").Where("conversation_id = ? AND username != ?", t.ConversationId, t.Username).Find(&usernames).Error; err != nil {
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

	payload := TypingNotification{
		Username:       t.Username,
		Value:          strconv.Itoa(t.Value),
		ConversationId: t.ConversationId,
		Tokens:         *tokens,
	}

	return SendTypingNotification(payload)
}

func SendTypingNotification(t TypingNotification) error {
	fcmClient := database.GetFcmClient()

	for _, token := range t.Tokens {
		msg := &fcm.Message{
			To: token,
			Data: map[string]interface{}{
				"username":       t.Username,
				"value":          t.Value,
				"type":           "typing",
				"conversationId": t.ConversationId,
			},
			ContentAvailable: true,
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
