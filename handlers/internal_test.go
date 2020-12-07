package handlers

import (
	"bytes"
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

func TestCreateRoles(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-post requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/internal/roles", CreateRoles, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Content-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/internal/roles", CreateRoles, nil, nil, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusBadRequest)
	})

	t.Run("returns UnprocessableEntity when Role already exists", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"roles": ["duplicate"]}`))
		recorder := performRequest(t, "POST", "/internal/roles", CreateRoles, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns UnprocessableEntity with blank roles array", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"roles": []}`))
		recorder := performRequest(t, "POST", "/internal/roles", CreateRoles, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns UnprocessableEntity with invalid roles type", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"roles": "invalid"}`))
		recorder := performRequest(t, "POST", "/internal/roles", CreateRoles, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns UnprocessableEntity with invalid Role Name", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"roles": [""]}`))
		recorder := performRequest(t, "POST", "/internal/roles", CreateRoles, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns OK with valid Role", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"roles": ["test"]}`))
		recorder := performRequest(t, "POST", "/internal/roles", CreateRoles, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusOK)
	})
}

func TestDeleteRole(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-delete requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/internal/roles/delete", DeleteRoles, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Content-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "DELETE", "/internal/roles/delete", DeleteRoles, nil, nil, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusBadRequest)
	})

	t.Run("returns UnprocessableEntity when Role with name doesn't exist", func(t *testing.T) {
		recorder := performRequest(t, "DELETE", "/internal/roles/delete?roles=not_found", DeleteRoles, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns UnprocessableEntity with invalid Role Name", func(t *testing.T) {
		recorder := performRequest(t, "DELETE", "/internal/roles/delete?roles=", DeleteRoles, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns OK with valid Role Name", func(t *testing.T) {
		recorder := performRequest(t, "DELETE", "/internal/roles/delete?roles=test", DeleteRoles, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusOK)
	})
}
