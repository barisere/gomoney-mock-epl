package teams

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Team struct {
	ID          string    `json:"id" bson:"_id"`
	City        string    `json:"city" bson:"city"`
	HomeStadium string    `json:"home_stadium" bson:"home_stadium"`
	LogoURL     string    `json:"logo_url" bson:"logo_url"`
	Name        string    `json:"name" bson:"name"`
	NameAbbr    string    `json:"name_abbr" bson:"name_abbr"`
	ShortName   string    `json:"short_name" bson:"short_name"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type TeamsDB struct {
	*mongo.Collection
}

// Create adds a new team to the database.
func (t TeamsDB) Create(ctx context.Context, team Team) (*Team, error) {
	team.ID = primitive.NewObjectID().Hex()
	team.CreatedAt = time.Now()
	team.UpdatedAt = team.CreatedAt
	_, err := t.InsertOne(ctx, &team, options.InsertOne().SetBypassDocumentValidation(false))
	return &team, err
}

// Update changes a team's information in the database.
func (t TeamsDB) Update(ctx context.Context, team Team) (*Team, error) {
	team.UpdatedAt = time.Now()
	filter := bson.D{bson.E{Key: "_id", Value: team.ID}}
	_, err := t.UpdateOne(ctx, filter, &team, options.Update().SetBypassDocumentValidation(false))
	return &team, err
}

// List fetches all the teams in the database. It's currently
// not paginated.
func (t TeamsDB) List(ctx context.Context) ([]Team, error) {
	cursor, err := t.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	teams := []Team{}
	if err := cursor.All(ctx, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

// Delete removes a team from the database.
func (t TeamsDB) Delete(ctx context.Context, id string) error {
	_, err := t.DeleteOne(ctx, bson.E{Key: "_id", Value: id})
	return err
}
