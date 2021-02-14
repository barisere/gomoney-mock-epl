package web

import (
	"fmt"

	"gomoney-mock-epl/config"
	"gomoney-mock-epl/database"
	"gomoney-mock-epl/teams"
	"gomoney-mock-epl/users"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
)

type RouteProvider func(*echo.Echo)

type Application struct {
	*config.Config
	DBClient  *mongo.Client
	DefaultDB *mongo.Database
	AdminDB   users.AdminsDB
	UsersDB   users.UsersDB
	TeamsDB   teams.TeamsDB
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

	e := echo.New()
	e.HTTPErrorHandler = DefaultErrorHandler
	e.Server.Addr = fmt.Sprintf(":%d", cfg.HttpBindPort)

	app := &Application{
		AdminDB:   adminsDB,
		Config:    &cfg,
		DBClient:  db,
		DefaultDB: defaultDB,
		Echo:      e,
		UsersDB:   usersDB,
		TeamsDB:   teamsDB,
	}

	AdminSignupRoute(app.AdminDB)(app.Echo)
	UserAuthRoute(app.UsersDB)(app.Echo)
	TeamRoutes(app.TeamsDB)(app.Echo)

	return app, nil
}
