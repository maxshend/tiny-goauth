package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const invalidUserID = handlerErr("Invalid User ID")
const blankRoles = handlerErr("Blank Roles")
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

// CreateRoles creates a role record
func CreateRoles(deps *Deps) http.Handler {
	return logHandler(deps, jsonHandler(postHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

		body := make(map[string][]string)
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&body)
		if err != nil {
			respondError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		var ok bool
		var roles []string

		if roles, ok = body["roles"]; !ok || len(roles) == 0 {
			respondError(w, http.StatusUnprocessableEntity, blankRoles)
			return
		}

		for _, role := range roles {
			if len(role) == 0 {
				respondError(w, http.StatusUnprocessableEntity, blankRole)
				return
			}
		}

		if err := deps.DB.CreateRoles(roles); err != nil {
			respondError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
	}))))
}

// DeleteRoles removes a role record
func DeleteRoles(deps *Deps) http.Handler {
	return logHandler(deps, jsonHandler(deleteHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		roles := r.Form["roles"]
		if len(roles) == 0 {
			respondError(w, http.StatusUnprocessableEntity, blankRoles)
			return
		}

		for _, role := range roles {
			if len(role) == 0 {
				respondError(w, http.StatusUnprocessableEntity, blankRole)
				return
			}
		}

		if err := deps.DB.DeleteRoles(roles); err != nil {
			respondError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
	}))))
}
