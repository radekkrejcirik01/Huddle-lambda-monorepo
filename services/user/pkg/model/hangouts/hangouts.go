package hangouts

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type HangoutsTable struct {
	Id        uint   `gorm:"primary_key;auto_increment;not_null" json:"id"`
	CreatedBy string `json:"createdBy"`
	Title     string `json:"title"`
	Time      string `json:"time"`
	Place     string `json:"place"`
	Picture   string `json:"picture"`
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
	List []HangoutsTable `json:"list"`
}

func (HangoutsTable) TableName() string {
	return "hangouts"
}

// Create new hangout in DB
func CreateHangout(db *gorm.DB, t *HangoutInvite) error {
	hangout := HangoutsTable{CreatedBy: t.User, Time: t.Time, Place: t.Place}
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
		condition = ` AND T1.time > '` + today.Format("2006-01-02") + `' `

		order = ``
	}
	queryGetHangouts :=
		`
							SELECT
								T1.*
							FROM
							hangouts T1
								INNER JOIN hangouts_invitations T2 ON T1.id = T2.hangout_id
							WHERE (T1.created_by = '` + t.Username + `'
								OR T2.username = '` + t.Username + `')
							AND T2.confirmed = 1` + condition + `
							GROUP BY
								id
							ORDER BY
								T1.time ` + order + `
		`
	hangouts, err := GetHangoutsFromQuery(db, queryGetHangouts)
	if err != nil {
		return nil, err
	}

	titles := GetTitles(hangouts)

	var data []Hangouts
	for _, title := range titles {
		var hangoutsArray []HangoutsTable
		for _, hangout := range hangouts {
			if strings.Contains(hangout.Time, title) {
				hangout.Time = GetTime(hangout.Time)
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

func GetTime(datetime string) string {
	string := strings.Fields(datetime)
	time := string[1]
	return time[0:5]
}

func reverseHangoutsArray(hangoutsArray []HangoutsTable) []HangoutsTable {
	for i, j := 0, len(hangoutsArray)-1; i < j; i, j = i+1, j-1 {
		hangoutsArray[i], hangoutsArray[j] = hangoutsArray[j], hangoutsArray[i]
	}
	return hangoutsArray
}
