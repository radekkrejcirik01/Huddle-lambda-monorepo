package users

import (
	"bytes"
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
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
}

type Notification struct {
	FriendsInvitesNotifications int `json:"friendsInvitesNotifications"`
	NewHuddlesNotifications     int `json:"newHuddlesNotifications"`
	InteractionsNotifications   int `json:"interactionsNotifications"`
	CommentsNotifications       int `json:"commentsNotifications"`
	MentionsNotifications       int `json:"mentionsNotifications"`
	MessagesNotifications       int `json:"messagesNotifications"`
}

func (User) TableName() string {
	return "users"
}

type UpdateNotification struct {
	Username     string
	Notification string
	Value        int
}

type UploadProfilePhotoBody struct {
	Username string
	Buffer   string
	FileName string
}

// CreateUser in users table
func CreateUser(db *gorm.DB, t *User) error {
	return db.Create(t).Error
}

// GetUser from users table
func GetUser(db *gorm.DB, username string) (User, error) {
	var user User
	err := db.Table("users").Where("username = ?", username).First(&user).Error
	if err != nil {
		return User{}, err
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
func UpdateUserNotification(db *gorm.DB, t *UpdateNotification) error {
	return db.
		Table("users").
		Where("username = ?", t.Username).
		Update(t.Notification, t.Value).
		Error
}

// UploadProfilePhoto to S3 bucket
func UploadProfilePhoto(db *gorm.DB, t *UploadProfilePhotoBody) (string, error) {
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
		Key:         aws.String("profile-images/" + t.Username + "/" + t.FileName),
		Body:        bytes.NewReader(decode),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", err
	}

	if err := db.Table("users").Where("username = ?", t.Username).Update("profile_photo", result.Location).Error; err != nil {
		return "", err
	}

	return result.Location, nil
}
