package database

import (
	"context"
	"os"
	"padimage/models"
	"strconv"

	"github.com/joho/godotenv"
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
	godotenv.Load()
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")

	mongoURI := "mongodb://" + mongoUser + ":" + mongoPassword + "@" + mongoHost + ":" + mongoPort
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to MongoDB")
		log.Info().Msg(mongoURI)
		return nil
	}

	log.Info().Msg("Connected to MongoDB")

	userDB := client.Database("ImageDB")
	userCollection := userDB.Collection("images")

	return &ImageMongoDB{
		userDatabase:   userDB,
		userCollection: userCollection,
	}
}

func (db *ImageMongoDB) GetImage(imageID int) (*models.Image, error) {
	var image models.Image

	err := db.userCollection.FindOne(
		context.Background(),
		bson.D{
			{Key: "imageID", Value: imageID},
		},
	).Decode(&image)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving image from MongoDB. Image Id: " + strconv.Itoa(imageID))
		return nil, err
	}

	return &image, nil
}

func (db *ImageMongoDB) CreateImage(image models.Image) (int, error) {
	id, err := db.genNextID()
	if err != nil {
		log.Error().Err(err).Msg("Error generating next image ID")
		return 0, err
	}

	res, err := db.userCollection.InsertOne(
		context.Background(),
		bson.D{
			{Key: "imageID", Value: id},
			{Key: "author", Value: image.Author},
			{Key: "title", Value: image.Title},
			{Key: "description", Value: image.Description},
			{Key: "tags", Value: image.Tags},
			{Key: "imageChunk", Value: image.ImageChunk},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting image into MongoDB")
		return 0, err
	}

	err = db.userCollection.FindOne(context.Background(), bson.D{
		{Key: "_id", Value: res.InsertedID},
	}).Decode(&image)
	
	return image.ImageID, nil
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

func (db *ImageMongoDB) GetAuthorImages(author string) ([]int, error) {
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

	var imageIDs []int

	for _, image := range images {
		imageIDs = append(imageIDs, image.ImageID)
	}

	return imageIDs, nil
}

func (db *ImageMongoDB) GetTagImages(tag string) ([]int, error) {
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

	var imageIDs []int

	for _, image := range images {
		imageIDs = append(imageIDs, image.ImageID)
	}

	return imageIDs, nil
}

func (db *ImageMongoDB) AddViews(imageID int, views int) error {
	_, err := db.userCollection.UpdateOne(
		context.Background(),
		bson.D{
			{Key: "imageID", Value: imageID},
		},
		bson.D{
			{Key: "$inc", Value: bson.D{
				{Key: "views", Value: views},
			}},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error adding views to image in MongoDB")
		return err
	}

	return nil
}

func (db *ImageMongoDB) AddLikes(imageID int, likes int) error {
	_, err := db.userCollection.UpdateOne(
		context.Background(),
		bson.D{
			{Key: "imageID", Value: imageID},
		},
		bson.D{
			{Key: "$inc", Value: bson.D{
				{Key: "likes", Value: likes},
			}},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error adding likes to image in MongoDB")
		return err
	}

	return nil
}

func (db *ImageMongoDB) ModifyImage(imageID int, image models.Image) error {
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

func (db *ImageMongoDB) genNextID() (int, error) {
	var newestImage models.Image

	err := db.userCollection.FindOne(
		context.Background(),
		bson.D{},
		options.FindOne().SetSort(bson.D{{Key: "imageID", Value: -1}}),
	).Decode(&newestImage)

	if err != nil && err != mongo.ErrNoDocuments {
		log.Error().Err(err).Msg("Error retrieving newest image from MongoDB")
		return -1, err
	}

	return newestImage.ImageID + 1, nil
}