package utils

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

func IntLength(number int64) int {
	if number == 0 {
		return 1
	}

	length := 0

	for number != 0 {
		number /= 10
		length++
	}

	return length
}

func IsValidDateFormat(date string) bool {
	var dateRegex = regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])\.(0[1-9]|1[0-2])\.(\d{4})$`)

	matches := dateRegex.FindStringSubmatch(date)

	if matches == nil {
		return false
	}

	year, err := strconv.Atoi(matches[3])
	if err != nil {
		return false
	}

	return year <= time.Now().Year()
}

func IsValidEmail(email string) error {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if email == "" {
		return errors.New("invalid input: email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("invalid input: email must contain a valid domain")
	}

	return nil
}
