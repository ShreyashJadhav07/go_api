// package middleware

// import (
// 	"auth_api/utils"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func AuthMiddleWare() gin.HandlerFunc{

// 	return func(c * gin.Context){
// 		tokenString,err:= c.Cookie("token")
// 		if err!=nil{
// 			c.JSON(http.StatusUnauthorized,gin.H{
// 				"error":"to access user should login First"})
// 				c.Abort()
// 				return 
// 		}

// 		claims,err:=utils.ValidateToken(tokenString)

// 		if err!=nil{
// 			c.JSON(http.StatusUnauthorized,gin.H{
// 				"error":"Invalid or experied token"})
// 					c.Abort()
// 					return 
// 		}

// 		c.Set("userEmail",claims.Email)
// 		c.Next()
// 	}

// }

package middleware

import (
	"auth_api/utils"
	"auth_api/database" // <-- NEW IMPORT
	"net/http"
	"log"              // <-- NEW IMPORT
	"strings"          // <-- NEW IMPORT

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		var err error

		// 1. Check Authorization Header for "Bearer <token>" (Supports Postman)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Fallback to checking for "token" cookie
		if tokenString == "" {
			tokenString, err = c.Cookie("token")
			if err != nil {
				// If token is missing from both sources, abort
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required (missing token)."})
				c.Abort()
				return
			}
		}

		// 3. Validate Token
		claims, err := utils.ValidateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 4. Fetch User ID from DB (CRITICAL for credential ownership)
		db := database.GetDB()
		var userID int
		query := "SELECT id FROM users WHERE email=$1 LIMIT 1"

		err = db.QueryRow(query, claims.Email).Scan(&userID)
		if err != nil {
			log.Printf("Error fetching user ID for email %s: %v", claims.Email, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User record not found."})
			c.Abort()
			return
		}
		
		// 5. Set Context values
		c.Set("userEmail", claims.Email)
		c.Set("userID", userID) // <-- CRITICAL: Now set for handlers
		c.Next()
	}
}