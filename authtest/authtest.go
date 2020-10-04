package authtest

import (
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

// import (
// 	"io"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/go-playground/validator"
// 	"github.com/maxshend/tiny_goauth/auth"
// 	"github.com/maxshend/tiny_goauth/handlers"
// 	"github.com/maxshend/tiny_goauth/logwrapper"
// 	"github.com/maxshend/tiny_goauth/models"
// 	"github.com/maxshend/tiny_goauth/validations"
// )

// GenerateFakeJWT generates fake JWT token with valid format
func GenerateFakeJWT(t *testing.T, sign []byte, method jwt.SigningMethod, claims jwt.Claims) string {
	t.Helper()

	jwtToken := jwt.NewWithClaims(method, claims)
	token, err := jwtToken.SignedString(sign)
	if err != nil {
		t.Fatal(err)
	}

	return token
}

// // PerformRequest performs http request using test recorder
// func PerformRequest(t *testing.T, method, path string, h func(deps *handlers.Deps) http.Handler, body io.Reader, headers map[string]string) (recorder *httptest.ResponseRecorder) {
// 	t.Helper()

// 	testUser := models.User{ID: 1, Email: "test@mail.com", Password: "password", CreatedAt: time.Now()}
// 	db := &TestDL{User: testUser}
// 	validator, translator, err := validations.Init(db)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	logger := logwrapper.New()
// 	logger.SetOutput(ioutil.Discard)

// 	deps := &handlers.Deps{DB: db, Validator: validator, Translator: translator, Logger: logger}

// 	request, err := http.NewRequest(method, path, body)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	for name, value := range headers {
// 		request.Header.Add(name, value)
// 	}

// 	recorder = httptest.NewRecorder()
// 	handler := h(deps)

// 	handler.ServeHTTP(recorder, request)

// 	return recorder
// }

// AssertStatusCode asserts HTTP status code
func AssertStatusCode(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if got := recorder.Code; got != expected {
		t.Errorf("Returned wrong status code. Expected %d, got %d", expected, got)
	}
}

// // TestDL represents data layer that can be used in tests
// type TestDL struct {
// 	User models.User
// }

// // CreateUser imitates creation of a user
// func (t *TestDL) CreateUser(user *models.User) error {
// 	user.ID = t.User.ID
// 	user.CreatedAt = t.User.CreatedAt

// 	return nil
// }

// // UserByEmail imitates search of the user by email
// func (t *TestDL) UserByEmail(email string) (*models.User, error) {
// 	var err error

// 	t.User.Password, err = auth.EncryptPassword(t.User.Password)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &t.User, nil
// }

// // Close imitates closing of datastores connections
// func (t *TestDL) Close() {}

// // UserExistsWithField imitates searching of the user by field/value
// func (t *TestDL) UserExistsWithField(fl validator.FieldLevel) (bool, error) {
// 	return false, nil
// }

// // StoreCache imitates storing data in the cache store
// func (t *TestDL) StoreCache(key string, payload interface{}, exp time.Duration) error {
// 	return nil
// }

// // DeleteCache imitates deleting data from the cache store
// func (t *TestDL) DeleteCache(key string) (int64, error) {
// 	return 1, nil
// }

// // GetCacheValue imitates retrieving from the cache store
// func (t *TestDL) GetCacheValue(key string) (string, error) {
// 	return "", nil
// }
