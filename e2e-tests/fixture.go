package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"gomoney-mock-epl/config"
	"gomoney-mock-epl/database"
	"gomoney-mock-epl/users"
	"gomoney-mock-epl/web"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type testFixtures struct {
	app *web.Application
}

var (
	testAdmin = &users.Administrator{
		Email:     "jon.doe@gomoney.local",
		FirstName: "Jon",
		LastName:  "Doe",
	}
	testUser = &users.User{
		Email:     "jane.doe@gomoney.local",
		FirstName: "Jane",
		LastName:  "Doe",
	}
	testPassword = "password"
)

func (t testFixtures) setUpAdminAccount() error {
	intent := users.SignUpIntent{
		Email:     testAdmin.Email,
		FirstName: testAdmin.FirstName,
		LastName:  testAdmin.LastName,
		Password:  testPassword,
	}
	var err error
	testAdmin, err = users.SignUpAdmin(context.Background(), intent, *t.app.AdminDB)
	return err
}

func (t testFixtures) setUpUserAccount() error {
	intent := users.SignUpIntent{
		Email:     testUser.Email,
		FirstName: testUser.FirstName,
		LastName:  testUser.LastName,
		Password:  testPassword,
	}
	var err error
	testUser, err = users.SignUpUser(context.Background(), intent, *t.app.UsersDB)
	return err
}

func (t testFixtures) destroy(ctx context.Context) {
	t.app.DBClient.Database(database.MockEPLDatabase).Drop(ctx)
	t.app.DBClient.Disconnect(ctx)
}

func setUpFixtures() *testFixtures {
	config := config.Config{
		Environment:  config.Testing,
		HttpBindPort: 8080,
		MongoURL:     "mongodb://localhost:27017/hf?ssl=false",
	}
	dbClient, err := database.ConnectToDB(config.MongoURL)
	if err != nil {
		panic(err)
	}
	app, err := web.NewApplication(dbClient, config)
	if err != nil {
		panic(err)
	}

	fixture := testFixtures{
		app: app,
	}
	fixture.setUpAdminAccount()
	fixture.setUpUserAccount()

	return &fixture
}

func jsonRequest(method, path string, body io.Reader, token string) (*http.Request, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	return req, rec
}

func readJsonResponse(body io.ReadCloser, dst interface{}) error {
	data, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}
