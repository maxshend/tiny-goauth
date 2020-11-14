package handlers

import (
	"bytes"
	"crypto/rsa"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/maxshend/tiny_goauth/auth"
	"github.com/maxshend/tiny_goauth/authtest"
	"github.com/maxshend/tiny_goauth/logwrapper"
	"github.com/maxshend/tiny_goauth/models"
	"github.com/maxshend/tiny_goauth/validations"
)

var jsonHeaders = map[string]string{contentTypeHeader: jsonContentType}

func TestEmailRegister(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-POST requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/email/register", EmailRegister, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Conten-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/register", EmailRegister, nil, nil, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusBadRequest)
	})

	t.Run("returns InternalServerError when body isn't valid json", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/register", EmailRegister, strings.NewReader("invalid"), jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusInternalServerError)
	})

	t.Run("returns UnprocessableEntity with invalid user data", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "invalid.mail.com", "password": "foobar123"}`))
		recorder := performRequest(t, "POST", "/email/register", EmailRegister, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns OK with valid user data", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "valid@mail.com", "password": "12345678"}`))

		externalApp := testServer()
		defer externalApp.Close()

		os.Setenv("API_HOST", externalApp.URL)

		recorder := performRequest(t, "POST", "/email/register", EmailRegister, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusOK)
	})
}

func TestEmailLogin(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-POST requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/email/login", EmailLogin, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Conten-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, nil, nil, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusBadRequest)
	})

	t.Run("returns InternalServerError when body isn't valid json", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, strings.NewReader("invalid"), jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusInternalServerError)
	})

	t.Run("returns Unauthorized with invalid user creds", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "invalid.mail.com", "password": "foobar123"}`))
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnauthorized)
	})

	t.Run("returns OK with valid user creds", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "test@mail.com", "password": "password"}`))
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusOK)
	})
}

func performRequest(t *testing.T, method, path string, h func(deps *Deps) http.Handler, body io.Reader, headers map[string]string, key *rsa.PrivateKey) (recorder *httptest.ResponseRecorder) {
	t.Helper()

	testUser := models.User{ID: 1, Email: "test@mail.com", Password: "password", CreatedAt: time.Now()}
	db := &testDL{User: testUser}
	validator, translator, err := validations.Init(db)
	if err != nil {
		t.Error(err)
	}

	logger := logwrapper.New()
	logger.SetOutput(ioutil.Discard)

	if key == nil {
		key, err = authtest.GeneratePrivateKey()
		if err != nil {
			t.Fatal(err)
		}
	}
	keys := &auth.RSAKeys{AccessSign: key, AccessVerify: &key.PublicKey, RefreshSign: key, RefreshVerify: &key.PublicKey}

	deps := &Deps{DB: db, Validator: validator, Translator: translator, Logger: logger, Keys: keys}

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

type testDL struct {
	User models.User
}

func (t *testDL) CreateUser(user *models.User) error {
	user.ID = t.User.ID
	user.CreatedAt = t.User.CreatedAt

	return nil
}

func (t *testDL) DeleteUser(id int64) error {
	return nil
}

func (t *testDL) UserByEmail(email string) (*models.User, error) {
	var err error

	t.User.Password, err = auth.EncryptPassword(t.User.Password)
	if err != nil {
		return nil, err
	}

	return &t.User, nil
}

func (t *testDL) Close() {}

func (t *testDL) UserExistsWithField(fl validator.FieldLevel) (bool, error) {
	return false, nil
}

func (t *testDL) StoreCache(key string, payload interface{}, exp time.Duration) error {
	return nil
}

func (t *testDL) DeleteCache(key string) (int64, error) {
	return 1, nil
}

func (t *testDL) GetCacheValue(key string) (string, error) {
	return "", nil
}
