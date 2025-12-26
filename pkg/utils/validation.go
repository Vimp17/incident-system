package utils

import (
	"regexp"
	"unicode"
)

// ValidateEmail проверяет валидность email
func ValidateEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}

// ValidatePhone проверяет валидность номера телефона
func ValidatePhone(phone string) bool {
    phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
    return phoneRegex.MatchString(phone)
}

// ValidatePassword проверяет сложность пароля
func ValidatePassword(password string) bool {
    if len(password) < 8 {
        return false
    }
    
    var (
        hasUpper   bool
        hasLower   bool
        hasNumber  bool
        hasSpecial bool
    )
    
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsDigit(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    return hasUpper && hasLower && hasNumber && hasSpecial
}

// ValidateCoordinates проверяет валидность координат
func ValidateCoordinates(lat, lng float64) bool {
    return lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180
}