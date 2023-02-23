package hangouts

import (
	"bytes"
	"encoding/base64"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

const hangoutType = "hangout"
const groupHangoutType = "group_hangout"

const timeFormat = "2006-01-02 15:04:05"

type HangoutsTable struct {
	Id        uint   `gorm:"primary_key;auto_increment;not_null" json:"id"`
	CreatedBy string `json:"createdBy"`
	Title     string `json:"title"`
	Time      string `json:"time"`
	Place     string `json:"place"`
	Picture   string `json:"picture"`
	Type      string `json:"type"`
}

type HangoutInvite struct {
	User     string
	Name     string
	Username string
	Time     string
	Place    string
}

type GroupHangoutInvite struct {
	User      string
	Name      string
	Title     string
	Usernames []string
	Time      string
	Place     string
	Buffer    string
	FileName  string
}

type GetHangout struct {
	Username string
	ShowAll  bool
}

type Hangouts struct {
	Title string `json:"title"`
	Data  []List `json:"data"`
}

type List struct {
	List []HangoutsTable `json:"list"`
}

func (HangoutsTable) TableName() string {
	return "hangouts"
}

type Notification struct {
	Sender  string
	Title   string
	Body    string
	Devices []string
}

// Create new hangout in DB
func CreateHangout(db *gorm.DB, t *HangoutInvite) error {
	hangout := HangoutsTable{
		CreatedBy: t.User,
		Time:      t.Time,
		Place:     t.Place,
		Type:      hangoutType,
	}
	if err := db.Table("hangouts").Create(&hangout).Error; err != nil {
		return err
	}

	now := time.Now().Format(timeFormat)
	hangoutInvitation := HangoutsInvitationTable{
		HangoutId: hangout.Id,
		User:      t.User,
		Username:  t.Username,
		Time:      now,
		Confirmed: 0,
		Type:      hangoutType,
	}

	if err := db.Table("hangouts_invitations").Create(&hangoutInvitation).Error; err != nil {
		return err
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Username); err != nil {
		return nil
	}
	hangoutNotification := service.FcmNotification{
		Sender:  t.User,
		Type:    hangoutType,
		Title:   t.Name + " sends a hangout!",
		Sound:   "notification.wav",
		Devices: *tokens,
	}
	service.SendNotification(&hangoutNotification)
	return nil
}

// Create new group hangout in DB
func CreateGroupHangout(db *gorm.DB, t *GroupHangoutInvite) error {
	title := "Group hangout"
	if len(t.Title) > 0 {
		title = t.Title
	}

	var photoUrl string
	if len(t.Buffer) > 0 {
		url, err := UplaodPhoto(db, t.User, t.Buffer, t.FileName)
		if err != nil {
			return err
		}
		photoUrl = url
	}

	hangout := HangoutsTable{
		CreatedBy: t.User,
		Title:     title,
		Time:      t.Time,
		Place:     t.Place,
		Picture:   photoUrl,
		Type:      groupHangoutType,
	}
	if err := db.Table("hangouts").Create(&hangout).Error; err != nil {
		return err
	}

	now := time.Now().Format(timeFormat)

	var hangoutInvitations []HangoutsInvitationTable
	for _, username := range t.Usernames {
		hangoutInvitations = append(hangoutInvitations, HangoutsInvitationTable{
			HangoutId: hangout.Id,
			User:      t.User,
			Username:  username,
			Time:      now,
			Type:      groupHangoutType,
		})
	}

	if err := db.Table("hangouts_invitations").Create(&hangoutInvitations).Error; err != nil {
		return err
	}

	var usernamesArray []string
	for _, username := range t.Usernames {
		usernamesArray = append(usernamesArray, `'`+username+`'`)
	}

	usernamesString := strings.Join(usernamesArray, ", ")

	tokens := &[]string{}
	if err := service.GetTokensByUsernames(db, tokens, usernamesString); err != nil {
		return nil
	}
	groupHangoutNotification := service.FcmNotification{
		Sender:  t.User,
		Type:    groupHangoutType,
		Title:   t.Name + " sends a group hangout!",
		Sound:   "notification.wav",
		Devices: *tokens,
	}

	service.SendNotification(&groupHangoutNotification)
	return nil
}

// Get all hangouts from DB
func GetHangouts(db *gorm.DB, t *GetHangout) ([]Hangouts, error) {
	condition := ``
	order := ` DESC`
	if !t.ShowAll {
		today := time.Now()
		condition = ` AND T1.time > '` + today.Format("2006-01-02") + `' `

		order = ``
	}
	queryGetHangouts :=
		`	SELECT
								T1.id,
								CASE T1.created_by
								WHEN '` + t.Username + `' THEN
									T2.username
								ELSE
									T1.created_by
								END AS created_by,
								T1.title,
								T1.time,
								T1.place,
								T1.picture,
								T1.type
									FROM
										hangouts T1
										INNER JOIN hangouts_invitations T2 ON T1.id = T2.hangout_id
									WHERE (T1.created_by = '` + t.Username + `'
										OR T2.username = '` + t.Username + `')
									AND T2.confirmed = 1 ` + condition + `
									GROUP BY
										T2.hangout_id
										ORDER BY
												T1.time ` + order + `
						`
	hangouts, err := GetHangoutsFromQuery(db, queryGetHangouts)
	if err != nil {
		return nil, err
	}

	hangoutsUsers := getHangoutUsers(hangouts)

	var users []people.People
	if len(hangoutsUsers) > 1 {
		query := `SELECT username, firstname, profile_picture FROM users WHERE username IN (` + hangoutsUsers + `)`
		usersFromQuery, err := GetUsersFromQuery(db, query)
		if err != nil {
			return nil, err
		}
		users = usersFromQuery
	}

	titles := GetTitles(hangouts)

	var data []Hangouts
	for _, title := range titles {
		var hangoutsArray []HangoutsTable
		for _, hangout := range hangouts {
			if strings.Contains(hangout.Time, title) {
				if hangout.Type == hangoutType {
					for _, user := range users {
						if user.Username == hangout.CreatedBy {
							hangout.Title = user.Firstname
							hangout.Picture = user.ProfilePicture
						}
					}
				}
				hangoutsArray = append(hangoutsArray, hangout)
			}
		}

		if t.ShowAll {
			reverseHangoutsArray(hangoutsArray)
		}

		var hangoutsData []List
		hangoutsData = append(hangoutsData, List{hangoutsArray})

		data = append(data,
			Hangouts{
				Title: title,
				Data:  hangoutsData,
			},
		)
	}

	return data, nil
}

func GetUsersFromQuery(db *gorm.DB, query string) ([]people.People, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []people.People
	for rows.Next() {
		db.ScanRows(rows, &users)
	}

	return users, nil
}

func getHangoutUsers(hangouts []HangoutsTable) string {
	var usersnames []string
	for _, hangout := range hangouts {
		if hangout.Type == hangoutType {
			usersnames = append(usersnames, `'`+hangout.CreatedBy+`'`)
		}
	}

	return strings.Join(usersnames, ", ")
}

func GetHangoutsFromQuery(db *gorm.DB, query string) ([]HangoutsTable, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var hangouts []HangoutsTable
	for rows.Next() {
		db.ScanRows(rows, &hangouts)
	}

	return hangouts, nil
}

func GetTitles(hangouts []HangoutsTable) []string {
	var titles []string
	lastDate := ""
	for _, hangout := range hangouts {
		string := strings.Fields(hangout.Time)
		date := string[0]
		if date != lastDate {
			lastDate = date
			titles = append(titles, date)
		}
	}

	return titles
}

func reverseHangoutsArray(hangoutsArray []HangoutsTable) []HangoutsTable {
	for i, j := 0, len(hangoutsArray)-1; i < j; i, j = i+1, j-1 {
		hangoutsArray[i], hangoutsArray[j] = hangoutsArray[j], hangoutsArray[i]
	}
	return hangoutsArray
}

func UplaodPhoto(db *gorm.DB, username string, buffer string, fileName string) (string, error) {
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

	decode, _ := base64.StdEncoding.DecodeString(buffer)
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("notify-bucket-images"),
		Key:         aws.String("hangout-images/" + username + "/" + fileName),
		Body:        bytes.NewReader(decode),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}
