package middleware

import (
	"auth_api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc{

	return func(c * gin.Context){
		tokenString,err:= c.Cookie("token")
		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"to access user should login First"})
				c.Abort()
				return 
		}

		claims,err:=utils.ValidateToken(tokenString)

		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"Invalid or experied token"})
					c.Abort()
					return 
		}

		c.Set("userEmail",claims.Email)
		c.Next()
	}

}