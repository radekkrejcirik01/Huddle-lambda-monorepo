package people

import (
	"time"

	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/notifications"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/service"
	"gorm.io/gorm"
)

const timeFormat = "2006-01-02 15:04:05"

type PeopleInvitationTable struct {
	Id        uint `gorm:"primary_key;auto_increment;not_null"`
	User      string
	Username  string
	Time      string
	Confirmed int `gorm:"default:0"`
	Seen      int `gorm:"default:0"`
}

type People struct {
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	ProfilePicture string `json:"profilePicture"`
}

type AcceptInvite struct {
	Id       uint
	Value    int
	User     string
	Username string
	Name     string
}
type Notification struct {
	Sender  string
	Title   string
	Body    string
	Sound   string
	Devices []string
}

type CheckIfFriend struct {
	Id        uint   `json:"id"`
	User      string `json:"user"`
	Username  string `json:"username"`
	Confirmed uint   `json:"confirmed"`
}

func (PeopleInvitationTable) TableName() string {
	return "people_invitations"
}

// Create new user invitatiom in DB
func CreatePeopleInvitation(db *gorm.DB, t *PeopleInvitationTable) (string, error) {
	errMessage := "We apologize, we couldn't send invite 😔"

	var exists bool
	if err := db.Table("users").Select("count(*) > 0").Where("username = ?", t.Username).Find(&exists).Error; err != nil {
		return errMessage, err
	}

	if exists {
		var user string
		if err := db.Table("people_invitations").Select("user").Where("(user = ? AND username = ?) OR (user = ? AND username = ?)", t.User, t.Username, t.Username, t.User).Find(&user).Error; err != nil {
			return errMessage, err
		}

		if user == t.User {
			return "Sorry, you already invite this user 🙂", nil
		}
		if user == t.Username {
			return "This user already invited you. Please check your notifications 🙂", nil
		}

		time := time.Now()
		t.Time = time.Format(timeFormat)

		if err := db.Create(t).Error; err != nil {
			return errMessage, err
		}

		tokens := &[]string{}
		if err := service.GetTokensByUsername(db, tokens, t.Username); err != nil {
			return "", nil
		}
		friendInviteNotification := service.FcmNotification{
			Sender:  t.User,
			Type:    "people",
			Body:    t.User + " sends friend invite",
			Sound:   "default",
			Devices: *tokens,
		}
		service.SendNotification(&friendInviteNotification)

		return "Invitation sent! ✅", nil
	}
	return "We apologize, this user doesn't exist 😔", nil
}

// Get people from DB
func GetPeople(db *gorm.DB, t *People) ([]People, error) {
	query := `	SELECT
					*
				FROM
					users
				WHERE username IN (
							SELECT
								CASE T1.username
								WHEN '` + t.Username + `' THEN
									T1.user
								ELSE
									T2.username
								END AS username
							FROM
								people_invitations T1
								INNER JOIN people_invitations T2 ON T1.id = T2.id
							WHERE (T1.username = '` + t.Username + `'
								AND T1.confirmed = 1)
								OR(T2.user = '` + t.Username + `'
									AND T2.confirmed = 1))`

	people, err := GetPeopleFromQuery(db, query)
	if err != nil {
		return nil, err
	}

	return people, nil
}

func AcceptInvitation(db *gorm.DB, t *AcceptInvite) error {
	if err := db.Table("people_invitations").Where("id = ?", t.Id).Updates(map[string]interface{}{"confirmed": t.Value, "seen": 1}).Error; err != nil {
		return err
	}

	now := time.Now().Format(timeFormat)
	acceptedInvitation := notifications.AcceptedInvitations{
		EventId:  t.Id,
		User:     t.User,
		Username: t.Username,
		Time:     now,
		Type:     "accepted_people",
	}

	if rowsAffected := db.Table("accepted_invitations").Where(notifications.AcceptedInvitations{EventId: t.Id}).FirstOrCreate(&acceptedInvitation).RowsAffected; rowsAffected == 0 {
		return nil
	}

	tokens := &[]string{}
	if err := service.GetTokensByUsername(db, tokens, t.Username); err != nil {
		return nil
	}
	acceptFriendInviteNotification := service.FcmNotification{
		Sender:  t.User,
		Type:    "people",
		Body:    t.Name + " accepted friend invite!",
		Sound:   "default",
		Devices: *tokens,
	}
	service.SendNotification(&acceptFriendInviteNotification)

	return nil
}

func GetPeopleFromQuery(db *gorm.DB, query string) ([]People, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var people []People
	for rows.Next() {
		db.ScanRows(rows, &people)
	}

	return people, nil
}

func CheckInvitations(db *gorm.DB, t *CheckIfFriend) (CheckIfFriend, error) {
	var record CheckIfFriend
	err := db.Table("people_invitations").
		Select("id, user, username, confirmed").
		Where("(user = ? AND username = ?) OR (user = ? AND username = ?)",
			t.User, t.Username, t.Username, t.User).
		First(&record).Error
	return record, err
}
