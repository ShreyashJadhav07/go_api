package routes

import (
	"auth_api/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	api:=r.Group("/api")

	{
		api.POST("/signup",handlers.SignUp)

	}

}