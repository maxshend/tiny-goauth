package handlers

import (
	"crypto/rsa"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/maxshend/tiny_goauth/auth"
	"github.com/maxshend/tiny_goauth/authtest"
	"github.com/maxshend/tiny_goauth/logwrapper"
	"github.com/maxshend/tiny_goauth/models"
	"github.com/maxshend/tiny_goauth/validations"
)

func TestLogout(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-DELETE requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/logout", Logout, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Content-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "DELETE", "/logout", Logout, nil, nil, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusBadRequest)
	})

	privateKey, err := authtest.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	claims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * 15).Unix()}
	expiredClaims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * -15).Unix()}
	token := authtest.GenerateFakeJWT(t, privateKey, jwt.SigningMethodRS256, claims)
	expired := authtest.GenerateFakeJWT(t, privateKey, jwt.SigningMethodRS256, expiredClaims)

	t.Run("returns Unauthorized with invalid 'Authorization' header", func(t *testing.T) {
		h := jsonHeaders
		h[auhtorizationHeader] = expired
		recorder := performRequest(t, "DELETE", "/logout", Logout, nil, h, privateKey)

		authtest.AssertStatusCode(t, recorder, http.StatusUnauthorized)
	})

	t.Run("returns OK with valid token", func(t *testing.T) {
		h := jsonHeaders
		h[auhtorizationHeader] = token
		recorder := performRequest(t, "DELETE", "/logout", Logout, nil, h, privateKey)

		authtest.AssertStatusCode(t, recorder, http.StatusOK)
	})
}

func TestRefresh(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-POST requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/refresh", Refresh, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Content-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/refresh", Refresh, nil, nil, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusBadRequest)
	})

	privateKey, err := authtest.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	claims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * 15).Unix()}
	expiredClaims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * -15).Unix()}
	token := authtest.GenerateFakeJWT(t, privateKey, jwt.SigningMethodRS256, claims)
	expired := authtest.GenerateFakeJWT(t, privateKey, jwt.SigningMethodRS256, expiredClaims)

	t.Run("returns Unauthorized with invalid 'Authorization' header", func(t *testing.T) {
		h := jsonHeaders
		h[auhtorizationHeader] = expired
		recorder := performRequest(t, "POST", "/refresh", Refresh, nil, h, privateKey)

		authtest.AssertStatusCode(t, recorder, http.StatusUnauthorized)
	})

	token = authtest.GenerateFakeJWT(t, privateKey, jwt.SigningMethodRS256, claims)

	t.Run("returns OK with valid Refresh token", func(t *testing.T) {
		h := jsonHeaders
		h[auhtorizationHeader] = token
		recorder := performRequest(t, "POST", "/refresh", Refresh, nil, h, privateKey)

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
	if id == 0 {
		return invalidUserID
	}

	return nil
}

func (t *testDL) GetRoles() ([]string, error) {
	return make([]string, 0), nil
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
func (t *testDL) Migrate() error {
	return nil
}

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

func (t *testDL) CreateRoles(names []string) error {
	if names[0] == "duplicate" {
		return errors.New("duplicate")
	}

	return nil
}

func (t *testDL) DeleteRole(name string) error {
	if name == "not_found" {
		return errors.New("not found")
	}

	return nil
}
