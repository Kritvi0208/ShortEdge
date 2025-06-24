package factory

import (
	"database/sql"
	"log"
	"os"
	"urlify/store"

	//"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gofr.dev/pkg/gofr"
)

var db *sql.DB

func NewURLStore(app *gofr.App) store.URL {
	return store.NewURLStore(GetDB())
}

func NewVisitStore(app *gofr.App) store.Visit {
	return store.NewVisitStore(GetDB())
}

func GetDB() *sql.DB {
	if db != nil {
		return db
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("❌ DB_URL is not set. Please set it in environment or .env file.")
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Could not ping DB: %v", err)
	}
 return db 
}
