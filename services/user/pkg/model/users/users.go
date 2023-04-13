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
	Id             uint   `gorm:"primary_key;auto_increment;not_null" json:"id"`
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	ProfilePicture string `json:"profilePicture"`
}

type UserGet struct {
	User           User  `json:"user"`
	People         int64 `json:"people"`
	Hangouts       int64 `json:"hangouts"`
	Notifications  int64 `json:"notifications"`
	UnreadMessages int64 `json:"unreadMessages"`
}

func (User) TableName() string {
	return "users"
}

type UplaodProfilePictureBody struct {
	Username string
	Buffer   string
	FileName string
}

// Create new User in DB
func CreateUser(db *gorm.DB, t *User) error {
	return db.Create(t).Error
}

// Get User from DB
func GetUser(db *gorm.DB, t *User) (UserGet, error) {
	err := db.Where("username = ?", t.Username).First(&t).Error
	if err != nil {
		return UserGet{}, err
	}

	var peopleCount int64
	if err := db.Table("people_invitations").
		Where("(username = ? AND confirmed = 1) OR (user = ? AND confirmed = 1)", t.Username, t.Username).
		Count(&peopleCount).Error; err != nil {
		return UserGet{}, err
	}

	var hangoutsCount int64 = 64

	var notificationsCount int64
	query := `
				SELECT
					COUNT(*)
				FROM (
					SELECT
					id
					FROM
						people_invitations 
					WHERE
						username = '` + t.Username + `' AND seen = 0
					UNION ALL
					SELECT
						id
					FROM
						hangouts_invitations
					WHERE
						receiver = '` + t.Username + `' AND seen = 0
					UNION ALL
					SELECT
						id
					FROM
						accepted_invitations
					WHERE
						username = '` + t.Username + `'
						AND seen = 0) T1`
	if err := db.Raw(query).Count(&notificationsCount).Error; err != nil {
		return UserGet{}, err
	}

	var unreadMessagesCount int64
	queryUnreadMessageCount := `SELECT COUNT(*) FROM (SELECT id AS message_id FROM messages WHERE id IN( SELECT MAX(id) FROM messages WHERE conversation_id IN( SELECT conversation_id FROM people_in_conversations WHERE username = '` + t.Username + `') GROUP BY conversation_id)) T1 WHERE T1.message_id NOT IN( SELECT message_id FROM last_read_messages WHERE username = '` + t.Username + `') GROUP BY message_id`
	if err := db.Raw(queryUnreadMessageCount).Scan(&unreadMessagesCount).Error; err != nil {
		return UserGet{}, err
	}

	result := UserGet{
		User:           *t,
		People:         peopleCount,
		Hangouts:       hangoutsCount,
		Notifications:  notificationsCount,
		UnreadMessages: unreadMessagesCount,
	}

	return result, nil
}

func UplaodProfilePicture(db *gorm.DB, t *UplaodProfilePictureBody) (string, error) {
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

	if err := db.Table("users").Where("username = ?", t.Username).Update("profile_picture", result.Location).Error; err != nil {
		return "", err
	}

	return result.Location, nil
}
