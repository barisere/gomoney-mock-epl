package tests

import (
	"context"
	"gomoney-mock-epl/users"
	"gomoney-mock-epl/web"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fixture *testFixtures
var adminToken string
var userToken string

func TestMain(m *testing.M) {
	fixture = setUpFixtures()
	loginResult := loginAsAdmin(users.LoginDto{
		Email:    testAdminEmail,
		Password: testPassword,
	}, *fixture).Result()
	loginResponse := web.DataDto{}
	readJsonResponse(loginResult.Body, &loginResponse)
	adminToken = loginResponse.Data.(map[string]interface{})["token"].(string)

	loginResult = loginAsUser(users.LoginDto{
		Email:    testUserEmail,
		Password: testPassword,
	}, *fixture).Result()
	loginResponse = web.DataDto{}
	readJsonResponse(loginResult.Body, &loginResponse)
	userToken = loginResponse.Data.(map[string]interface{})["token"].(string)

	exitCode := m.Run()
	fixture.destroy(context.Background())
	os.Exit(exitCode)
}

func assertThatAdminAccountWasCreated(t *testing.T, db users.AdminsDB, accountID string) {
	admin, err := db.ByID(context.Background(), accountID)
	assert.Nil(t, err, "Unexpected error retrieving account")
	assert.NotEmpty(t, admin, "Admin account was not saved to DB")
}

func Test_creating_admin_account(t *testing.T) {
	var intent = users.SignUpIntent{
		Email:     "admin@example.com",
		FirstName: "Victor",
		LastName:  "Alabi",
		Password:  "really strong, 123456",
	}

	t.Run("fails given invalid admin authentication", func(t *testing.T) {
		req, rec := jsonRequest(http.MethodPost, "/signup/admins/", &intent, "")
		fixture.app.ServeHTTP(rec, req)
		result := rec.Result()
		assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})

	t.Run("succeeds given valid admin authentication", func(t *testing.T) {
		req, rec := jsonRequest(http.MethodPost, "/signup/admins/", &intent, adminToken)
		fixture.app.ServeHTTP(rec, req)
		response := web.DataDto{}
		result := rec.Result()
		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.NoError(t, readJsonResponse(result.Body, &response))
		assertThatAdminAccountWasCreated(t, fixture.app.AdminDB, response.Data.(map[string]interface{})["id"].(string))
	})
}

func loginAsAdmin(dto users.LoginDto, fixture testFixtures) *httptest.ResponseRecorder {
	req, rec := jsonRequest(http.MethodPost, "/login/admins/", &dto, "")
	fixture.app.ServeHTTP(rec, req)
	return rec
}

func Test_admin_login(t *testing.T) {
	loginDto := users.LoginDto{
		Email:    testAdminEmail,
		Password: testPassword,
	}

	t.Run("fails for invalid credentials", func(t *testing.T) {
		result := loginAsAdmin(users.LoginDto{}, *fixture).Result()
		assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})

	t.Run("succeeds given correct credentials", func(t *testing.T) {
		result := loginAsAdmin(loginDto, *fixture).Result()
		response := web.DataDto{}
		assert.NoError(t, readJsonResponse(result.Body, &response))
		token := response.Data.(map[string]interface{})["token"].(string)
		assert.NotEmpty(t, token)
	})
}

func Test_creating_user_account(t *testing.T) {
	var intent = users.SignUpIntent{
		Email:     "victor@gomoney.local",
		FirstName: "Victor",
		LastName:  "Alabi",
		Password:  "really strong, 123456",
	}

	req, rec := jsonRequest(http.MethodPost, "/signup/users/", intent, "")
	fixture.app.ServeHTTP(rec, req)
	response := web.DataDto{}
	result := rec.Result()
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.NoError(t, readJsonResponse(result.Body, &response))
	assertThatUserAccountWasCreated(t, fixture.app.UsersDB, response.Data.(map[string]interface{})["id"].(string))
}

func loginAsUser(dto users.LoginDto, fixture testFixtures) *httptest.ResponseRecorder {
	req, rec := jsonRequest(http.MethodPost, "/login/users/", dto, "")
	fixture.app.ServeHTTP(rec, req)
	return rec
}

func assertThatUserAccountWasCreated(t *testing.T, db users.UsersDB, accountID string) {
	user, err := db.ByID(context.Background(), accountID)
	assert.Nil(t, err, "Unexpected error retrieving account")
	assert.NotEmpty(t, user, "User account was not saved to DB")
}

func Test_user_login(t *testing.T) {
	loginDto := users.LoginDto{
		Email:    testUserEmail,
		Password: testPassword,
	}

	t.Run("fails for invalid credentials", func(t *testing.T) {
		result := loginAsUser(users.LoginDto{}, *fixture).Result()
		assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})

	t.Run("succeeds given correct credentials", func(t *testing.T) {
		result := loginAsUser(loginDto, *fixture).Result()
		response := web.DataDto{}
		assert.NoError(t, readJsonResponse(result.Body, &response))
		assert.Equal(t, http.StatusOK, result.StatusCode)
		token := response.Data.(map[string]interface{})["token"].(string)
		assert.NotEmpty(t, token)
	})
}
