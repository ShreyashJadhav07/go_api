// This file contains the API handlers for authentication (signup, login, etc.).
package handlers

import (
	"log"
	"net/http"

	// Using relative path for internal packages
	"auth_api/database"
	"auth_api/models" 

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SignUp handles user registration. It receives JSON, validates it, hashes the password, and inserts the user into the database.
func SignUp(c *gin.Context) {
	// 1. Bind JSON input to User model for validation (Gin handles validation via 'binding' tags).
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		// Respond with bad request if JSON is invalid or fails model validation
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format. Check email and ensure password is at least 8 characters."})
		return
	}

	// 2. Hash the password securely (Industry Best Practice).
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password securely."})
		return
	}

	// 3. Prepare the secure SQL INSERT statement.
	// We use QueryRow to insert the data and immediately retrieve the automatically generated ID.
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`

	// Get the database connection
	db := database.GetDB()

	// 4. Execute the statement.
	var newID int
	err = db.QueryRow(query, user.Email, hashedPassword).Scan(&newID)
	
	if err != nil {
		// Handle database errors, particularly duplicate email (Unique Constraint Violation)
		log.Printf("Database insertion error: %v", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered. Please use a different email."})
		return
	}

	// 5. Success response.
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"id":      newID,
		"email":   user.Email,
	})
}