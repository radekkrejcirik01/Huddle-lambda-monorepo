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
	Id           uint   `gorm:"primary_key;auto_increment;not_null" json:"id"`
	Username     string `json:"username"`
	Firstname    string `json:"firstname"`
	ProfilePhoto string `json:"profilePhoto"`
}

func (User) TableName() string {
	return "users"
}

type UplaodProfilePhotoBody struct {
	Username string
	Buffer   string
	FileName string
}

// Create new user in users table
func CreateUser(db *gorm.DB, t *User) error {
	return db.Create(t).Error
}

// Get user from users table
func GetUser(db *gorm.DB, username string) (User, error) {
	var user User
	err := db.Table("users").Where("username = ?", username).First(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func UplaodProfilePhoto(db *gorm.DB, t *UplaodProfilePhotoBody) (string, error) {
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
