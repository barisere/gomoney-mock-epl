package tests

import (
	"bytes"
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

const (
	testAdminEmail = "jon.doe@gomoney.local"
	testPassword   = "password"
	testUserEmail  = "jane.doe@gomoney.local"
)

type testFixtures struct {
	app *web.Application
}

func (t testFixtures) setUpAdminAccount() error {
	intent := users.SignUpIntent{
		Email:     testAdminEmail,
		FirstName: "Jon",
		LastName:  "Doe",
		Password:  testPassword,
	}
	_, err := users.SignUpAdmin(context.Background(), intent, t.app.AdminDB)
	return err
}

func (t testFixtures) setUpUserAccount() error {
	intent := users.SignUpIntent{
		Email:     testUserEmail,
		FirstName: "Jane",
		LastName:  "Doe",
		Password:  testPassword,
	}
	_, err := users.SignUpUser(context.Background(), intent, t.app.UsersDB)
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

func jsonRequest(method, path string, data interface{}, token string) (*http.Request, *httptest.ResponseRecorder) {
	body, _ := json.Marshal(data)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
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
