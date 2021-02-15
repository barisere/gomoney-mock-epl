package web

import (
	"fmt"
	"gomoney-mock-epl/fixtures"
	"net/http"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createFixture(db fixtures.FixturesDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		dto := fixtures.CreateFixtureRequest{}
		if err := c.Bind(&dto); err != nil {
			return err
		}
		fixture, err := db.Create(c.Request().Context(), dto)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated,
			dataResponse("Fixture", "Fixture created successfully", fixture))
	}
}

func listFixtures(db fixtures.FixturesDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		fixtures, err := db.List(c.Request().Context(),
			fixtures.NewFixtureStatus(c.QueryParam("status")))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK,
			dataResponse("Fixtures", "Available EPL fixtures", fixtures))
	}
}

func deleteFixture(db fixtures.FixturesDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		fixtureID := c.Param("fixture_id")
		if err := db.Delete(c.Request().Context(), fixtureID); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, nil)
	}
}

var fixtureNotFound = errorDto("NotFound", "That fixture does not exist")

func editFixture(db fixtures.FixturesDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		fixtureID, err := primitive.ObjectIDFromHex(c.Param("fixture_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fixtureNotFound)
		}
		dto := fixtures.CreateFixtureRequest{}
		if err := c.Bind(&dto); err != nil {
			return err
		}
		fixture, err := db.Update(c.Request().Context(), fixtureID, dto)
		if err != nil {
			return err
		}
		if fixture == nil {
			return echo.NewHTTPError(http.StatusNotFound, fixtureNotFound)
		}
		return c.JSON(http.StatusOK,
			dataResponse("Fixture", "Fixture updated successfully", fixture))
	}
}

func viewFixture(db fixtures.FixturesDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		fixtureID, err := primitive.ObjectIDFromHex(c.Param("fixture_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fixtureNotFound)
		}
		fixture, err := db.ByID(c.Request().Context(), fixtureID)
		if err != nil {
			return err
		}
		if fixture == nil {
			return echo.NewHTTPError(http.StatusNotFound,
				errorDto("NotFound", "That fixture does not exist"))
		}
		return c.JSON(http.StatusOK,
			dataResponse("Fixture", fmt.Sprintf("%s - %s",
				fixture.HomeTeam.ShortName, fixture.AwayTeam.ShortName), fixture))
	}
}

func FixturesRoutes(db fixtures.FixturesDB) RouteProvider {
	return func(e *echo.Echo) {
		fixturesRoutes := e.Group("/fixtures", jwtMiddleware)
		fixturesRoutes.POST("/", createFixture(db), onlyAdmins)
		fixturesRoutes.GET("/", listFixtures(db))
		fixturesRoutes.DELETE("/:fixture_id", deleteFixture(db), onlyAdmins)
		fixturesRoutes.GET("/:fixture_id", viewFixture(db))
		fixturesRoutes.PATCH("/:fixture_id", editFixture(db), onlyAdmins)
	}
}
