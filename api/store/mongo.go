package store

import (
	"context"
	"github.com/adrianosela/padl/api/user"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB holds the MongoDB Collection
type MongoDB struct {
	collection *mongo.Collection
}

// NewMongoDB initializes MongoDB connection
// returns MongoDB object
func NewMongoDB(dbUri string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(dbUri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")

	ds := &MongoDB{
		collection: client.Database("padl").Collection("Users"),
	}
	return ds, nil
}

// PutUser adds a new user to the database
func (db *MongoDB) PutUser(user *user.User) error {
	_, err := db.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil
}

// GetUser gets a user from the database
func (db *MongoDB) GetUser(email string) (*user.User, error) {
	query := bson.D{{"email", email}}

	var user user.User
	err := db.collection.FindOne(context.TODO(), query).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
