package common

import (
	"fmt"
	"regexp"
)

var regexpEmail = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$")

func ValidateEmail(email string) error {
	if !regexpEmail.MatchString(email) {
		return fmt.Errorf("invalid email")
	}
	return nil
}

var reASCII = regexp.MustCompile(`^[\x20-\x7E]+$`)

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	if !reASCII.MatchString(password) {
		return fmt.Errorf("only ASCII characters allowed")
	}
	return nil
}
