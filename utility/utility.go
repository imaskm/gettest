package utility

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// ConvertStringDateToTime converts date string in format (yyyy-mm-dd) in golang's time.Time format
func ConvertStringDateToTime(date string) (time.Time, error) {

	dateArr := strings.Split(date, "-")

	if len(dateArr) != 3 {
		return time.Time{}, errors.New("invalid date format")
	}
	year, err := strconv.Atoi(dateArr[0])
	if err != nil {
		return time.Time{}, errors.New("invalid date format")
	}
	month, err := strconv.Atoi(dateArr[1])
	if err != nil {
		return time.Time{}, errors.New("invalid date format")
	}
	day, err := strconv.Atoi(dateArr[2])
	if err != nil {
		return time.Time{}, errors.New("invalid date format")
	}

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return t, nil
}

// GetSum returns sum of integer slice's elements
func GetSum(arr []int) int {
	s := 0

	for _, ar := range arr {
		s += ar
	}
	return s
}
