package web

import (
	"fmt"

	"gomoney-mock-epl/config"
	"gomoney-mock-epl/database"
	"gomoney-mock-epl/users"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
)

type RouteProvider func(*echo.Echo)

type Application struct {
	*config.Config
	DBClient  *mongo.Client
	DefaultDB *mongo.Database
	AdminDB   *users.AdminsDB
	*echo.Echo
}

func NewApplication(db *mongo.Client, cfg config.Config) (*Application, error) {
	defaultDB := db.Database(database.MockEPLDatabase)
	adminsCollection := defaultDB.Collection("admin_accounts")
	adminsDB := users.AdminsDB{Collection: adminsCollection}

	e := echo.New()
	e.HTTPErrorHandler = DefaultErrorHandler
	e.Server.Addr = fmt.Sprintf(":%d", cfg.HttpBindPort)

	return &Application{
		AdminDB:   &adminsDB,
		Config:    &cfg,
		DBClient:  db,
		DefaultDB: defaultDB,
		Echo:      e,
	}, nil
}
