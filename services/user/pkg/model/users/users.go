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

type UserGet struct {
	User           User  `json:"user"`
	People         int64 `json:"people"`
	Huddles        int64 `json:"huddles"`
	Notifications  int64 `json:"notifications"`
	UnreadMessages int64 `json:"unreadMessages"`
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
func GetUser(db *gorm.DB, username string) (UserGet, error) {
	var userGet UserGet

	var user User
	err := db.Table("users").Where("username = ?", username).First(&user).Error
	if err != nil {
		return userGet, err
	}

	var peopleCount int64
	if err := db.
		Table("notifications_people").
		Where("(sender = ? OR receiver = ?) AND type = 'person_invite' AND accepted = 1", username, username).
		Count(&peopleCount).Error; err != nil {
		return userGet, err
	}

	var huddlesCount int64 = 64

	var notificationsCount int64
	notificationsCountQuery :=
		`
			SELECT COUNT(*) FROM (SELECT
							id
						FROM
							notifications_people 
						WHERE
							receiver = ? AND seen = 0
						UNION ALL
						SELECT
							id
						FROM
							notifications_notify
						WHERE
							receiver = ? AND seen = 0
						UNION ALL
						SELECT
							id
						FROM
							notifications_huddles
						WHERE
							receiver = ? AND seen = 0) T
		`
	if err := db.
		Raw(notificationsCountQuery, username, username, username).
		First(&notificationsCount).Error; err != nil {
		return userGet, err
	}

	var unreadMessagesCount int64
	queryUnreadMessageCount := `SELECT COUNT(*) FROM (SELECT id AS message_id FROM messages WHERE id IN( SELECT MAX(id) FROM messages WHERE conversation_id IN( SELECT conversation_id FROM people_in_conversations WHERE username = '` + username + `') GROUP BY conversation_id)) T1 WHERE T1.message_id NOT IN( SELECT message_id FROM last_read_messages WHERE username = '` + username + `') GROUP BY message_id`
	if err := db.Raw(queryUnreadMessageCount).Scan(&unreadMessagesCount).Error; err != nil {
		return userGet, err
	}

	userGet = UserGet{
		User:           user,
		People:         peopleCount,
		Huddles:        huddlesCount,
		Notifications:  notificationsCount,
		UnreadMessages: unreadMessagesCount,
	}

	return userGet, nil
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
