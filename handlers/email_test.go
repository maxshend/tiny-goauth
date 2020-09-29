package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/maxshend/tiny_goauth/auth"
	"github.com/maxshend/tiny_goauth/models"
	"github.com/maxshend/tiny_goauth/validations"
)

var jsonHeaders = map[string]string{contentTypeHeader: jsonContentType}

func TestEmailRegister(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-POST requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/email/register", EmailRegister, nil, jsonHeaders)

		validateStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Conten-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/register", EmailRegister, nil, nil)

		validateStatusCode(t, recorder, http.StatusBadRequest)
	})

	t.Run("returns InternalServerError when body isn't valid json", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/register", EmailRegister, strings.NewReader("invalid"), jsonHeaders)

		validateStatusCode(t, recorder, http.StatusInternalServerError)
	})

	t.Run("returns UnprocessableEntity with invalid user data", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "invalid.mail.com", "password": "foobar123"}`))
		recorder := performRequest(t, "POST", "/email/register", EmailRegister, body, jsonHeaders)

		validateStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns OK with valid user data", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "valid@mail.com", "password": "12345678"}`))
		recorder := performRequest(t, "POST", "/email/register", EmailRegister, body, jsonHeaders)

		validateStatusCode(t, recorder, http.StatusOK)
	})
}

func TestEmailLogin(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-POST requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/email/login", EmailLogin, nil, jsonHeaders)

		validateStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Conten-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, nil, nil)

		validateStatusCode(t, recorder, http.StatusBadRequest)
	})

	t.Run("returns InternalServerError when body isn't valid json", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, strings.NewReader("invalid"), jsonHeaders)

		validateStatusCode(t, recorder, http.StatusInternalServerError)
	})

	t.Run("returns Unauthorized with invalid user creds", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "invalid.mail.com", "password": "foobar123"}`))
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, body, jsonHeaders)

		validateStatusCode(t, recorder, http.StatusUnauthorized)
	})

	t.Run("returns OK with valid user creds", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "test@mail.com", "password": "password"}`))
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, body, jsonHeaders)

		validateStatusCode(t, recorder, http.StatusOK)
	})
}

func performRequest(t *testing.T, method, path string, h func(deps *Deps) http.Handler, body io.Reader, headers map[string]string) (recorder *httptest.ResponseRecorder) {
	t.Helper()

	testUser := models.User{ID: 1, Email: "test@mail.com", Password: "password", CreatedAt: time.Now()}
	db := &TestDL{User: testUser}
	validator, translator, err := validations.Init(db)
	if err != nil {
		t.Error(err)
	}

	deps := &Deps{DB: db, Validator: validator, Translator: translator}

	request, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Error(err)
	}

	for name, value := range headers {
		request.Header.Add(name, value)
	}

	recorder = httptest.NewRecorder()
	handler := h(deps)

	handler.ServeHTTP(recorder, request)

	return recorder
}

func validateStatusCode(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if got := recorder.Code; got != expected {
		t.Errorf("Returned wrong status code. Expected %d, got %d", expected, got)
	}
}

type TestDL struct {
	User models.User
}

func (t *TestDL) CreateUser(user *models.User) error {
	user.ID = t.User.ID
	user.CreatedAt = t.User.CreatedAt

	return nil
}

func (t *TestDL) UserByEmail(email string) (*models.User, error) {
	var err error

	t.User.Password, err = auth.EncryptPassword(t.User.Password)
	if err != nil {
		return nil, err
	}

	return &t.User, nil
}

func (t *TestDL) Close() {}

func (t *TestDL) UserExistsWithField(fl validator.FieldLevel) (bool, error) {
	return false, nil
}

func (t *TestDL) StoreCache(key string, payload interface{}, exp time.Duration) error {
	return nil
}

func (t *TestDL) DeleteCache(key string) (int64, error) {
	return 1, nil
}

func (t *TestDL) GetCacheValue(key string) (string, error) {
	return "", nil
}
