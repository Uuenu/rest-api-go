package db

import (
	"context"
	"errors"
	"fmt"
	"rest-api-go/internal/user"
	"rest-api-go/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error %v", err)
	}

	d.logger.Debug("convert insertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to conver oid to hex oid: %s", oid)
}
func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to ObjectID. hex: %s", id)
	}
	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter, nil)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			//TODO ErrEntityNotFound
			return u, fmt.Errorf("not found")
		}
		return u, fmt.Errorf("failed to FindOne user by id: %s due to error: %v", id, err)
	}
	if err := result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user (id: %s) due to error: %v", id, err)
	}

	return u, nil
}

func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return u, fmt.Errorf("failed to find all users due error: %v", err)
	}

	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents from cursor. error: %v", err)
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	ObjectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to create ObjectID from hex user.ID. error: %v", err)
	}

	filter := bson.M{"_id": ObjectID}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user. error: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal userBytes. error: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query. error: %v", err)
	}
	if result.MatchedCount == 0 {
		//TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}

	d.logger.Tracef("Matched %d documents and Modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}
func (d *db) Delete(ctx context.Context, id string) error {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert ObjectID from hex(id: %s). error: %v", id, err)
	}

	filter := bson.M{"_id": ObjectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute delet user by id. error: %v", err)
	}

	if result.DeletedCount == 0 {
		//TODO ErrEntityNotFound
		return fmt.Errorf("not found document for deleted. error: %v", err)
	}

	d.logger.Tracef("Deleted %d documents", result.DeletedCount)

	return nil

}
