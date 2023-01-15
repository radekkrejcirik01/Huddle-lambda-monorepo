package hangouts

import (
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type HangoutsTable struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	Username string
	Time     string
	Place    string
}

type HangoutInvite struct {
	User     string
	Username string
	Time     string
	Place    string
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
	List []Hangout `json:"list"`
}

type Hangout struct {
	Id        uint   `json:"id"`
	CreatedBy User   `json:"createdBy"`
	Users     []User `json:"users"`
	Time      string `json:"time"`
	Place     string `json:"place"`
}

type User struct {
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	ProfilePicture string `json:"profilePicture"`
}

func (HangoutsTable) TableName() string {
	return "hangouts"
}

// Create new hangout in DB
func CreateHangout(db *gorm.DB, t *HangoutInvite) error {
	hangout := HangoutsTable{Username: t.User, Time: t.Time, Place: t.Place}
	if err := db.Create(&hangout).Error; err != nil {
		return err
	}

	hangoutInvitation := HangoutsInvitationTable{
		HangoutId: hangout.Id,
		Username:  t.Username,
		Confirmed: 0,
	}
	return db.Table("hangouts_invitations").Create(&hangoutInvitation).Error
}

// Get all hangouts from DB
func GetHangouts(db *gorm.DB, t *GetHangout) ([]Hangouts, error) {
	condition := ``
	order := ` DESC`
	if !t.ShowAll {
		today := time.Now()
		condition = ` AND time > '` + today.Format("2006-01-02") + `' `

		order = ``
	}
	queryGetAllHangouts :=
		`
						SELECT
							*
						FROM
							hangouts
						WHERE
							(username = '` + t.Username + `'
							OR id IN(
								SELECT
									hangout_id FROM hangouts_invitations
								WHERE
									username = '` + t.Username + `'
									AND confirmed = 1)) ` + condition + `
						ORDER BY
							time` + order + `
		`
	allHangouts, err := GetAllHangoutsFromQuery(db, queryGetAllHangouts)
	if err != nil {
		return nil, err
	}

	idsArray := getIdsArray(allHangouts)

	queryGetUserAndConfirms := `SELECT * FROM hangouts_invitations WHERE hangout_id IN (` + idsArray + `) AND username != '` + t.Username + `'`
	usersAndConfirms, err := GetUserAndConfirmsFromQuery(db, queryGetUserAndConfirms)
	if err != nil {
		return nil, err
	}

	titles := GetTitles(allHangouts)

	var data []Hangouts
	for _, title := range titles {
		var hangouts []Hangout
		for _, hangout := range allHangouts {
			if strings.Contains(hangout.Time, title) {
				if !hasConfirmed(hangout, usersAndConfirms) && hangout.Username == t.Username {
					continue
				}

				usersByHangout := GetUsersByHangout(db, hangout, usersAndConfirms)
				createdBy := User{}
				var users []User
				for _, user := range usersByHangout {
					if user.Username != hangout.Username {
						users = append(users, user)
					} else {
						createdBy = user
					}
				}

				hangouts = append(hangouts, Hangout{
					Id:        hangout.Id,
					CreatedBy: createdBy,
					Users:     users,
					Time:      GetTime(hangout.Time),
					Place:     hangout.Place,
				})
			}
		}

		var hangoutsData []List
		if len(hangouts) > 0 {
			hangoutsData = append(hangoutsData, List{hangouts})

			data = append(data,
				Hangouts{
					Title: title,
					Data:  hangoutsData,
				},
			)
		}
	}

	return data, nil
}

func GetTime(datetime string) string {
	string := strings.Fields(datetime)
	time := string[1]
	return time[0:5]
}

func hasConfirmed(hangout HangoutsTable, usersAndConfirms []HangoutsInvitationTable) bool {
	confirmed := false
	for _, user := range usersAndConfirms {
		if user.HangoutId == hangout.Id {
			if user.Confirmed == 1 {
				confirmed = true
			}
		}
	}

	return confirmed
}

func GetUsersByHangout(db *gorm.DB, hangout HangoutsTable, usersAndConfirms []HangoutsInvitationTable) []User {
	var usernames []string
	usernames = append(usernames, "'"+hangout.Username+"'")
	for _, user := range usersAndConfirms {
		if user.HangoutId == hangout.Id {
			usernames = append(usernames, "'"+user.Username+"'")
		}
	}

	usernamesArray := strings.Join(usernames, ", ")

	query := `SELECT username, firstname, profile_picture FROM users WHERE username IN (` + usernamesArray + `)`

	users, err := GetUsersFromQuery(db, query)
	if err != nil {
		return []User{}
	}
	return users
}

func GetUsersFromQuery(db *gorm.DB, query string) ([]User, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		db.ScanRows(rows, &users)
	}

	return users, nil
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

func GetAllHangoutsFromQuery(db *gorm.DB, query string) ([]HangoutsTable, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var allHangouts []HangoutsTable
	for rows.Next() {
		db.ScanRows(rows, &allHangouts)
	}

	return allHangouts, nil
}

func getIdsArray(hangouts []HangoutsTable) string {
	var ids []string
	for _, hangout := range hangouts {
		ids = append(ids, strconv.Itoa(int(hangout.Id)))
	}
	idsArray := strings.Join(ids, ", ")

	return idsArray
}

func GetUserAndConfirmsFromQuery(db *gorm.DB, query string) ([]HangoutsInvitationTable, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var usersAndConfirms []HangoutsInvitationTable
	for rows.Next() {
		db.ScanRows(rows, &usersAndConfirms)
	}

	return usersAndConfirms, nil
}
