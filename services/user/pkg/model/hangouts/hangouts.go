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
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/notifications"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

const hangoutType = "hangout"
const groupHangoutType = "group_hangout"

const timeFormat = "2006-01-02 15:04:05"

type HangoutsTable struct {
	Id               uint   `gorm:"primary_key;auto_increment;not_null" json:"id"`
	CreatedBy        string `json:"createdBy"`
	Title            string `json:"title"`
	Time             string `json:"time"`
	Place            string `json:"place"`
	Picture          string `json:"picture"`
	Type             string `json:"type"`
	CreatorConfirmed int    `gorm:"default:1" json:"creatorConfirmed"`
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
	today := time.Now().Format("2006-01-02")

	var hangoutsInvitations []HangoutsInvitationTable
	if err := db.Table("hangouts_invitations").
		Where("(user = ? OR username = ?) AND confirmed = 1", t.Username, t.Username).
		Group("hangout_id").
		Find(&hangoutsInvitations).Error; err != nil {
		return nil, err
	}

	var hangoutIds []uint
	for _, hangoutInvite := range hangoutsInvitations {
		if hangoutInvite.Type == hangoutType && hangoutInvite.Username == t.Username {
			hangoutIds = append(hangoutIds, hangoutInvite.HangoutId)
		}
	}

	var hangouts []HangoutsTable
	if err := db.Table("hangouts").
		Where("(created_by = ? OR id IN ?) AND time > ?", t.Username, hangoutIds, today).
		Group("id").
		Find(&hangouts).Error; err != nil {
		return nil, err
	}

	var groupHangoutIds []uint
	for _, hangout := range hangouts {
		if hangout.Type == groupHangoutType {
			groupHangoutIds = append(groupHangoutIds, hangout.Id)
		}
	}

	var acceptedInvitations []notifications.AcceptedInvitations
	if len(groupHangoutIds) > 0 {
		if err := db.Table("accepted_invitations").
			Where("event_id IN ? AND (type = 'accepted_group_hangout' OR type = 'accepted_hangout')", groupHangoutIds).
			Group("event_id").
			Find(&acceptedInvitations).Error; err != nil {
			return nil, err
		}
	}

	var resultHangouts []HangoutsTable
out:
	for _, hangout := range hangouts {
		// When user created group hangout and one or more users accepted
		if hangout.Type == groupHangoutType && hangout.CreatedBy == t.Username {
			for _, acceptedInvite := range acceptedInvitations {
				if acceptedInvite.EventId == hangout.Id {
					resultHangouts = append(resultHangouts, hangout)
					continue out
				}
			}
		}

		// When user created hangout and other side accepted
		if hangout.CreatedBy == t.Username && hangout.CreatorConfirmed == 1 {
			for _, hangoutInvite := range hangoutsInvitations {
				if hangout.Id == hangoutInvite.HangoutId {
					resultHangouts = append(resultHangouts, hangout)
					continue out
				}
			}
		}

		// When user confirmed and it is group hangout
		for _, hangoutInvite := range hangoutsInvitations {
			if hangout.Id == hangoutInvite.HangoutId && hangoutInvite.Username == t.Username && hangout.Type == groupHangoutType {
				resultHangouts = append(resultHangouts, hangout)
				continue out
			}
		}

		// When user and creator both confirmed
		for _, hangoutInvite := range hangoutsInvitations {
			if hangout.Id == hangoutInvite.HangoutId && hangoutInvite.Username == t.Username && hangout.CreatorConfirmed == 1 {
				resultHangouts = append(resultHangouts, hangout)
				continue out
			}
		}
	}

	hangoutsUsersString := getHangoutUsers(resultHangouts, t.Username)

	var users []people.People
	if len(hangoutsUsersString) > 0 {
		query := `SELECT username, firstname, profile_picture FROM users WHERE username IN (` + hangoutsUsersString + `)`
		usersFromQuery, err := GetUsersFromQuery(db, query)
		if err != nil {
			return nil, err
		}
		users = usersFromQuery
	}

	titles := GetTitles(resultHangouts)

	var data []Hangouts
	for _, title := range titles {
		var hangoutsArray []HangoutsTable
		for _, hangout := range resultHangouts {
			if strings.Contains(hangout.Time, title) {
				if hangout.Type == hangoutType {
					for _, user := range users {
						if user.Username != t.Username {
							hangout.Title = user.Firstname
							hangout.Picture = user.ProfilePicture
						}
					}
				}
				hangoutsArray = append(hangoutsArray, hangout)
			}
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

// Get history hangouts from DB
func GetHistoryHangouts(db *gorm.DB, t *GetHangout) ([]Hangouts, error) {
	today := time.Now().Format("2006-01-02")

	var hangouts []HangoutsTable
	if err := db.Table("hangouts T1").
		Select("T1.*").
		Joins(`JOIN hangouts_invitations T2 ON ((T1.created_by = ? AND T1.creator_confirmed = 1) OR T2.username = ?) AND T1.id = T2.hangout_id AND T2.confirmed = 1 AND T1.time < ?`, t.Username, t.Username, today).
		Group("T1.id").
		Order("T1.time DESC").
		Find(&hangouts).Error; err != nil {
		return nil, err
	}

	hangoutsUsersString := getHangoutUsers(hangouts, t.Username)

	var users []people.People
	if len(hangoutsUsersString) > 0 {
		query := `SELECT username, firstname, profile_picture FROM users WHERE username IN (` + hangoutsUsersString + `)`
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
						if user.Username != t.Username {
							hangout.Title = user.Firstname
							hangout.Picture = user.ProfilePicture
						}
					}
				}
				hangoutsArray = append(hangoutsArray, hangout)
			}
		}

		reverseHangoutsArray(hangoutsArray)

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

func getHangoutUsers(hangouts []HangoutsTable, username string) string {
	var usersnames []string
	for _, hangout := range hangouts {
		if hangout.CreatedBy != username && hangout.Type == hangoutType {
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
