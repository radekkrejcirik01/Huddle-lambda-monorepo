package users

import (
	"bytes"
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/middleware"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/devices"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/huddles"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/messaging"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"gorm.io/gorm"
)

type User struct {
	Id                          uint   `gorm:"primary_key;auto_increment;not_null" json:"id"`
	Username                    string `json:"username"`
	Firstname                   string `json:"firstname"`
	ProfilePhoto                string `json:"profilePhoto"`
	FriendsInvitesNotifications int    `gorm:"default:1"`
	NewHuddlesNotifications     int    `gorm:"default:1"`
	InteractionsNotifications   int    `gorm:"default:1"`
	CommentsNotifications       int    `gorm:"default:1"`
	MentionsNotifications       int    `gorm:"default:1"`
	MessagesNotifications       int    `gorm:"default:1"`
	Password                    string
}

func (User) TableName() string {
	return "users"
}

type Login struct {
	Username string
	Password string
}

type UserData struct {
	Username     string `json:"username"`
	Firstname    string `json:"firstname"`
	ProfilePhoto string `json:"profilePhoto"`
}

type Notification struct {
	FriendsInvitesNotifications int `json:"friendsInvitesNotifications"`
	NewHuddlesNotifications     int `json:"newHuddlesNotifications"`
	InteractionsNotifications   int `json:"interactionsNotifications"`
	CommentsNotifications       int `json:"commentsNotifications"`
	MentionsNotifications       int `json:"mentionsNotifications"`
	MessagesNotifications       int `json:"messagesNotifications"`
}

type UpdateNotification struct {
	Notification string
	Value        int
}

type UploadProfilePhotoBody struct {
	Buffer   string
	FileName string
}

// CreateUser in users table
func CreateUser(db *gorm.DB, t *User) (string, error) {
	t.Password = middleware.GetHashPassword(t.Password)

	if rows := db.
		Table("users").
		Where("username = ?", t.Username).
		FirstOrCreate(&t).
		RowsAffected; rows == 0 {
		return "User already exists", nil
	}

	return "", nil
}

// LoginUser in users table
func LoginUser(db *gorm.DB, t *Login) error {
	var user User
	t.Password = middleware.GetHashPassword(t.Password)

	if err := db.
		Table("users").
		Where("username = ? AND password = ?", t.Username, t.Password).
		First(&user).
		Error; err != nil {
		return err
	}

	return nil
}

// GetUser from users table
func GetUser(db *gorm.DB, username string) (UserData, error) {
	var user UserData

	err := db.
		Table("users").
		Select("username, firstname, profile_photo").
		Where("username = ?", username).
		Find(&user).
		Error
	if err != nil {
		return UserData{}, err
	}

	return user, nil
}

// GetUserNotifications from users table
func GetUserNotifications(db *gorm.DB, username string) (Notification, error) {
	var notifications Notification

	if err := db.
		Table("users").
		Select(`
			friends_invites_notifications,
			new_huddles_notifications,
			interactions_notifications,
			comments_notifications,
			mentions_notifications,
			messages_notifications`).
		Where("username = ?", username).
		First(&notifications).
		Error; err != nil {
		return Notification{}, err
	}

	return notifications, nil
}

// UpdateUserNotification in users table
func UpdateUserNotification(db *gorm.DB, username string, t *UpdateNotification) error {
	return db.
		Table("users").
		Where("username = ?", username).
		Update(t.Notification, t.Value).
		Error
}

// DeleteAccount in database
func DeleteAccount(db *gorm.DB, username string) error {
	var deleteConversations []int

	if err := db.
		Table("people_in_conversations").
		Select("conversation_id").
		Where("username = ?", username).
		Find(&deleteConversations).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("people_in_conversations").
		Where("conversation_id IN ?", deleteConversations).
		Delete(&messaging.PersonInConversation{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("conversations").
		Where("id IN ?", deleteConversations).
		Delete(&messaging.Conversation{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("messages").
		Where("conversation_id IN ?", deleteConversations).
		Delete(&messaging.Message{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("invites").
		Where("sender = ? OR receiver = ?", username, username).
		Delete(&people.Invite{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("conversations_likes").
		Where("sender = ?", username).
		Delete(&people.Invite{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("devices").
		Where("username = ?", username).
		Delete(&devices.Device{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("huddles").
		Where("created_by = ?", username).
		Delete(&huddles.Huddle{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("huddles_comments").
		Where("sender = ?", username).
		Delete(&huddles.HuddleComment{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("huddles_comments_likes").
		Where("sender = ?", username).
		Delete(&huddles.HuddleCommentLike{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("huddles_interacted").
		Where("sender = ?", username).
		Delete(&huddles.HuddleInteracted{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("last_seen_messages").
		Where("conversation_id IN ?", deleteConversations).
		Delete(&messaging.LastSeenMessage{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("last_seen_huddles").
		Where("conversation_id IN ?", deleteConversations).
		Delete(&huddles.LastSeenHuddle{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("messages_reactions").
		Where("conversation_id IN ?", deleteConversations).
		Delete(&messaging.MessageReaction{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("muted_conversations").
		Where("conversation_id IN ?", deleteConversations).
		Delete(&people.MutedConversation{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("muted_huddles").
		Where("user = ?", username).
		Delete(&people.MutedHuddle{}).
		Error; err != nil {
		return err
	}

	if err := db.
		Table("hides").
		Where("user = ?", username).
		Delete(&people.Hide{}).
		Error; err != nil {
		return err
	}

	return db.
		Table("users").
		Where("username = ?", username).
		Delete(User{}).
		Error
}

// UploadProfilePhoto to S3 bucket
func UploadProfilePhoto(db *gorm.DB, username string, t *UploadProfilePhotoBody) (string, error) {
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

	decode, _ := base64.StdEncoding.DecodeString(t.Buffer)
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("notify-bucket-images"),
		Key:         aws.String("profile-images/" + username + "/" + t.FileName),
		Body:        bytes.NewReader(decode),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", err
	}

	if err := db.Table("users").Where("username = ?", username).Update("profile_photo", result.Location).Error; err != nil {
		return "", err
	}

	return result.Location, nil
}
