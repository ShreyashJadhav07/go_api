package main

import (

	"auth_api/database"
	"auth_api/routes"
	"log"

	"os"

	"github.com/gin-gonic/gin"
	
)


func main() {

	database.InitDB()

	router:=gin.Default()

	routes.AuthRoutes(router)

   

	port := os.Getenv("PORT")
	if port ==""{
		port="8080"
	}
    log.Printf("Server running on :%s",port)
	log.Fatal(router.Run(":" + port))

}