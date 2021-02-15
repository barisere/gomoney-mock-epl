package tests

import (
	"context"
	"gomoney-mock-epl/fixtures"
	"gomoney-mock-epl/web"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func clearFixtures() {
	testApp.app.FixturesDB.DeleteMany(context.Background(), bson.D{})
}

func Test_actions_on_fixtures(t *testing.T) {
	clearTeamsDB()
	clearFixtures()
	createTeam(manUtd)
	createTeam(liverpool)
	teams, err := testApp.app.TeamsDB.List(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, teams)
	assert.Len(t, teams, 2)

	createFixtureDto := fixtures.CreateFixtureRequest{
		HomeTeam:  teams[0].ID,
		AwayTeam:  teams[1].ID,
		MatchDate: time.Now().Add(48 * time.Hour),
	}

	t.Run("admins can create fixtures", func(t *testing.T) {
		result := createFixture(createFixtureDto)
		assert.Equal(t, http.StatusCreated, result.StatusCode)
		responseBody := web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		assert.Equal(t, "Fixture", responseBody.Type)
	})

	t.Run("admins can list fixtures", func(t *testing.T) {
		req, rec := jsonRequest(http.MethodGet, "/fixtures/", nil, adminToken)
		testApp.app.ServeHTTP(rec, req)
		result := rec.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		responseBody := web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		assert.Equal(t, "Fixtures", responseBody.Type)
		assert.Equal(t, 1, len(responseBody.Data.([]interface{})))
	})

	t.Run("admins can view single fixtures", func(t *testing.T) {
		req, rec := jsonRequest(http.MethodGet, "/fixtures/", nil, adminToken)
		testApp.app.ServeHTTP(rec, req)
		result := rec.Result()
		responseBody := web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		id := responseBody.Data.([]interface{})[0].(map[string]interface{})["id"].(string)

		req, rec = jsonRequest(http.MethodGet, "/fixtures/"+id, nil, adminToken)
		testApp.app.ServeHTTP(rec, req)
		result = rec.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		responseBody = web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		viewedID := responseBody.Data.(map[string]interface{})["id"].(string)
		assert.Equal(t, id, viewedID)
	})

	t.Run("admins can edit fixtures", func(t *testing.T) {
		req, rec := jsonRequest(http.MethodGet, "/fixtures/", nil, adminToken)
		testApp.app.ServeHTTP(rec, req)
		result := rec.Result()
		responseBody := web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		id := responseBody.Data.([]interface{})[0].(map[string]interface{})["id"].(string)

		update := createFixtureDto
		update.MatchDate = update.MatchDate.Add(48 * time.Hour)
		req, rec = jsonRequest(http.MethodPatch, "/fixtures/"+id, update, adminToken)
		testApp.app.ServeHTTP(rec, req)
		result = rec.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		responseBody = web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		updatedTimeString := responseBody.Data.(map[string]interface{})["match_date"].(string)
		updateTime, _ := time.Parse(time.RFC3339, updatedTimeString)
		assert.Equal(t, update.MatchDate.UTC().Day(), updateTime.UTC().Day())
	})

	t.Run("admins can delete fixtures", func(t *testing.T) {
		req, rec := jsonRequest(http.MethodGet, "/fixtures/", nil, adminToken)
		testApp.app.ServeHTTP(rec, req)
		result := rec.Result()
		responseBody := web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		id := responseBody.Data.([]interface{})[0].(map[string]interface{})["id"].(string)

		req, rec = jsonRequest(http.MethodDelete, "/fixtures/"+id, nil, adminToken)
		testApp.app.ServeHTTP(rec, req)
		result = rec.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})

	t.Run("users can list fixtures", func(t *testing.T) {
		clearFixtures()
		createFixture(createFixtureDto)
		req, rec := jsonRequest(http.MethodGet, "/fixtures/", nil, userToken)
		testApp.app.ServeHTTP(rec, req)
		result := rec.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		responseBody := web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		assert.Equal(t, "Fixtures", responseBody.Type)
		assert.Equal(t, 1, len(responseBody.Data.([]interface{})))
	})

	t.Run("users can view single fixtures", func(t *testing.T) {
		clearFixtures()
		createFixture(createFixtureDto)
		req, rec := jsonRequest(http.MethodGet, "/fixtures/", nil, userToken)
		testApp.app.ServeHTTP(rec, req)
		result := rec.Result()
		responseBody := web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		id := responseBody.Data.([]interface{})[0].(map[string]interface{})["id"].(string)

		req, rec = jsonRequest(http.MethodGet, "/fixtures/"+id, nil, userToken)
		testApp.app.ServeHTTP(rec, req)
		result = rec.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		responseBody = web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		viewedID := responseBody.Data.(map[string]interface{})["id"].(string)
		assert.Equal(t, id, viewedID)
	})

	t.Run("users can view completed or pending fixtures", func(t *testing.T) {
		clearFixtures()
		pendingFixture := createFixtureDto
		completedFixture := createFixtureDto
		completedFixture.MatchDate = time.Now().Add(-72 * time.Hour)
		createFixture(pendingFixture)
		createFixture(completedFixture)

		req, rec := jsonRequest(http.MethodGet, "/fixtures/?status=pending", nil, userToken)
		testApp.app.ServeHTTP(rec, req)
		result := rec.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		responseBody := web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		data := responseBody.Data.([]interface{})
		assert.Len(t, data, 1)
		homeTeamID := data[0].(map[string]interface{})["home_team"].(map[string]interface{})["id"].(string)
		assert.Equal(t, pendingFixture.HomeTeam, homeTeamID)

		req, rec = jsonRequest(http.MethodGet, "/fixtures/?status=completed", nil, userToken)
		testApp.app.ServeHTTP(rec, req)
		result = rec.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		responseBody = web.DataDto{}
		readJsonResponse(result.Body, &responseBody)
		data = responseBody.Data.([]interface{})
		assert.Len(t, data, 1)
		homeTeamID = data[0].(map[string]interface{})["home_team"].(map[string]interface{})["id"].(string)
		assert.Equal(t, completedFixture.HomeTeam, homeTeamID)
	})
}

func createFixture(dto fixtures.CreateFixtureRequest) *http.Response {
	req, rec := jsonRequest(http.MethodPost, "/fixtures/", dto, adminToken)
	testApp.app.ServeHTTP(rec, req)
	result := rec.Result()
	return result
}
