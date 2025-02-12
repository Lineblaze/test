package domain

import "regexp"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func ValidateUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._-]{3,32}$`)
	return re.MatchString(username)
}

func ValidatePassword(password string) bool {
	re := regexp.MustCompile(`^[A-Za-z\d!@#$%^&*()\-_+=]{8,}$`)
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	return re.MatchString(password) && hasLetter && hasDigit
}
