package people

import "gorm.io/gorm"

type PeopleTable struct {
	Id       uint `gorm:"primary_key;auto_increment;not_null"`
	User     string
	Username string
}

type People struct {
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	ProfilePicture string `json:"profilePicture"`
}

func (PeopleTable) TableName() string {
	return "people"
}

// Get people from DB
func GetPeople(db *gorm.DB, t *People) ([]People, error) {
	query := `SELECT * FROM users WHERE username IN (SELECT username FROM people WHERE user = '` + t.Username + `')`

	people, err := GetPeopleFromQuery(db, query)
	if err != nil {
		return nil, err
	}

	return people, nil
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
