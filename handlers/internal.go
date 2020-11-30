package handlers

import (
	"net/http"
	"strconv"
)

const invalidUserID = handlerErr("Invalid User ID")

// DeleteUser removes a user record
func DeleteUser(deps *Deps) http.Handler {
	return logHandler(deps, jsonHandler(deleteHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if err != nil {
			respondError(w, http.StatusUnprocessableEntity, invalidUserID)
			return
		}

		err = deps.DB.DeleteUser(userID)
		if err != nil {
			respondError(w, http.StatusUnprocessableEntity, err)
			return
		}
	}))))
}
