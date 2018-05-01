package presentation

import (
	"net/mail"
	"unicode"
)

// Password ...
type Password struct {
	isValidLenght             bool
	isValidNumberCharacter    bool
	isValidUpperCaseCharacter bool
	isValidLowerCaseCharacter bool
	isValidSpecialCharacter   bool
}

// IsValidEmail ...
func IsValidEmail(emailAddress string) bool {
	_, err := mail.ParseAddress(emailAddress)
	return err != nil
}

// NewPassword ..
func NewPassword(
	password string,
	minLenght int,
	minNumberCharacter int,
	minUpperCase int,
	minLowerCase int,
	minSpecialCharacter int,
) *Password {
	numberChar := 0
	upperCase := 0
	lowerCase := 0
	specialCharacter := 0
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			numberChar++
		case unicode.IsLower(c):
			lowerCase++
		case unicode.IsUpper(c):
			upperCase++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			specialCharacter++
		}
	}
	return &Password{
		isValidLenght:             len(password) >= minLenght,
		isValidLowerCaseCharacter: lowerCase >= minLowerCase,
		isValidUpperCaseCharacter: upperCase >= minUpperCase,
		isValidSpecialCharacter:   specialCharacter >= minSpecialCharacter,
		isValidNumberCharacter:    numberChar >= minNumberCharacter,
	}
}
