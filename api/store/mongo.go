package store

import (
	"context"
	"github.com/adrianosela/padl/api/user"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Datastore holds the mongoDb Collection
type Datastore struct {
	collection *mongo.Collection
}

// Initialize initializes mongoDb connection
// returns datastore object
func Initialize(dbUri string) (*Datastore, error) {
	clientOptions := options.Client().ApplyURI(dbUri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")

	ds := &Datastore{
		collection: client.Database("test").Collection("Users"),
	}
	return ds, nil
}

// PutUser adds a new user to the database
func (ds *Datastore) PutUser(user *user.User) error {
	_, err := ds.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil
}

// GetUser gets a user from the database
func (ds *Datastore) GetUser(email string) (*user.User, error) {
	query := bson.D{{"email", email}}

	var user user.User
	err := ds.collection.FindOne(context.TODO(), query).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
