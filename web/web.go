package web

import (
	"fmt"

	"gomoney-mock-epl/config"
	"gomoney-mock-epl/database"
	"gomoney-mock-epl/fixtures"
	"gomoney-mock-epl/teams"
	"gomoney-mock-epl/users"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type RouteProvider func(*echo.Echo)

type Application struct {
	*config.Config
	DBClient   *mongo.Client
	DefaultDB  *mongo.Database
	AdminDB    users.AdminsDB
	FixturesDB fixtures.DB
	UsersDB    users.UsersDB
	TeamsDB    teams.TeamsDB
	*echo.Echo
}

func NewApplication(db *mongo.Client, cfg config.Config) (*Application, error) {
	defaultDB := db.Database(database.MockEPLDatabase)
	adminsCollection := defaultDB.Collection(database.AdminsCollection)
	adminsDB := users.AdminsDB{Collection: adminsCollection}
	usersCollection := defaultDB.Collection(database.UsersCollection)
	usersDB := users.UsersDB{Collection: usersCollection}
	teamsCollection := defaultDB.Collection(database.TeamsCollection)
	teamsDB := teams.TeamsDB{Collection: teamsCollection}
	fixturesCollection := defaultDB.Collection(database.FixturesCollection)
	fixturesDB := fixtures.DB{Collection: fixturesCollection, TeamsDB: teamsDB}

	e := echo.New()
	e.Use(middleware.Logger(),
		middleware.Recover(),
		middleware.CORS(),
		middleware.BodyLimit("8K"))
	e.HTTPErrorHandler = DefaultErrorHandler
	e.Server.Addr = fmt.Sprintf("0.0.0.0:%d", cfg.HttpBindPort)

	app := &Application{
		AdminDB:    adminsDB,
		Config:     &cfg,
		DBClient:   db,
		DefaultDB:  defaultDB,
		Echo:       e,
		UsersDB:    usersDB,
		TeamsDB:    teamsDB,
		FixturesDB: fixturesDB,
	}

	AdminSignupRoute(app.AdminDB)(app.Echo)
	UserAuthRoute(app.UsersDB)(app.Echo)
	TeamRoutes(app.TeamsDB)(app.Echo)
	FixturesRoutes(app.FixturesDB)(app.Echo)
	searchRoutesProvider(app.TeamsDB, app.FixturesDB)(app.Echo)
	app.GET("/", func(c echo.Context) error {
		return c.File("docs/index.html")
	})
	app.GET("/openapi.yaml", func(c echo.Context) error {
		return c.File("docs/openapi.yaml")
	})

	return app, nil
}
