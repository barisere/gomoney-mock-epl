package web

import (
	"context"
	"errors"
	"fmt"
	"gomoney-mock-epl/fixtures"
	"gomoney-mock-epl/teams"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SearchResults collects search responses for a query.
type SearchResults struct {
	Query    string             `json:"query"`
	Fixtures []fixtures.Fixture `json:"fixures"`
	Teams    []teams.Team       `json:"teams"`
}

func searchTeams(ctx context.Context, db teams.TeamsDB, query string) ([]teams.Team, error) {
	q := bson.D{
		{Key: "$text", Value: bson.D{
			{Key: "$search", Value: query},
		}},
	}
	score := bson.D{
		{Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}}},
	}
	teamsCursor, err := db.Find(ctx, q, options.Find().SetProjection(score).SetSort(score))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	teams := []teams.Team{}
	if err := teamsCursor.All(ctx, &teams); err != nil {
		return nil, err
	}

	return teams, nil
}

// runSearch executes a text search on both the teams
// database and the fixtures database.
func runSearch(ctx context.Context, teamsDB teams.TeamsDB, fixturesDB fixtures.DB, query string) (*SearchResults, error) {
	teams, err := searchTeams(ctx, teamsDB, query)
	if err != nil {
		return nil, err
	}
	fixtures, err := fixturesDB.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	return &SearchResults{
		Query:    query,
		Fixtures: fixtures,
		Teams:    teams,
	}, nil
}

func searchRoutesProvider(teamsDB teams.TeamsDB, fixturesDB fixtures.DB) RouteProvider {
	return func(e *echo.Echo) {
		e.GET("/search", func(c echo.Context) error {
			query := c.QueryParam("q")
			results, err := runSearch(c.Request().Context(), teamsDB, fixturesDB, query)
			if err != nil {
				return err
			}
			message := query
			if len(message) > 20 {
				message = message[:20] + "..."
			}
			return c.JSON(http.StatusOK,
				dataResponse("SearchResults",
					fmt.Sprintf("Search results for %q", message), results))
		})
	}
}
