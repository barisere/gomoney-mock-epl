package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MockEPLDatabase  = "mock_epl"
	AdminsCollection = "admins"
	UsersCollection  = "users"
	TeamsCollection  = "teams"
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
	uniqueUserEmails := "unique_user_emails"
	unique := true
	ctx := context.Background()
	adminIndexes := db.Collection(AdminsCollection).Indexes()
	adminIndexes.DropAll(ctx)
	_, err := adminIndexes.CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
		Options: &options.IndexOptions{
			Name:   &uniqueAdminEmail,
			Unique: &unique,
		},
	})
	if err != nil {
		return err
	}
	userIndexes := db.Collection(UsersCollection).Indexes()
	userIndexes.DropAll(ctx)
	_, err = userIndexes.CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
		Options: &options.IndexOptions{
			Name:   &uniqueUserEmails,
			Unique: &unique,
		},
	})

	return err
}
