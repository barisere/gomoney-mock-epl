package tests

import (
	"context"
	"gomoney-mock-epl/fixtures"
	"gomoney-mock-epl/web"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_anyone_can_search_on_fixtures_and_teams(t *testing.T) {
	clearTeamsDB()
	clearFixtures()

	ctx := context.Background()
	lvpl, _ := testApp.app.TeamsDB.Create(ctx, liverpool.ToTeam(""))
	mct, _ := testApp.app.TeamsDB.Create(ctx, manCity.ToTeam(""))
	mutd, _ := testApp.app.TeamsDB.Create(ctx, manUtd.ToTeam(""))
	testApp.app.FixturesDB.Create(ctx, fixtures.CreateFixtureRequest{
		HomeTeam:  lvpl.ID,
		AwayTeam:  mct.ID,
		MatchDate: time.Now().Add(1 * time.Hour),
	})
	testApp.app.FixturesDB.Create(ctx, fixtures.CreateFixtureRequest{
		HomeTeam:  mutd.ID,
		AwayTeam:  lvpl.ID,
		MatchDate: time.Now().Add(24 * time.Hour),
	})
	testApp.app.FixturesDB.Create(ctx, fixtures.CreateFixtureRequest{
		HomeTeam:  mct.ID,
		AwayTeam:  mutd.ID,
		MatchDate: time.Now().Add(1 * time.Hour),
	})

	req, rec := jsonRequest(http.MethodGet, "/search?q=manchester%20united", nil, "")
	testApp.app.ServeHTTP(rec, req)
	result := rec.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	body := web.DataDto{}
	readJsonResponse(result.Body, &body)
	assert.Equal(t, "SearchResults", body.Type)
	assert.NotEmpty(t, body.Data)
}
