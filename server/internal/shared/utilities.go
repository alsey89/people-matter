package shared

import (
	"net/mail"
	"time"
)

func EmailValidator(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GetCurrentDateTime() string {
	return time.Now().Format("2006/01/02 15:04:05")
}
