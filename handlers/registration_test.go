package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maxshend/tiny_goauth/validations"
)

func TestEmailRegister(t *testing.T) {
	validator, translator, err := validations.Init()
	if err != nil {
		log.Fatal(err)
	}

	deps := &Deps{DB: nil, Validator: validator, Translator: translator}
	performRequest := func(t *testing.T, method string, body io.Reader, headers map[string]string) (recorder *httptest.ResponseRecorder) {
		t.Helper()

		request, err := http.NewRequest(method, "/email/register", body)
		if err != nil {
			t.Fatal(err)
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
}
