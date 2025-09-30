package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db*sql.DB


type User struct{
	ID int `json:"id"`
	Email string `json:"email`
	password string `json:"password`
}


func initDB(){
    dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")


	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	var err error
	db, err=sql.Open("postgres",connStr)
	if err!= nil{
		log.Fatalf("Error opening database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v. Please check your .env credentials and ensure Postgres is running.", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL database!")
}

func main() {

    err := godotenv.Load() 
    if err != nil {       
    log.Println("No .env file found, using system env vars")
    }

	initDB()

	port := os.Getenv("PORT")
	if port ==""{
		port="8080"
	}
    log.Printf("Server running on :%s",port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}