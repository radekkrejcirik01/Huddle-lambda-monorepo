package hangouts

import (
	"sort"
	"strconv"
	"strings"

	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"gorm.io/gorm"
)

type HangoutId struct {
	Id       uint
	Username string
}

type HangoutById struct {
	CreatedBy        string      `json:"createdBy"`
	Title            string      `json:"title"`
	Time             string      `json:"time"`
	Place            string      `json:"place"`
	Picture          string      `json:"picture"`
	Usernames        []Usernames `json:"usernames"`
	Type             string      `json:"type"`
	CreatorConfirmed int         `json:"creatorConfirmed"`
}

type UpdateHangout struct {
	Id       uint
	Buffer   string
	FileName string
	Title    string
	Time     string
	Plan     string
}

type HangoutUsernames struct {
	Username  string
	Confirmed int
}

type Usernames struct {
	Username       string `json:"username"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profilePicture"`
	Confirmed      int    `json:"confirmed"`
}

type User struct {
	Firstname      string
	ProfilePicture string
}

// Get hangout by id from DB
func GetHangoutById(db *gorm.DB, t *HangoutId) (HangoutById, error) {
	var hangout HangoutsTable
	if err := db.Table("hangouts").Where("id = ?", t.Id).First(&hangout).Error; err != nil {
		return HangoutById{}, err
	}

	var hangoutUsernames []HangoutUsernames
	if err := db.Table("hangouts_invitations").Select("username, confirmed").Where("hangout_id = ?", t.Id).Find(&hangoutUsernames).Error; err != nil {
		return HangoutById{}, err
	}
	hangoutUsernames = append(hangoutUsernames, HangoutUsernames{
		Username:  hangout.CreatedBy,
		Confirmed: hangout.CreatorConfirmed,
	})

	title := hangout.Title
	picture := hangout.Picture
	if hangout.Type == hangoutType {
		hangoutUsername := hangoutUsernames[0].Username
		if hangoutUsername == t.Username {
			hangoutUsername = hangoutUsernames[1].Username
		}

		var user User
		if err := db.Table("users").Select("firstname, profile_picture").Where("username = ?", hangoutUsername).First(&user).Error; err != nil {
			return HangoutById{}, err
		}

		title = user.Firstname
		picture = user.ProfilePicture
	}

	var usernamesArray []string
	for _, user := range hangoutUsernames {
		usernamesArray = append(usernamesArray, `'`+user.Username+`'`)
	}
	usernamesString := strings.Join(usernamesArray, ", ")

	var users []people.People
	if err := db.Table("users").Select("username, firstname, profile_picture").Where("username IN (" + usernamesString + ")").Find(&users).Error; err != nil {
		return HangoutById{}, err
	}

	usernames := []Usernames{}
	for _, hangoutUsername := range hangoutUsernames {
		for _, user := range users {
			if user.Username == hangoutUsername.Username {
				usernames = append(usernames, Usernames{
					Username:       hangoutUsername.Username,
					Name:           user.Firstname,
					ProfilePicture: user.ProfilePicture,
					Confirmed:      hangoutUsername.Confirmed,
				})
			}
		}
	}

	sort.SliceStable(usernames, func(i, j int) bool {
		return usernames[i].Name < usernames[j].Name
	})

	result := HangoutById{
		CreatedBy:        hangout.CreatedBy,
		Title:            title,
		Time:             hangout.Time,
		Place:            hangout.Place,
		Picture:          picture,
		Usernames:        usernames,
		Type:             hangout.Type,
		CreatorConfirmed: hangout.CreatorConfirmed,
	}

	return result, nil
}

// Update hangout in DB
func UpdateHangoutById(db *gorm.DB, t *UpdateHangout) error {
	update := map[string]interface{}{}

	var photoUrl string
	if len(t.Buffer) > 0 {
		url, err := UplaodPhoto(db, strconv.Itoa(int(t.Id)), t.Buffer, t.FileName)
		if err != nil {
			return err
		}
		photoUrl = url
	}

	if len(photoUrl) > 0 {
		update["picture"] = photoUrl
	}
	if len(t.Title) > 0 {
		update["title"] = t.Title
	}
	if len(t.Time) > 0 {
		update["time"] = t.Time
	}
	if len(t.Plan) > 0 {
		update["place"] = t.Plan
	}

	return db.Table("hangouts").Where("id = ?", t.Id).Updates(update).Error
}

// Remove user from hangout in DB
func RemoveUserFromHangout(db *gorm.DB, t *HangoutId) error {
	return db.Exec(`
	DELETE T1,
	T2 FROM hangouts_invitations T1
		LEFT JOIN accepted_invitations T2 ON T1.hangout_id = T2.event_id
		WHERE (T1.hangout_id = ?
				AND T1.username = ?)
			OR(T2.event_id = ?
				AND T2.user = ? AND T2.type = 'accepted_hangout')`, t.Id, t.Username, t.Id, t.Username).Error
}

// Delete hangout by id from DB
func DeleteHangoutById(db *gorm.DB, t *HangoutId) error {
	return db.Exec(`
		DELETE hangouts, hangouts_invitations, accepted_invitations FROM hangouts
			LEFT JOIN hangouts_invitations ON hangouts.id = hangouts_invitations.hangout_id
			LEFT JOIN accepted_invitations ON hangouts.id = accepted_invitations.event_id
		WHERE hangouts.id = ?`, t.Id).Error
}
