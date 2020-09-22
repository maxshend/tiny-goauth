package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/maxshend/tiny_goauth/models"
	"github.com/maxshend/tiny_goauth/validations"
)

type TestDL struct {
	User models.User
}

func (t *TestDL) CreateUser(user *models.User) error {
	user.ID = t.User.ID
	user.CreatedAt = t.User.CreatedAt

	return nil
}

func (t *TestDL) Close() {}

func (t *TestDL) ExistsWithField(field, value string) (bool, error) {
	return false, nil
}

func TestEmailRegister(t *testing.T) {
	testUser := models.User{ID: 1, Email: "test@mail.com", Password: "password", CreatedAt: time.Now()}
	db := &TestDL{User: testUser}
	validator, translator, err := validations.Init(db)
	if err != nil {
		t.Error(err)
	}

	deps := &Deps{DB: db, Validator: validator, Translator: translator}
	performRequest := func(t *testing.T, method string, body io.Reader, headers map[string]string) (recorder *httptest.ResponseRecorder) {
		t.Helper()

		request, err := http.NewRequest(method, "/email/register", body)
		if err != nil {
			t.Error(err)
		}

		for name, value := range headers {
			request.Header.Add(name, value)
		}

		recorder = httptest.NewRecorder()
		handler := http.HandlerFunc(EmailRegister(deps))

		handler.ServeHTTP(recorder, request)

		return recorder
	}
	validateStatusCode := func(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
		t.Helper()

		if got := recorder.Code; got != expected {
			t.Errorf("Returned wrong status code. Expected %d, got %d", expected, got)
		}
	}

	t.Run("returns MethodNotAllowed for non-POST requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", nil, nil)
		validateStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Conten-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "POST", nil, nil)
		validateStatusCode(t, recorder, http.StatusBadRequest)
	})

	t.Run("returns InternalServerError when body isn't valid json", func(t *testing.T) {
		headers := map[string]string{contentTypeHeader: jsonContentType}
		recorder := performRequest(t, "POST", strings.NewReader("invalid"), headers)
		validateStatusCode(t, recorder, http.StatusInternalServerError)
	})

	t.Run("returns UnprocessableEntity with invalid user data", func(t *testing.T) {
		headers := map[string]string{contentTypeHeader: jsonContentType}
		body := bytes.NewBuffer([]byte(`{"email": "invalid.mail.com", "password": "foobar123"}`))
		recorder := performRequest(t, "POST", body, headers)
		validateStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns OK with valid user data", func(t *testing.T) {
		headers := map[string]string{contentTypeHeader: jsonContentType}
		body := bytes.NewBuffer([]byte(`{"email": "valid@mail.com", "password": "12345678"}`))
		recorder := performRequest(t, "POST", body, headers)
		validateStatusCode(t, recorder, http.StatusOK)
	})
}
