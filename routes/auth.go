package routes

import (
	"auth_api/handlers"
	"auth_api/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	api:=r.Group("/api")

	{
		api.POST("/signup",handlers.SignUp)
		api.POST("/login",handlers.Login)
	    api.GET("/profile",middleware.AuthMiddleWare(),handlers.Profile)
		api.POST("/logout",handlers.Logout)
	

	}



}