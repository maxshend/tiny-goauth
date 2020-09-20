package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maxshend/tiny_goauth/db"
	"github.com/maxshend/tiny_goauth/handlers"
	"github.com/maxshend/tiny_goauth/validations"
)

func main() {
	db, err := db.Init(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	validator, translator, err := validations.Init()
	if err != nil {
		log.Fatal(err)
	}

	deps := &handlers.Deps{DB: db, Validator: validator, Translator: translator}
	server := http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		Handler:      nil,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	http.HandleFunc("/email/register", handlers.EmailRegister(deps))

	log.Fatal(server.ListenAndServe())
}
