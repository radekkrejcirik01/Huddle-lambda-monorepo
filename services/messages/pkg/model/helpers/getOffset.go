package helpers

import "strconv"

func GetOffset(page string) string {
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return "0"
	}
	offset := pageNumber * 10

	return strconv.Itoa(offset)
}
