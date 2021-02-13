package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gomoney-mock-epl/config"
	"gomoney-mock-epl/database"
	"gomoney-mock-epl/web"

	"github.com/tylerb/graceful"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("server not configured properly: %v", err)
	}
	dbClient, err := database.ConnectToDB(config.MongoURL)
	if err != nil {
		log.Fatal(err)
	}
	app, err := web.NewApplication(dbClient, *config)
	if err != nil {
		log.Fatal(err)
	}
	defer app.DBClient.Disconnect(context.Background())

	web.AdminSignupRoute(*app.AdminDB)(app.Echo)

	if err := graceful.ListenAndServe(app.Echo.Server, 10*time.Second); err != nil {
		log.Fatalf("error: %v\n", err)
	}

	fmt.Println("Exiting...")
}
