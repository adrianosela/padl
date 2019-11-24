package store

import (
	"context"
	"log"

	"github.com/adrianosela/padl/api/project"
	"github.com/adrianosela/padl/api/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB holds the MongoDB Collection
type MongoDB struct {
	usersCollection    *mongo.Collection
	projectsCollection *mongo.Collection
}

// NewMongoDB initializes MongoDB connection
// returns MongoDB object
func NewMongoDB(connStr, dbName, usersCollName, projectsCollName string) (*MongoDB, error) {
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
		usersCollection:    client.Database(dbName).Collection(usersCollName),
		projectsCollection: client.Database(dbName).Collection(projectsCollName),
	}
	return ds, nil
}

// PutUser adds a new user to the database
func (db *MongoDB) PutUser(user *user.User) error {
	_, err := db.usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil
}

// GetUser gets a user from the database
func (db *MongoDB) GetUser(email string) (*user.User, error) {
	query := bson.D{{Key: "email", Value: email}}

	var user user.User
	err := db.usersCollection.FindOne(context.TODO(), query).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates a user in the database
func (db *MongoDB) UpdateUser(user *user.User) error {
	query := bson.D{{Key: "email", Value: user.Email}}

	update := bson.M{
		"$set": bson.M{
			"hashedpass": user.HashedPass,
			"keyid":      user.KeyID,
			"projects":   user.Projects,
		},
	}
	_, err := db.usersCollection.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return err
	}

	return nil
}

// UserExists returns true if a user with given email exists
func (db *MongoDB) UserExists(email string) (bool, error) {
	query := bson.D{{Key: "email", Value: email}}

	var user user.User
	err := db.usersCollection.FindOne(context.TODO(), query).Decode(&user)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// PutProject adds a new project to the database
func (db *MongoDB) PutProject(project *project.Project) error {
	_, err := db.projectsCollection.InsertOne(context.TODO(), project)
	if err != nil {
		return err
	}

	return nil
}

// GetProject gets a project from the database
func (db *MongoDB) GetProject(projectName string) (*project.Project, error) {
	query := bson.D{{Key: "name", Value: projectName}}

	var project project.Project
	err := db.projectsCollection.FindOne(context.TODO(), query).Decode(&project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// UpdateProject updates a project in the database
func (db *MongoDB) UpdateProject(project *project.Project) error {
	query := bson.D{{Key: "name", Value: project.Name}}

	update := bson.M{
		"$set": bson.M{
			"description":     project.Description,
			"members":         project.Members,
			"projectkey":      project.ProjectKey,
			"serviceAccounts": project.ServiceAccounts,
		},
	}
	_, err := db.projectsCollection.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProject deletes a project from the database
func (db *MongoDB) DeleteProject(projectName string) error {
	query := bson.D{{Key: "name", Value: projectName}}
	_, err := db.projectsCollection.DeleteOne(context.TODO(), query)
	if err != nil {
		return err
	}

	return nil
}

// ProjectExists returns true if a project with that name already exists
func (db *MongoDB) ProjectExists(projectName string) (bool, error) {
	query := bson.D{{Key: "name", Value: projectName}}

	var project project.Project
	err := db.projectsCollection.FindOne(context.TODO(), query).Decode(&project)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// ListProjects returns all projects in the db based on provided project names
func (db *MongoDB) ListProjects(projectNames []string) ([]*project.Project, error) {
	query := bson.M{
		"name": bson.M{
			"$in": projectNames,
		},
	}

	cur, err := db.projectsCollection.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}

	projects := []*project.Project{}
	for cur.Next(context.TODO()) {
		var elem project.Project
		err := cur.Decode(&elem)
		if err == nil {
			projects = append(projects, &elem)
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
