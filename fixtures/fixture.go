package fixtures

import (
	"context"
	"gomoney-mock-epl/database"
	customErrors "gomoney-mock-epl/errors"
	"gomoney-mock-epl/teams"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Fixture struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	HomeTeam  *teams.Team        `json:"home_team" bson:"home_team"`
	AwayTeam  *teams.Team        `json:"away_team" bson:"away_team"`
	MatchDate time.Time          `json:"match_date" bson:"match_date"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type FixtureWriteModel struct {
	ID        primitive.ObjectID `bson:"_id"`
	HomeTeam  string             `bson:"home_team"`
	AwayTeam  string             `bson:"away_team"`
	MatchDate time.Time          `bson:"match_date"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type CreateFixtureRequest struct {
	HomeTeam  string    `json:"home_team"`
	AwayTeam  string    `json:"away_team"`
	MatchDate time.Time `json:"match_date"`
}

type FixturesDB struct {
	*mongo.Collection
	teams.TeamsDB
}

func (db FixturesDB) Create(ctx context.Context, dto CreateFixtureRequest) (*Fixture, error) {
	validationErrs := customErrors.ValidationError{
		Code:    "fixtures/cannot-create-fixture",
		Message: "Your request to create a fixture failed",
		Details: []customErrors.ValidationErrorDetails{},
	}
	if dto.HomeTeam == dto.AwayTeam {
		validationErrs.Message = "home team and away team must be different"
		return nil, validationErrs
	}
	homeTeam, err := db.TeamsDB.ByID(ctx, dto.HomeTeam)
	if err != nil {
		return nil, err
	}
	if homeTeam == nil {
		validationErrs.Details = append(validationErrs.Details, customErrors.ValidationErrorDetails{
			Field:   "home_team",
			Message: "Unknown home team",
		})
	}
	awayTeam, err := db.TeamsDB.ByID(ctx, dto.AwayTeam)
	if awayTeam == nil {
		validationErrs.Details = append(validationErrs.Details, customErrors.ValidationErrorDetails{
			Field:   "away_team",
			Message: "Unknown away team",
		})
	}
	if len(validationErrs.Details) > 0 {
		return nil, validationErrs
	}
	fixture := FixtureWriteModel{
		ID:        primitive.NewObjectID(),
		HomeTeam:  homeTeam.ID,
		AwayTeam:  awayTeam.ID,
		MatchDate: dto.MatchDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = db.InsertOne(ctx, fixture)
	return &Fixture{
		ID:        fixture.ID,
		HomeTeam:  homeTeam,
		AwayTeam:  awayTeam,
		MatchDate: fixture.MatchDate,
		CreatedAt: fixture.CreatedAt,
		UpdatedAt: fixture.UpdatedAt,
	}, err
}

func restFindStages() mongo.Pipeline {
	return mongo.Pipeline{
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: database.TeamsCollection},
				{Key: "localField", Value: "home_team"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "home_team"},
			}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: database.TeamsCollection},
			{Key: "localField", Value: "away_team"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "away_team"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$home_team"}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$away_team"}}}},
	}
}

func findFixtureQuery(fixtureID primitive.ObjectID) mongo.Pipeline {
	return append(mongo.Pipeline{bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "$eq", Value: fixtureID},
			}},
		}},
	}}, restFindStages()...)
}

func findFixturesQuery(fixtureIDs ...string) mongo.Pipeline {
	inIDMatch := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "$in", Value: fixtureIDs},
			}},
		}}
	return append(inIDMatch, restFindStages()...)
}

func listFixturesQuery() mongo.Pipeline {
	inIDMatch := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{}},
		}}
	return append(inIDMatch, restFindStages()...)
}

func listFixturesByStatusQuery(status fixtureStatus) mongo.Pipeline {
	comparison := "$gt"
	if status == Completed {
		comparison = "$lt"
	}
	now := time.Now().UTC()
	match := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "match_date", Value: bson.D{{Key: comparison, Value: now}}},
			}},
		},
	}
	return append(match, restFindStages()...)
}

type fixtureStatus string

const (
	Completed = fixtureStatus("completed")
	Pending   = fixtureStatus("pending")
)

func NewFixtureStatus(s string) fixtureStatus {
	if s == "completed" {
		return Completed
	}
	if s == "pending" {
		return Pending
	}
	return ""
}

func (db FixturesDB) List(ctx context.Context, status fixtureStatus) ([]Fixture, error) {
	query := listFixturesQuery()
	if status != "" {
		query = listFixturesByStatusQuery(status)
	}
	cursor, err := db.Collection.Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}
	fixtures := []Fixture{}
	if err := cursor.All(ctx, &fixtures); err != nil {
		return nil, err
	}
	return fixtures, nil
}

func (db FixturesDB) ByID(ctx context.Context, id primitive.ObjectID) (*Fixture, error) {
	cursor, err := db.Collection.Aggregate(ctx, findFixtureQuery(id))
	if err != nil {
		return nil, err
	}
	fixture := []Fixture{}
	if err := cursor.All(ctx, &fixture); err != nil {
		return nil, err
	}
	return &fixture[0], nil
}

func (db FixturesDB) Delete(ctx context.Context, id string) error {
	_, err := db.Collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	return err
}

func (db FixturesDB) Update(ctx context.Context, id primitive.ObjectID, update CreateFixtureRequest) (*Fixture, error) {
	fixture, err := db.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if fixture == nil {
		return nil, nil
	}
	writeModel := FixtureWriteModel{
		ID:        id,
		HomeTeam:  fixture.HomeTeam.ID,
		AwayTeam:  fixture.AwayTeam.ID,
		MatchDate: fixture.MatchDate,
		CreatedAt: fixture.CreatedAt,
		UpdatedAt: fixture.UpdatedAt,
	}
	if update.AwayTeam != "" {
		writeModel.AwayTeam = update.AwayTeam
	}
	if update.HomeTeam != "" {
		writeModel.HomeTeam = update.HomeTeam
	}
	if !update.MatchDate.IsZero() {
		writeModel.MatchDate = update.MatchDate
	}
	_, err = db.Collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: id}}, writeModel)

	return db.ByID(ctx, id)
}
