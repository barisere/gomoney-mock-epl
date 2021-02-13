package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"gomoney-mock-epl/config"
	"gomoney-mock-epl/users"
	"gomoney-mock-epl/web"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func connectToDB(mongoURL string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		panic(err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	return client
}

type testFixtures struct {
	adminsDB *users.AdminsDB
	app      *web.Application
}

var (
	testAdmin = &users.Administrator{
		Email:     "jon.doe@gomoney.local",
		FirstName: "Jon",
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
	testAdmin, err = users.SignUpAdmin(context.Background(), intent, *t.adminsDB)
	return err
}

func (t testFixtures) destroy(ctx context.Context) {
	t.app.DBClient.Disconnect(ctx)
}

func setUpFixtures() *testFixtures {
	config := config.Config{
		Environment:  config.Testing,
		HttpBindPort: 8080,
		MongoURL:     "mongodb://localhost:27017/hf?ssl=false",
	}
	dbClient := connectToDB(config.MongoURL)
	app, err := web.NewApplication(dbClient, config)
	if err != nil {
		panic(err)
	}

	web.AdminSignupRoute(*app.AdminDB)(app.Echo)
	fixture := testFixtures{
		app:      app,
		adminsDB: app.AdminDB,
	}

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
