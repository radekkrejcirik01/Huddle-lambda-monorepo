package hangouts

import (
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
	if err := db.Table("hangouts_invitations").Select("username").Where("hangout_id = ?", t.Id).Find(&usernames).Error; err != nil {
		return HangoutById{}, err
	}
	usernames = append(usernames, hangout.CreatedBy)

	title := hangout.Title
	picture := hangout.Picture
	if hangout.Type == hangoutType {
		username := usernames[0]
		if username == t.Username {
			username = usernames[1]
		}

		var user User
		if err := db.Table("users").Select("firstname, profile_picture").Where("username = ?", username).First(&user).Error; err != nil {
			return HangoutById{}, err
		}

		title = user.Firstname
		picture = user.ProfilePicture
	}

	result := HangoutById{
		CreatedBy: hangout.CreatedBy,
		Title:     title,
		Time:      hangout.Time,
		Place:     hangout.Place,
		Picture:   picture,
		Usernames: usernames,
	}

	return result, nil
}
