package keystore

import (
	"context"
	"log"
	"strings"

	"github.com/adrianosela/padl/api/kms"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB holds the MongoDB Collection
type MongoDB struct {
	privKeysCollection *mongo.Collection
	pubKeysCollection  *mongo.Collection
}

// NewMongoDB initializes MongoDB connection
// returns MongoDB object
func NewMongoDB(connStr, dbName, privKeysCollName, pubKeysCollName string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(connStr)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	log.Println("[info] successfully connected to MongoDB")

	ds := &MongoDB{
		privKeysCollection: client.Database(dbName).Collection(privKeysCollName),
		pubKeysCollection:  client.Database(dbName).Collection(pubKeysCollName),
	}
	return ds, nil
}

// PutPrivKey adds a new private key to the database
func (db *MongoDB) PutPrivKey(key *kms.PrivateKey) error {
	_, err := db.privKeysCollection.InsertOne(context.TODO(), key)
	if err != nil {
		if strings.LastIndex(err.Error(), "multiple write errors:") != -1 {
			return ErrKeyExists
		}
		return err
	}

	return nil
}

// GetPrivKey gets a private key from the database
func (db *MongoDB) GetPrivKey(id string) (*kms.PrivateKey, error) {
	query := bson.D{{"id", id}}

	var key kms.PrivateKey
	err := db.privKeysCollection.FindOne(context.TODO(), query).Decode(&key)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, ErrKeyNotFound
		}

		return nil, err
	}

	return &key, nil
}

// UpdatePrivKey updates a private key in the database
func (db *MongoDB) UpdatePrivKey(key *kms.PrivateKey) error {
	query := bson.D{{"id", key.ID}}

	update := bson.M{
		"$set": bson.M{
			"project": key.Project,
			"pem":     key.PEM,
		},
	}
	res, err := db.privKeysCollection.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrKeyNotFound
	}

	return nil
}

// DeletePrivKey deletes a private key from the database
func (db *MongoDB) DeletePrivKey(id string) error {
	query := bson.D{{"id", id}}
	res, err := db.privKeysCollection.DeleteOne(context.TODO(), query)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrKeyNotFound
	}

	return nil
}

// PutPubKey adds a new public key to the database
func (db *MongoDB) PutPubKey(key *kms.PublicKey) error {
	_, err := db.pubKeysCollection.InsertOne(context.TODO(), key)
	if err != nil {
		if strings.LastIndex(err.Error(), "multiple write errors:") != -1 {
			return ErrKeyExists
		}

		return err
	}

	return nil
}

//
func (db *MongoDB) GetPubKey(id string) (*kms.PublicKey, error) {
	query := bson.D{{"id", id}}

	var key kms.PublicKey
	err := db.pubKeysCollection.FindOne(context.TODO(), query).Decode(&key)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, ErrKeyNotFound
		}

		return nil, err
	}

	return &key, nil
}

// DeletePubKey deletes a public key from the database
func (db *MongoDB) DeletePubKey(id string) error {
	query := bson.D{{"id", id}}
	res, err := db.pubKeysCollection.DeleteOne(context.TODO(), query)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrKeyNotFound
	}

	return nil
}
