package fixtures

import (
	"context"
	"errors"
	"gomoney-mock-epl/database"
	customErrors "gomoney-mock-epl/errors"
	"gomoney-mock-epl/teams"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Fixture is a match between two teams.
type Fixture struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	HomeTeam  *teams.Team        `json:"home_team" bson:"home_team"`
	AwayTeam  *teams.Team        `json:"away_team" bson:"away_team"`
	MatchDate time.Time          `json:"match_date" bson:"match_date"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// fixtureWriteModel defines the shape of the data we save to MongoDB.
// We store the home team name and the away teamn name in addition
// to their IDs. The names are stored to be used in text search only.
// Because the actual teams referenced can be updated, this information
// can get out of sync. It's an optimisation for the search because the
// text match stage has to be the first stage of the pipeline.
type fixtureWriteModel struct {
	ID           primitive.ObjectID `bson:"_id"`
	HomeTeam     string             `bson:"home_team"`
	HomeTeamName string             `bson:"home_team_name"`
	AwayTeam     string             `bson:"away_team"`
	AwayTeamName string             `bson:"away_team_name"`
	MatchDate    time.Time          `bson:"match_date"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

// CreateFixtureRequest is the DTO we receive from the
// clients when creating or updating fixtures.
type CreateFixtureRequest struct {
	HomeTeam  string    `json:"home_team"`
	AwayTeam  string    `json:"away_team"`
	MatchDate time.Time `json:"match_date"`
}

// DB provides methods for storing and accessing fixtures
// in the database. It uses the teams database for lookups.
type DB struct {
	*mongo.Collection
	teams.TeamsDB
}

// Create adds a new fixture to the system. The basic validations done
// is to ensure that the teams referenced actually exist, and that the
// same team is not paired with itself.
func (db DB) Create(ctx context.Context, dto CreateFixtureRequest) (*Fixture, error) {
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
	now := time.Now()
	fixture := fixtureWriteModel{
		ID:           primitive.NewObjectID(),
		HomeTeam:     homeTeam.ID,
		HomeTeamName: homeTeam.Name,
		AwayTeam:     awayTeam.ID,
		AwayTeamName: awayTeam.Name,
		MatchDate:    dto.MatchDate,
		CreatedAt:    now,
		UpdatedAt:    now,
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

func (db DB) List(ctx context.Context, status fixtureStatus) ([]Fixture, error) {
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

func textSearchQuery(q string) mongo.Pipeline {
	textMatch := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "$text", Value: bson.D{
					{Key: "$search", Value: q},
				}},
			}},
		},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "score", Value: bson.D{
				{Key: "$meta", Value: "textScore"},
			}},
		}}},
	}
	return append(textMatch, restFindStages()...)
}

func (db DB) Search(ctx context.Context, query string) ([]Fixture, error) {
	cursor, err := db.Aggregate(ctx, textSearchQuery(query))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	fixtures := []Fixture{}
	if err := cursor.All(ctx, &fixtures); err != nil {
		return nil, err
	}
	return fixtures, nil
}

func (db DB) ByID(ctx context.Context, id primitive.ObjectID) (*Fixture, error) {
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

func (db DB) Delete(ctx context.Context, id string) error {
	_, err := db.Collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	return err
}

func (db DB) Update(ctx context.Context, id primitive.ObjectID, dto CreateFixtureRequest) (*Fixture, error) {
	validationErrs := customErrors.ValidationError{
		Code:    "fixtures/cannot-create-fixture",
		Message: "Your request to create a fixture failed",
		Details: []customErrors.ValidationErrorDetails{},
	}
	if dto.HomeTeam == dto.AwayTeam && dto.HomeTeam != "" {
		validationErrs.Message = "home team and away team must be different"
		return nil, validationErrs
	}
	fixture, err := db.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if fixture == nil {
		return nil, nil
	}
	writeModel := fixtureWriteModel{
		ID:           id,
		HomeTeam:     fixture.HomeTeam.ID,
		HomeTeamName: fixture.HomeTeam.Name,
		AwayTeam:     fixture.AwayTeam.ID,
		AwayTeamName: fixture.AwayTeam.Name,
		MatchDate:    fixture.MatchDate,
		CreatedAt:    fixture.CreatedAt,
		UpdatedAt:    time.Now(),
	}
	if dto.HomeTeam != "" {
		homeTeam, err := db.TeamsDB.ByID(ctx, dto.HomeTeam)
		if err != nil {
			return nil, err
		}
		if homeTeam == nil {
			validationErrs.Details = append(validationErrs.Details, customErrors.ValidationErrorDetails{
				Field:   "home_team",
				Message: "Unknown home team",
			})
		} else {
			writeModel.HomeTeam = dto.HomeTeam
			writeModel.HomeTeamName = homeTeam.Name
		}
	}
	if dto.AwayTeam != "" {
		awayTeam, err := db.TeamsDB.ByID(ctx, dto.AwayTeam)
		if err != nil {
			return nil, err
		}
		if awayTeam == nil {
			validationErrs.Details = append(validationErrs.Details, customErrors.ValidationErrorDetails{
				Field:   "away_team",
				Message: "Unknown away team",
			})
		} else {
			writeModel.AwayTeam = dto.AwayTeam
			writeModel.AwayTeamName = awayTeam.Name
		}
	}
	if len(validationErrs.Details) > 0 {
		return nil, validationErrs
	}
	if !dto.MatchDate.IsZero() {
		writeModel.MatchDate = dto.MatchDate
	}
	_, err = db.Collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: id}}, writeModel)

	return db.ByID(ctx, id)
}
