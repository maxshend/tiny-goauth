package handlers

import (
	"net/http"
	"testing"

	"github.com/maxshend/tiny_goauth/authtest"
)

func TestDeleteUser(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-delete requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/internal/users/delete", DeleteUser, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Content-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "DELETE", "/internal/users/delete", DeleteUser, nil, nil, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusBadRequest)
	})

	t.Run("returns UnprocessableEntity when User with ID doesn't exist", func(t *testing.T) {
		recorder := performRequest(t, "DELETE", "/internal/users/delete?id=0", DeleteUser, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns UnprocessableEntity with invalid User ID", func(t *testing.T) {
		recorder := performRequest(t, "DELETE", "/internal/users/delete?id=", DeleteUser, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns OK with valid User ID", func(t *testing.T) {
		recorder := performRequest(t, "DELETE", "/internal/users/delete?id=1", DeleteUser, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusOK)
	})
}
