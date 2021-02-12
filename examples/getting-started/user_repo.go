// Code generated by repogen. DO NOT EDIT.
package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &UserRepositoryMongo{
		collection: collection,
	}
}

type UserRepositoryMongo struct {
	collection *mongo.Collection
}

func (r *UserRepositoryMongo) InsertOne(arg0 context.Context, arg1 *UserModel) (interface{}, error) {
	result, err := r.collection.InsertOne(arg0, arg1)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *UserRepositoryMongo) FindByUsername(arg0 context.Context, arg1 string) (*UserModel, error) {
	var entity UserModel
	if err := r.collection.FindOne(arg0, bson.M{
		"username": arg1,
	}).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *UserRepositoryMongo) UpdateDisplayNameByID(arg0 context.Context, arg1 string, arg2 primitive.ObjectID) (bool, error) {
	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg2,
	}, bson.M{
		"$set": bson.M{
			"display_name": arg1,
		},
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, err
}

func (r *UserRepositoryMongo) DeleteByCity(arg0 context.Context, arg1 string) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"city": arg1,
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}

func (r *UserRepositoryMongo) CountByCity(arg0 context.Context, arg1 string) (int, error) {
	count, err := r.collection.CountDocuments(arg0, bson.M{
		"city": arg1,
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}