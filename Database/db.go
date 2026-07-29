package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func DatabaseConnection() (*sql.DB, error) {

	godotenv.Load(".env")
	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("Failed to load Url")
	}

	db, err := sql.Open("pgx", db_url)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to connect database %v", err)
	}
	
	db.SetMaxOpenConns(25)
	db.SetConnMaxIdleTime(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("Database connection failed")
	}
	fmt.Println("Database Connected successfully")
	return db, nil
}
