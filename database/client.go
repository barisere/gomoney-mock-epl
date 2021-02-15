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
	MockEPLDatabase    = "mock_epl"
	AdminsCollection   = "admins"
	UsersCollection    = "users"
	TeamsCollection    = "teams"
	FixturesCollection = "fixtures"
)

func ConnectToDB(mongoURL string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

var unique = true
var uniqueAdminEmail = "unique_admin_emails"
var adminIndexModel = mongo.IndexModel{
	Keys: bson.D{{Key: "email", Value: 1}},
	Options: &options.IndexOptions{
		Name:   &uniqueAdminEmail,
		Unique: &unique,
	},
}

var uniqueUserEmails = "unique_user_emails"
var userIndexModel = mongo.IndexModel{
	Keys: bson.D{{Key: "email", Value: 1}},
	Options: &options.IndexOptions{
		Name:   &uniqueUserEmails,
		Unique: &unique,
	},
}

var defaultSearchLanguage = "english"
var teamsSearch = "teams_search"
var teamIndexModel = []mongo.IndexModel{
	{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: &options.IndexOptions{Unique: &unique},
	},
	{
		Keys:    bson.D{{Key: "short_name", Value: 1}},
		Options: &options.IndexOptions{Unique: &unique},
	},
	{
		Keys: bson.D{{Key: "name", Value: "text"},
			{Key: "short_name", Value: "text"},
			{Key: "city", Value: "text"}},
		Options: &options.IndexOptions{
			DefaultLanguage: &defaultSearchLanguage,
			Name:            &teamsSearch,
			Weights: bson.D{
				{Key: "name", Value: 4},
				{Key: "short_name", Value: 2},
				{Key: "city", Value: 1},
			},
		},
	},
}

var fixturesSearch = "fixtures_search"
var fixturesIndexModel = mongo.IndexModel{
	Keys: bson.D{
		{Key: "home_team_name", Value: "text"},
		{Key: "away_team_name", Value: "text"}},
	Options: &options.IndexOptions{
		DefaultLanguage: &defaultSearchLanguage,
		Name:            &fixturesSearch,
	},
}

func createIndexes(db *mongo.Database) error {
	ctx := context.Background()
	adminIndexes := db.Collection(AdminsCollection).Indexes()
	adminIndexes.DropAll(ctx)
	_, err := adminIndexes.CreateOne(ctx, adminIndexModel)
	if err != nil {
		return err
	}
	userIndexes := db.Collection(UsersCollection).Indexes()
	userIndexes.DropAll(ctx)
	_, err = userIndexes.CreateOne(context.Background(), userIndexModel)
	teamIndexes := db.Collection(TeamsCollection).Indexes()
	teamIndexes.DropAll(ctx)
	_, err = teamIndexes.CreateMany(ctx, teamIndexModel)
	if err != nil {
		return err
	}
	fixturesIndexes := db.Collection(FixturesCollection).Indexes()
	fixturesIndexes.DropAll(ctx)
	_, err = fixturesIndexes.CreateOne(ctx, fixturesIndexModel)
	if err != nil {
		return err
	}

	return nil
}
