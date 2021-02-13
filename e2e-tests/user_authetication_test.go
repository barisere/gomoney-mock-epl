package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"gomoney-mock-epl/users"
	"gomoney-mock-epl/web"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertThatAdminAccountWasCreated(t *testing.T, db *users.AdminsDB, accountID string) {
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

	reqBody, _ := json.Marshal(&intent)
	fixture := setUpFixtures()
	fixture.setUpAdminAccount()
	defer fixture.destroy(context.Background())

	t.Run("fails given invalid admin authentication", func(t *testing.T) {
		req, rec := jsonRequest(http.MethodPost, "/signup/admins/", bytes.NewReader(reqBody), "")
		fixture.app.ServeHTTP(rec, req)
		result := rec.Result()
		assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})

	t.Run("succeeds given valid admin authentication", func(t *testing.T) {
		loginResult := loginAsAdmin(users.LoginDto{
			Email:    testAdmin.Email,
			Password: testPassword,
		}, *fixture).Result()
		loginResponse := web.DataDto{}
		assert.NoError(t, readJsonResponse(loginResult.Body, &loginResponse))
		token := loginResponse.Data.(map[string]interface{})["token"].(string)
		req, rec := jsonRequest(http.MethodPost, "/signup/admins/", bytes.NewReader(reqBody), token)
		fixture.app.ServeHTTP(rec, req)
		response := web.DataDto{}
		result := rec.Result()
		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.NoError(t, readJsonResponse(result.Body, &response))
		assertThatAdminAccountWasCreated(t, fixture.app.AdminDB, response.Data.(map[string]interface{})["id"].(string))
	})
}

func loginAsAdmin(dto users.LoginDto, fixture testFixtures) *httptest.ResponseRecorder {
	reqBody, _ := json.Marshal(&dto)
	req, rec := jsonRequest(http.MethodPost, "/login/admins/", bytes.NewReader(reqBody), "")
	fixture.app.ServeHTTP(rec, req)
	return rec
}

func Test_admin_login(t *testing.T) {
	fixture := setUpFixtures()
	assert.NoError(t, fixture.setUpAdminAccount())
	loginDto := users.LoginDto{
		Email:    testAdmin.Email,
		Password: testPassword,
	}
	defer fixture.destroy(context.Background())

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

	reqBody, _ := json.Marshal(&intent)
	fixture := setUpFixtures()
	fixture.setUpUserAccount()
	defer fixture.destroy(context.Background())

	req, rec := jsonRequest(http.MethodPost, "/signup/users/", bytes.NewReader(reqBody), "")
	fixture.app.ServeHTTP(rec, req)
	response := web.DataDto{}
	result := rec.Result()
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.NoError(t, readJsonResponse(result.Body, &response))
	assertThatAdminAccountWasCreated(t, fixture.app.AdminDB, response.Data.(map[string]interface{})["id"].(string))
}

func loginAsUser(dto users.LoginDto, fixture testFixtures) *httptest.ResponseRecorder {
	reqBody, _ := json.Marshal(&dto)
	req, rec := jsonRequest(http.MethodPost, "/login/users/", bytes.NewReader(reqBody), "")
	fixture.app.ServeHTTP(rec, req)
	return rec
}

func assertThatUserAccountWasCreated(t *testing.T, db *users.UsersDB, accountID string) {
	user, err := db.ByID(context.Background(), accountID)
	assert.Nil(t, err, "Unexpected error retrieving account")
	assert.NotEmpty(t, user, "User account was not saved to DB")
}

func Test_user_login(t *testing.T) {
	fixture := setUpFixtures()
	assert.NoError(t, fixture.setUpUserAccount())
	loginDto := users.LoginDto{
		Email:    testAdmin.Email,
		Password: testPassword,
	}
	defer fixture.destroy(context.Background())

	t.Run("fails for invalid credentials", func(t *testing.T) {
		result := loginAsUser(users.LoginDto{}, *fixture).Result()
		assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})

	t.Run("succeeds given correct credentials", func(t *testing.T) {
		result := loginAsUser(loginDto, *fixture).Result()
		response := web.DataDto{}
		assert.NoError(t, readJsonResponse(result.Body, &response))
		token := response.Data.(map[string]interface{})["token"].(string)
		assert.NotEmpty(t, token)
	})
}
