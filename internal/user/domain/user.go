package domain

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

type (
	UserID       uint
	Phone        string
	Email        string
	NationalCode string
)
type UserFilter struct {
	ID    UserID
	Phone string
}

func (p Phone) IsValid() bool {
	// Define the regex pattern for Iranian mobile numbers
	// Iranian mobile numbers start with +98 or 0, followed by 9 and 9 more digits
	regexPattern := `^(?:\+98|0)?9\d{9}$`

	// Compile the regex
	regex := regexp.MustCompile(regexPattern)

	// Validate the phone number using the regex
	return regex.MatchString(string(p))
}
func (e Email) IsValid() bool {
	// Define a regex pattern for validating email addresses
	// This pattern checks for a basic structure of local@domain
	regexPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	regex := regexp.MustCompile(regexPattern)

	// Validate the email using the regex
	return regex.MatchString(string(e))
}
func (nc NationalCode) IsValid() bool {
	// Check if the national code is a 10-digit number
	regex := regexp.MustCompile(`^\d{10}$`)
	if !regex.MatchString(string(nc)) {
		return false
	}

	// Check if all digits are the same (e.g., "1111111111")
	if isAllSameDigits(string(nc)) {
		return false
	}

	// Perform checksum validation
	return validateChecksum(string(nc))
}

// Helper function to check if all digits are the same
func isAllSameDigits(nc string) bool {
	for i := 1; i < len(nc); i++ {
		if nc[i] != nc[0] {
			return false
		}
	}
	return true
}

// Helper function to validate checksum
func validateChecksum(nc string) bool {
	// Convert the first 9 digits to integers
	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(nc[i]))
		sum += digit * (10 - i)
	}

	// Calculate the remainder
	remainder := sum % 11

	// Get the control digit (the last digit)
	controlDigit, _ := strconv.Atoi(string(nc[9]))

	// Validate based on the rules
	return (remainder < 2 && remainder == controlDigit) || (remainder >= 2 && controlDigit == 11-remainder)
}

type User struct {
	ID                UserID
	FirstName         string
	LastName          string
	Phone             Phone
	Email             Email
	PasswordHash      string
	NationalCode      NationalCode
	BirthDate         time.Time
	City              string
	Gender            bool
	SurveyLimitNumber int
	CreatedAt         time.Time
	DeletedAt         time.Time
	UpdatedAt         time.Time
	Balance           int
}

func (u *User) Validate() error {
	if !u.Phone.IsValid() {
		return errors.New("phone is not valid")
	}
	if !u.Email.IsValid() {
		return errors.New("email is not valid")
	}
	if !u.NationalCode.IsValid() {
		return errors.New("national code is not valid")
	}
	return nil
}
