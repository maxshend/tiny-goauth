package handlers

import (
	"net/http"
	"strconv"
)

const invalidUserID = handlerErr("Invalid User ID")
const blankRole = handlerErr("Blank Role Name")

// DeleteUser removes a user record
func DeleteUser(deps *Deps) http.Handler {
	return logHandler(deps, jsonHandler(deleteHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if err != nil {
			respondError(w, http.StatusUnprocessableEntity, invalidUserID)
			return
		}

		if err = deps.DB.DeleteUser(userID); err != nil {
			respondError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
	}))))
}

// CreateRole creates a role record
func CreateRole(deps *Deps) http.Handler {
	return logHandler(deps, jsonHandler(postHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")

		if len(name) == 0 {
			respondError(w, http.StatusUnprocessableEntity, blankRole)
			return
		}

		if err := deps.DB.CreateRole(name); err != nil {
			respondError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
	}))))
}
