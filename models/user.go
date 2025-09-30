package models

type User struct{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}