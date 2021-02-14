package users

import (
	"context"
	"errors"
	"gomoney-mock-epl/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

type UsersDB struct {
	*mongo.Collection
}

func (db UsersDB) Create(ctx context.Context, user User) (*User, error) {
	user.ID = primitive.NewObjectID().Hex()
	_, err := db.InsertOne(ctx, &user, options.InsertOne().SetBypassDocumentValidation(false))
	if database.IsDuplicateKeyError(err) {
		return nil, ErrEmailTaken
	}
	return &user, err
}

func (db UsersDB) ByEmail(ctx context.Context, email string) (*User, error) {
	user := User{}
	filter := bson.D{bson.E{Key: "email", Value: email}}

	if err := db.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (db UsersDB) ByID(ctx context.Context, ID string) (*User, error) {
	user := User{}
	filter := bson.D{bson.E{Key: "_id", Value: ID}}

	if err := db.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
