package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maxshend/tiny_goauth/db"
	"github.com/maxshend/tiny_goauth/handlers"
)

func main() {
	db, err := db.Init(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	deps := &handlers.Deps{DB: db}
	server := http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		Handler:      nil,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	http.HandleFunc("/email/register", handlers.EmailRegister(deps))

	log.Fatal(server.ListenAndServe())
}
