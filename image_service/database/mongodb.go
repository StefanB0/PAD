package database

import (
	"context"
	"os"
	"padimage/models"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ImageMongoDB struct {
	userDatabase   *mongo.Database
	userCollection *mongo.Collection
}

func NewMongoDB() *ImageMongoDB {
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

	userDB := client.Database("ImageDB")
	userCollection := userDB.Collection("images")

	return &ImageMongoDB{
		userDatabase:   userDB,
		userCollection: userCollection,
	}
}

func (db *ImageMongoDB) GetImage(imageID int64) (*models.Image, error) {
	var image models.Image

	err := db.userCollection.FindOne(
		context.Background(),
		bson.D{
			{Key: "imageID", Value: imageID},
		},
	).Decode(&image)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving image from MongoDB")
		return nil, err
	}

	return &image, nil
}

func (db *ImageMongoDB) CreateImage(image models.Image) error {
	_, err := db.userCollection.InsertOne(
		context.Background(),
		bson.D{
			{Key: "imageID", Value: image.ImageID},
			{Key: "author", Value: image.Author},
			{Key: "title", Value: image.Title},
			{Key: "description", Value: image.Description},
			{Key: "tags", Value: image.Tags},
			{Key: "imageChunk", Value: image.ImageChunk},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting image into MongoDB")
		return err
	}

	return nil
}

func (db *ImageMongoDB) DeleteImage(imageID int64) error {
	_, err := db.userCollection.DeleteOne(
		context.Background(),
		bson.D{
			{Key: "imageID", Value: imageID},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting image from MongoDB")
		return err
	}

	return nil
}

func (db *ImageMongoDB) GetAuthorImages(author string) ([]int64, error) {
	var images []models.Image

	cursor, err := db.userCollection.Find(
		context.Background(),
		bson.D{
			{Key: "author", Value: author},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving images from MongoDB")
		return nil, err
	}

	if err = cursor.All(context.Background(), &images); err != nil {
		log.Error().Err(err).Msg("Error decoding images from MongoDB")
	}

	var imageIDs []int64

	for _, image := range images {
		imageIDs = append(imageIDs, image.ImageID)
	}

	return imageIDs, nil
}

func (db *ImageMongoDB) GetTagImages(tag string) ([]int64, error) {
	var images []models.Image

	cursor, err := db.userCollection.Find(
		context.Background(),
		bson.D{
			{Key: "tags", Value: tag},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving images from MongoDB")
		return nil, err
	}

	if err = cursor.All(context.Background(), &images); err != nil {
		log.Error().Err(err).Msg("Error decoding images from MongoDB")
	}

	var imageIDs []int64

	for _, image := range images {
		imageIDs = append(imageIDs, image.ImageID)
	}

	return imageIDs, nil
}

func (db *ImageMongoDB) ModifyImage(imageID int64, image models.Image) error {
	_, err := db.userCollection.UpdateOne(
		context.Background(),
		bson.D{
			{Key: "imageID", Value: imageID},
		},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "author", Value: image.Author},
				{Key: "title", Value: image.Title},
				{Key: "description", Value: image.Description},
				{Key: "tags", Value: image.Tags},
				{Key: "imageChunk", Value: image.ImageChunk},
			}},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error modifying image in MongoDB")
		return err
	}

	return nil
}

func (db *ImageMongoDB) DeleteAll() error {
	_, err := db.userCollection.DeleteMany(
		context.Background(),
		bson.D{},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting images from MongoDB")
		return err
	}

	return nil
}

// func (db *UserMongoDB) GetAll() ([]models.Image, error) {
// 	var images []models.Image

// 	cursor, err := db.userCollection.Find(
// 		context.Background(),
// 		bson.D{},
// 	)

// 	if err != nil {
// 		log.Error().Err(err).Msg("Error retrieving images from MongoDB")
// 		return nil, err
// 	}

// 	if err = cursor.All(context.Background(), &images); err != nil {
// 		log.Error().Err(err).Msg("Error decoding images from MongoDB")
// 	}

// 	return images, nil
// }

// func (db *UserMongoDB) GetUser(username string) (*models.Image, error) {
// 	var user models.Image

// 	err := db.userCollection.FindOne(
// 		context.Background(),
// 		bson.D{
// 			{Key: "username", Value: username},
// 		},
// 	).Decode(&user)
// 	if err != nil {
// 		log.Error().Err(err).Msg("Error retrieving user from MongoDB")
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func (db *UserMongoDB) CreateUser(user models.Image) error {
// 	result, err := db.userCollection.InsertOne(
// 		context.Background(),
// 		bson.D{
// 			{Key: "ID", Value: user.ID},
// 			{Key: "Username", Value: user.Username},
// 			{Key: "Password", Value: user.Password},
// 		},
// 	)
// 	if err != nil {
// 		log.Error().Err(err).Msg("Error inserting user into MongoDB")
// 		return err
// 	}

// 	fmt.Printf("result: %v", result)

// 	return nil
// }

// func (db *UserMongoDB) DeleteUser(username string) error {
// 	_, err := db.userCollection.DeleteOne(
// 		context.Background(),
// 		bson.D{
// 			{Key: "username", Value: username},
// 		},
// 	)
// 	if err != nil {
// 		log.Error().Err(err).Msg("Error deleting user from MongoDB")
// 		return err
// 	}

// 	return nil
// }

// func (db *UserMongoDB) GetAll() ([]models.Image, error) {
// 	var users []models.Image

// 	cursor, err := db.userCollection.Find(
// 		context.Background(),
// 		bson.D{},
// 	)
// 	if err != nil {
// 		log.Error().Err(err).Msg("Error retrieving users from MongoDB")
// 		return nil, err
// 	}

// 	if err = cursor.All(context.Background(), &users); err != nil {
// 		log.Error().Err(err).Msg("Error decoding users from MongoDB")
// 	}

// 	return users, nil
// }
