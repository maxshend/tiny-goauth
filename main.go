package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/maxshend/tiny_goauth/models"
)

// HandlerDeps contains dependencies of the http handlers
type HandlerDeps struct {
	DB *pgxpool.Pool
}

func main() {
	db, err := models.InitDB(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	deps := &HandlerDeps{DB: db}
	server := http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		Handler:      nil,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	http.HandleFunc("/email/register", emailRegister(deps))

	log.Fatal(server.ListenAndServe())
}

func emailRegister(deps *HandlerDeps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var user models.User
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&user)
		if err != nil {
			http.Error(w, "Error while decoding json body", http.StatusInternalServerError)
			return
		}

		err = models.CreateUser(deps.DB, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		output, err := json.Marshal(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(output)

		return
	}
}
