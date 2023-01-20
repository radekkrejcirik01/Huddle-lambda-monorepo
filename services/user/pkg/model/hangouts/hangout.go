package hangouts

import (
	"strings"

	"gorm.io/gorm"
)

type HangoutId struct {
	Id       uint
	Username string
}

type HangoutById struct {
	CreatedBy string   `json:"createdBy"`
	Title     string   `json:"title"`
	Time      string   `json:"time"`
	Place     string   `json:"place"`
	Picture   string   `json:"picture"`
	Usernames []string `json:"usernames"`
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

	var usernames []string
	if hangout.Type == hangoutType && hangout.CreatedBy != t.Username {
		usernames = append(usernames, hangout.CreatedBy)
	} else {
		if err := db.Table("hangouts_invitations").Select("username").Where("hangout_id = ? AND username != ?", t.Id, t.Username).Find(&usernames).Error; err != nil {
			return HangoutById{}, err
		}
	}

	title := hangout.Title
	picture := hangout.Picture
	if hangout.Type == hangoutType {
		var user User
		if err := db.Table("users").Select("firstname, profile_picture").Where("username = ?", usernames[0]).First(&user).Error; err != nil {
			return HangoutById{}, err
		}

		title = user.Firstname
		picture = user.ProfilePicture
	}

	string := strings.Fields(hangout.Time)
	time := string[1]

	result := HangoutById{
		CreatedBy: hangout.CreatedBy,
		Title:     title,
		Time:      time,
		Place:     hangout.Place,
		Picture:   picture,
		Usernames: usernames,
	}

	return result, nil
}
