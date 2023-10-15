package database

import (
	"context"
	"fmt"
	"os"
	"padimage/models"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserMongoDB struct {
	userDatabase   *mongo.Database
	userCollection *mongo.Collection
}

func NewMongoDB() *UserMongoDB {
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")

	mongoURI := "mongodb://" + mongoUser + ":" + mongoPassword + "@" + mongoHost + ":" + mongoPort
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to MongoDB")
		return nil
	}

	userDB := client.Database("UserDB")
	userCollection := userDB.Collection("Users")

	return &UserMongoDB{
		userDatabase:   userDB,
		userCollection: userCollection,
	}
}

func (db *UserMongoDB) GetUser(username string) (*models.User, error) {
	var user models.User

	err := db.userCollection.FindOne(
		context.Background(),
		bson.D{
			{Key: "username", Value: username},
		},
	).Decode(&user)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving user from MongoDB")
		return nil, err
	}

	return &user, nil
}

func (db *UserMongoDB) CreateUser(user models.User) error {
	result, err := db.userCollection.InsertOne(
		context.Background(),
		bson.D{
			{Key: "ID", Value: user.ID},
			{Key: "Username", Value: user.Username},
			{Key: "Password", Value: user.Password},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting user into MongoDB")
		return err
	}

	fmt.Printf("result: %v", result)

	return nil
}

func (db *UserMongoDB) DeleteUser(username string) error {
	_, err := db.userCollection.DeleteOne(
		context.Background(),
		bson.D{
			{Key: "username", Value: username},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting user from MongoDB")
		return err
	}

	return nil
}

func (db *UserMongoDB) GetAll() ([]models.User, error) {
	var users []models.User

	cursor, err := db.userCollection.Find(
		context.Background(),
		bson.D{},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving users from MongoDB")
		return nil, err
	}

	if err = cursor.All(context.Background(), &users); err != nil {
		log.Error().Err(err).Msg("Error decoding users from MongoDB")
	}

	return users, nil
}

func (db *UserMongoDB) DeleteAll() error {
	_, err := db.userCollection.DeleteMany(
		context.Background(),
		bson.D{},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting users from MongoDB")
		return err
	}

	return nil
}
