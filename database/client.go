package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MockEPLDatabase  = "mock_epl"
	AdminsCollection = "admins"
)

func ConnectToDB(mongoURL string) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	defaultDB := client.Database(MockEPLDatabase)
	return client, createIndexes(defaultDB)
}

func createIndexes(db *mongo.Database) error {
	uniqueAdminEmail := "unique_admin_emails"
	unique := true
	background := true
	_, err := db.Collection(AdminsCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: nil,
		Options: &options.IndexOptions{
			Background: &background,
			Name:       &uniqueAdminEmail,
			Unique:     &unique,
		},
	})
	return err
}
