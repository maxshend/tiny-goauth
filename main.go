package main

import (
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/maxshend/tiny_goauth/models"
)

// HandlerDeps contains dependencies of the http handlers
type HandlerDeps struct {
	db *pgxpool.Pool
}

func main() {
	db, err := models.InitDB(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
