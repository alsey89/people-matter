package common

import (
	"fmt"
	"net/mail"
	"strconv"
	"time"
)

// validate email address
func EmailValidator(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// get current date and time in string format
func GetCurrentDateTime() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

// format time to string
func FormatTime(t time.Time) string {
	return t.Format("2006/01/02 15:04:05")
}

// take string of numbers and convert to uint
func ConvertStringOfNumbersToUint(str string) (uint, error) {
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("common.convert_string_of_numbers_to_uint: %w", err)
	}
	uintNum := uint(num)
	return uintNum, nil
}
