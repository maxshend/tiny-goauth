package main

import (
	"net/http"
	"os"
	"time"

	"github.com/maxshend/tiny_goauth/auth"
	"github.com/maxshend/tiny_goauth/db"
	"github.com/maxshend/tiny_goauth/handlers"
	"github.com/maxshend/tiny_goauth/logwrapper"
	"github.com/maxshend/tiny_goauth/validations"
)

func main() {
	logger := logwrapper.New()

	db, err := db.Init()
	if err != nil {
		logger.FatalError(err)
	}
	defer db.Close()

	validator, translator, err := validations.Init(db)
	if err != nil {
		logger.FatalError(err)
	}

	keys, err := auth.Keys()
	if err != nil {
		logger.FatalError(err)
	}

	deps := &handlers.Deps{
		DB:         db,
		Validator:  validator,
		Translator: translator,
		Logger:     logger,
		Keys:       keys,
	}
	server := http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		Handler:      nil,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	http.Handle("/email/register", handlers.EmailRegister(deps))
	http.Handle("/email/login", handlers.EmailLogin(deps))
	http.Handle("/logout", handlers.Logout(deps))
	http.Handle("/refresh", handlers.Refresh(deps))

	logger.FatalError(server.ListenAndServe())
}
