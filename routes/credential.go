package routes

import (
	"auth_api/handlers"
	"auth_api/middleware"

	"github.com/gin-gonic/gin"
)

func CredentialRoutes(r *gin.Engine){
	authRequired := r.Group("/api/credentials")

	authRequired.Use(middleware.AuthMiddleWare())
	{
		authRequired.POST("/",handlers.CreateCredential)
		authRequired.GET("/",handlers.GetCredentials)
	

	}
}