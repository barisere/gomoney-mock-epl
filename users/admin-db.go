package users

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

type AdminsDB struct {
	*mongo.Collection
}

func (db AdminsDB) Create(ctx context.Context, admin Administrator) (*Administrator, error) {
	admin.ID = primitive.NewObjectID().Hex()
	_, err := db.InsertOne(ctx, &admin, options.InsertOne().SetBypassDocumentValidation(false))
	return &admin, err
}

func (db AdminsDB) ByEmail(ctx context.Context, email string) (*Administrator, error) {
	admin := Administrator{}
	filter := bson.D{{"email", email}}

	if err := db.FindOne(ctx, filter).Decode(&admin); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
}

func (db AdminsDB) ByID(ctx context.Context, ID string) (*Administrator, error) {
	admin := Administrator{}
	filter := bson.D{{"_id", ID}}

	if err := db.FindOne(ctx, filter).Decode(&admin); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
}
