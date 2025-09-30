package models

import (
	"fmt"
	"regexp"
	"strings"
)

type User struct {
	Email           string `json:"email" binding:"required,email,max=255"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=72"`
}

func (u *User) Validate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	if len(u.Email) > 255 {
		return fmt.Errorf("email too long")
	}

	if err:= ValidatePasswordStrength(u.Password); err !=nil{
		return err
	}

	return nil
}

func ValidatePasswordStrength(password string) error{
	if len(password) < 8 {
		return fmt.Errorf("Password must be at least 8 Characters")

	}
	if len(password) > 72{
		return fmt.Errorf("Password must be at least 72 characters")
	}

  var (
        hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
        hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
        hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
        hasSpecial = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
    )

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
        return fmt.Errorf("password must contain uppercase, lowercase, number, and special character")
    }
    
    return nil
}
