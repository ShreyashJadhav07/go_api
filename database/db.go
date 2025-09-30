
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

var DB *sql.DB


func InitDB() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}

	
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")


	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)


	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}


	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v. Please check your .env credentials and ensure Postgres is running.", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL database!")
}


func GetDB() *sql.DB {
    return DB
}
