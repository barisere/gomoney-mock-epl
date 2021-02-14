package tests

import (
	"context"
	"gomoney-mock-epl/web"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func createTeam(dto web.CreateTeamRequest) *httptest.ResponseRecorder {
	req, rec := jsonRequest(http.MethodPost, "/teams/", &dto, adminToken)
	fixture.app.ServeHTTP(rec, req)
	return rec
}

func clearTeamsDB() {
	fixture.app.TeamsDB.DeleteMany(context.Background(), bson.D{})
}

func Test_admins_can_create_teams(t *testing.T) {
	clearTeamsDB()
	result := createTeam(manUtd).Result()
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	responseBody := web.DataDto{}
	assert.NoError(t, readJsonResponse(result.Body, &responseBody))
	assert.Equal(t, "Team", responseBody.Type)
}

var manUtd = web.CreateTeamRequest{
	City:        "Manchester",
	HomeStadium: "Old Trafford",
	LogoURL:     "https://resources.premierleague.com/premierleague/badges/50/t1.png",
	Name:        "Manchester United",
	NameAbbr:    "MUN",
	ShortName:   "Man Utd",
}

var manCity = web.CreateTeamRequest{
	City:        "Manchester",
	HomeStadium: "Etihad Stadiun",
	LogoURL:     "https://resources.premierleague.com/premierleague/badges/70/t43@x2.png",
	Name:        "Manchester City",
	NameAbbr:    "MCI",
	ShortName:   "Man City",
}

var liverpool = web.CreateTeamRequest{
	City:        "Liverpool",
	HomeStadium: "Anfield",
	LogoURL:     "https://resources.premierleague.com/premierleague/badges/70/t14@x2.png",
	Name:        "Liverpool",
	NameAbbr:    "LIV",
	ShortName:   "Liverpool",
}

func Test_admins_can_view_teams(t *testing.T) {
	clearTeamsDB()
	createTeam(manCity)
	createTeam(liverpool)
	req, rec := jsonRequest(http.MethodGet, "/teams/", nil, adminToken)
	fixture.app.ServeHTTP(rec, req)
	result := rec.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	responseBody := web.DataDto{}
	assert.NoError(t, readJsonResponse(result.Body, &responseBody))
	assert.Equal(t, "Teams", responseBody.Type)
}

func Test_admins_can_view_single_team(t *testing.T) {
	clearTeamsDB()
	createTeam(liverpool)
	req, rec := jsonRequest(http.MethodGet, "/teams/", nil, adminToken)
	fixture.app.ServeHTTP(rec, req)
	result := rec.Result()
	responseBody := web.DataDto{}
	assert.NoError(t, readJsonResponse(result.Body, &responseBody))
	firstID := responseBody.Data.([]interface{})[0].(map[string]interface{})["id"].(string)

	req, rec = jsonRequest(http.MethodGet, "/teams/"+firstID, nil, adminToken)
	fixture.app.ServeHTTP(rec, req)
	result = rec.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	responseBody = web.DataDto{}
	assert.NoError(t, readJsonResponse(result.Body, &responseBody))
	teamName := responseBody.Data.(map[string]interface{})["name"].(string)
	assert.Equal(t, liverpool.Name, teamName)
}

func Test_admins_can_remove_teams(t *testing.T) {
	clearTeamsDB()
	createTeam(manCity)
	createTeam(liverpool)
	req, rec := jsonRequest(http.MethodGet, "/teams/", nil, adminToken)
	fixture.app.ServeHTTP(rec, req)
	result := rec.Result()
	responseBody := web.DataDto{}
	assert.NoError(t, readJsonResponse(result.Body, &responseBody))
	firstID := responseBody.Data.([]interface{})[0].(map[string]interface{})["id"].(string)

	req, rec = jsonRequest(http.MethodDelete, "/teams/"+firstID, nil, adminToken)
	fixture.app.ServeHTTP(rec, req)
	result = rec.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func Test_admins_can_edit_teams(t *testing.T) {
	clearTeamsDB()
	createTeam(liverpool)
	req, rec := jsonRequest(http.MethodGet, "/teams/", nil, adminToken)
	fixture.app.ServeHTTP(rec, req)
	result := rec.Result()
	responseBody := web.DataDto{}
	assert.NoError(t, readJsonResponse(result.Body, &responseBody))
	firstID := responseBody.Data.([]interface{})[0].(map[string]interface{})["id"].(string)

	liverpoolCopy := liverpool
	liverpoolCopy.HomeStadium = "Stamford Bridge" // travesty!
	req, rec = jsonRequest(http.MethodPatch, "/teams/"+firstID, liverpoolCopy, adminToken)
	fixture.app.ServeHTTP(rec, req)
	result = rec.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	responseBody = web.DataDto{}
	assert.NoError(t, readJsonResponse(result.Body, &responseBody))
	newStadium := responseBody.Data.(map[string]interface{})["home_stadium"].(string)
	assert.Equal(t, liverpoolCopy.HomeStadium, newStadium)
}

func Test_users_can_view_teams(t *testing.T) {
	clearTeamsDB()
	createTeam(manCity)
	createTeam(liverpool)
	req, rec := jsonRequest(http.MethodGet, "/teams/", nil, userToken)
	fixture.app.ServeHTTP(rec, req)
	result := rec.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	responseBody := web.DataDto{}
	assert.NoError(t, readJsonResponse(result.Body, &responseBody))
	assert.Equal(t, "Teams", responseBody.Type)
}

func Test_users_can_view_single_team(t *testing.T) {
	clearTeamsDB()
	createTeam(liverpool)
	req, rec := jsonRequest(http.MethodGet, "/teams/", nil, userToken)
	fixture.app.ServeHTTP(rec, req)
	result := rec.Result()
	responseBody := web.DataDto{}
	assert.NoError(t, readJsonResponse(result.Body, &responseBody))
	firstID := responseBody.Data.([]interface{})[0].(map[string]interface{})["id"].(string)

	req, rec = jsonRequest(http.MethodGet, "/teams/"+firstID, nil, adminToken)
	fixture.app.ServeHTTP(rec, req)
	result = rec.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	responseBody = web.DataDto{}
	assert.NoError(t, readJsonResponse(result.Body, &responseBody))
	teamName := responseBody.Data.(map[string]interface{})["name"].(string)
	assert.Equal(t, liverpool.Name, teamName)
}
