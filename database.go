package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	DATABASE = Database{"", "args", nil} // Global variable to access database
)

// Public mongodb database struct
type Database struct {
	MongoUrl, DatabaseName string
	MongoClient            *mongo.Client
}

// Start database connection
func (db *Database) StartConnection(database_uri string) error {
	db.MongoUrl = database_uri
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(db.MongoUrl))

	if err != nil {
		return err
	}

	// Ping database
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}

	db.MongoClient = client
	return nil
}

// Get mongo collection from database
func (db *Database) GetCollection(collection string) *mongo.Collection {
	return db.MongoClient.Database(db.DatabaseName).Collection(collection)
}
