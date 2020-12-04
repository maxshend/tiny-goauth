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

	dbInst, err := db.Init()
	if err != nil {
		logger.FatalError(err)
	}
	defer dbInst.Close()

	migrateDB, present := os.LookupEnv("MIGRATE_DB")
	if present && migrateDB == "true" {
		if err = dbInst.Migrate(); err != nil {
			logger.FatalError(err)
		}
	}

	validator, translator, err := validations.Init(dbInst)
	if err != nil {
		logger.FatalError(err)
	}

	keys, err := auth.Keys()
	if err != nil {
		logger.FatalError(err)
	}

	deps := &handlers.Deps{
		DB:         dbInst,
		Validator:  validator,
		Translator: translator,
		Logger:     logger,
		Keys:       keys,
	}
	server := http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	http.Handle("/email/register", handlers.EmailRegister(deps))
	http.Handle("/email/login", handlers.EmailLogin(deps))
	http.Handle("/logout", handlers.Logout(deps))
	http.Handle("/refresh", handlers.Refresh(deps))
	http.Handle("/internal/users/delete", handlers.DeleteUser(deps))
	http.Handle("/internal/roles", handlers.CreateRole(deps))

	logger.FatalError(server.ListenAndServe())
}
