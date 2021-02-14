package tests

import (
	"gomoney-mock-epl/web"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_admins_can_create_teams(t *testing.T) {
	dto := web.CreateTeamRequest{
		City:        "Manchester",
		HomeStadium: "Old Trafford",
		LogoURL:     "https://resources.premierleague.com/premierleague/badges/50/t1.png",
		Name:        "Manchester United",
		NameAbbr:    "MUN",
		ShortName:   "Man Utd",
	}
	req, rec := jsonRequest(http.MethodPost, "/teams/", &dto, adminToken)
	fixture.app.ServeHTTP(rec, req)

	result := rec.Result()
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	responseBody := web.DataDto{}
	assert.NoError(t, readJsonResponse(result.Body, &responseBody))
	assert.Equal(t, "Team", responseBody.Type)
}
