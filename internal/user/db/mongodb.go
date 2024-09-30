package db

import (
	"artOfDevPractise/internal/apperror"
	"artOfDevPractise/internal/user"
	"artOfDevPractise/pkg/logging"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("error inserting user: %v", err)
	}

	d.logger.Debug("convert insertedID to objectiD")
	oid, ok := result.InsertedID.(bson.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectid to HEX")
}

func (d db) FindOne(ctx context.Context, id string) (u user.User, err error) {

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert HEX to objectid. Hex: %s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {

			return u, apperror.ErrorNotFound
		}
		return u, fmt.Errorf("failed to find user by id: %s, error: %v", id, result.Err())
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user (%s) frob DB error: %v", id, err)
	}
	return u, nil
}

func (d db) Update(ctx context.Context, user user.User) error {
	objectID, err := bson.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert HEX to objectid. Hex: %s", user.ID)
	}

	filter := bson.M{"_id": objectID}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user. error: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user. error: %v", err)
	}

	delete(updateUserObj, "_id")
	update := bson.M{"$set": updateUserObj}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user by id: %s, error: %v", objectID, err)
	}

	if result.MatchedCount == 0 {
		return apperror.ErrorNotFound
	}

	d.logger.Tracef("Matched: %d, documents and Modifide: %d", result.MatchedCount, result.ModifiedCount)
	return nil

}

func (d db) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert HEX to objectid. Hex: %s", id)
	}

	filter := bson.M{"_id": objectID}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user by id: %s, error: %v", objectID, err)
	}
	if result.DeletedCount == 0 {
		return apperror.ErrorNotFound
	}
	d.logger.Tracef("Deleted: %d documents", result.DeletedCount)
	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {

	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
