package handlers

import (
	"auth_api/database"
	"auth_api/models"
	"log"
	"net/http"
	

	"github.com/gin-gonic/gin"
)


func getUserIDFromContext(c *gin.Context) (int ,bool){
	userIDVal , exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized,gin.H{"error":"User Session invalid."})
		return 0, false
	}
	userId , ok :=userIDVal.(int)

	if !ok {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error":"Failed to get User ID."})
			return 0, false
		
	}
	return  userId,true

}

func CreateCredential(c *gin.Context){
	var req models.CredentialRequest
	if err := c.ShouldBindJSON(&req); err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid input format or missing fields."})
		return
	}

	userID ,ok := getUserIDFromContext(c)

	if !ok {
		return
	}

    query := `INSERT INTO credentials (user_id, service_name, url, username, password_cipher, nonce, notes) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) 
			  RETURNING id`

	db:=database.GetDB()
	var newID int
	err:= db.QueryRow(query,
		userID,
		req.ServiceName,
		req.URL,
		req.Username,
		req.PasswordCipher,
		req.Nonce,
		req.Notes).
		Scan(&newID)

	if err!=nil{
		log.Printf("Database insertion error: %v", err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": "failed to save credential."})
			return
	}

	c.JSON(http.StatusCreated,gin.H{
		"message":"Credential stored successfully",
		"id": newID,
	})

	

}

func GetCredentials(c *gin.Context){
	userID, ok :=getUserIDFromContext(c)
	if !ok {
		return
	}
	query := `SELECT id, service_name, url, username, notes, created_at, updated_at 
	          FROM credentials 
			  WHERE user_id = $1 
			  ORDER BY service_name`
		
	db := database.GetDB()
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Database query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve credentials."})
		return
	}

		var credentials []models.CredentialResponse
	for rows.Next() {
		var cred models.CredentialResponse
		err := rows.Scan(&cred.ID, &cred.ServiceName, &cred.URL, &cred.Username, &cred.Notes, &cred.CreatedAt, &cred.UpdatedAt)
		if err != nil {
			log.Printf("Scanning row error: %v", err)

			continue 
		}
		credentials = append(credentials, cred)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing results."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Credentials fetched successfully",
		"count": len(credentials),
		"data": credentials,
	})
}


