package handlers

import (
	
	"log"
	"net/http"

	"auth_api/database"
	"auth_api/models"
	"auth_api/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


func SignUp(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
	
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format. Check email and ensure password is at least 8 characters."})
		return
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password securely."})
		return
	}


	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`


	db := database.GetDB()

	
	var newID int
	err = db.QueryRow(query, user.Email, hashedPassword).Scan(&newID)
	
	if err != nil {
		// Handle database errors, particularly duplicate email (Unique Constraint Violation)
		log.Printf("Database insertion error: %v", err)
		// For a duplicate email, PostgreSQL returns a specific error we catch and report as a 409 Conflict.
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered. Please use a different email."})
		return
	}


	tokenString,err :=utils.GenerateToken(user.Email)
	if err !=nil{
		log.Printf("Error Generating JWT on signup: %V",err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to create session token after registration."})
			return
	}
	c.SetCookie("token", tokenString, 60*60*24, "/", "localhost" ,false ,true)
	

	// 5. Success response.
	// NOTE: We are intentionally using StatusCreated (201) here, not 200 OK.
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"id":      newID,
		"email":   user.Email,
	})
}

	
