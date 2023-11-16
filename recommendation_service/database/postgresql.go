package database

import (
	"os"
	"padrecommendations/models"
	"padrecommendations/utils"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnalyticsPostgresDB struct {
	imageDB *gorm.DB
}

// Create new user repository
func NewPostgresDatabase() *AnalyticsPostgresDB {

	godotenv.Load()

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_ADDRESS") // "localhost"
	port := os.Getenv("POSTGRES_PORT")    // ":5432"

	dbURL := "postgres://" + user + ":" + password + "@" + host + port + "/" + dbName

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Error().Err(err).Msg("Error connecting to Postgres")
		return nil
	}

	db.AutoMigrate(&models.TagSQL{})
	db.AutoMigrate(&models.ImageSQL{})

	log.Info().Msg("Connected to Postgres")
	return &AnalyticsPostgresDB{
		imageDB: db,
	}
}

func (db *AnalyticsPostgresDB) GetTags() []string {
	var tagnames []string
	var tags []models.TagSQL

	result := db.imageDB.Find(&tags)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error retrieving taglists from Postgres")
		return tagnames
	}

	for _, tag := range tags {
		tagnames = append(tagnames, tag.Tagname)
	}

	return tagnames
}

func (db *AnalyticsPostgresDB) GetTagEngagement(tag string) (int, error) {
	var taglist models.TagSQL

	result := db.imageDB.First(&taglist, "tagname = ?", tag)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error retrieving taglist from Postgres")
		return 0, result.Error
	}

	return taglist.TotalEngagement, nil
}

func (db *AnalyticsPostgresDB) GetImageList(tag string) ([]models.Image, error) {
	var images []models.Image
	var sqlimages []models.ImageSQL

	err := db.imageDB.Preload("Tags", "tagname = ?", tag).Find(&sqlimages).Error
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving images from Postgres")
		return nil, err
	}

	for _, sqlimage := range sqlimages {
		images = append(images, models.Image{
			ImageID:    sqlimage.ImageID,
			Views:      sqlimage.Views,
			Likes:      sqlimage.Likes,
			Engagement: sqlimage.Engagement,
		})
	}
	return images, nil
}

func (db *AnalyticsPostgresDB) GetImage(imageID int) (*models.Image, error) {
	var imagesql models.ImageSQL

	result := db.imageDB.First(&imagesql, "image_id = ?", imageID)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error retrieving image from Postgres")
		return nil, result.Error
	}

	image := models.Image{
		ImageID:    imagesql.ImageID,
		Views:      imagesql.Views,
		Likes:      imagesql.Likes,
		Engagement: imagesql.Engagement,
	}

	return &image, nil
}

func (db *AnalyticsPostgresDB) AddImage(image models.Image) error {
	sqlImage := models.ImageSQL{
		ImageID:    image.ImageID,
		Views:      image.Views,
		Likes:      image.Likes,
		Engagement: image.Engagement,
	}

	result := db.imageDB.Create(&sqlImage)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error creating image in Postgres")
		return result.Error
	}

	existingTags := db.GetTags()

	for _, tag := range image.Tags {
		if !utils.ContainsString(existingTags, tag) {
			db.imageDB.Create(&models.TagSQL{
				Tagname:         tag,
				TotalEngagement: image.Engagement,
			})
		} else {
			var taglist models.TagSQL

			result := db.imageDB.First(&taglist, "tagname = ?", tag)

			if result.Error != nil {
				log.Error().Err(result.Error).Msg("Error retrieving taglist from Postgres")
				return result.Error
			}

			taglist.TotalEngagement += image.Engagement

			result = db.imageDB.Save(&taglist)

			if result.Error != nil {
				log.Error().Err(result.Error).Msg("Error updating taglist in Postgres")
				return result.Error
			}

		}
	}

	return nil
}

func (db *AnalyticsPostgresDB) AddViews(imageID int, views int, engagement int) error {
	var image models.ImageSQL

	result := db.imageDB.Preload("Tags").First(&image, "image_id = ?", imageID)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error retrieving image from Postgres")
		return result.Error
	}

	image.Views += views
	image.Engagement += engagement

	result = db.imageDB.Save(&image)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error updating image in Postgres")
		return result.Error
	}

	tags := image.Tags

	for _, tag := range tags {
		tag.TotalEngagement += engagement
		db.imageDB.Save(&tag)
	}

	return nil
}

func (db *AnalyticsPostgresDB) AddLikes(imageID int, likes int, engagement int) error {
	var image models.ImageSQL
	result := db.imageDB.Preload("Tags").First(&image, "image_id = ?", imageID)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error retrieving image from Postgres")
		return result.Error
	}

	image.Likes += likes
	image.Engagement += engagement

	result = db.imageDB.Save(&image)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error updating image in Postgres")
		return result.Error
	}

	tags := image.Tags

	for _, tag := range tags {
		tag.TotalEngagement += engagement
		db.imageDB.Save(&tag)
	}

	return nil
}

func (db *AnalyticsPostgresDB) DeleteImage(imageID int) error {
	var image models.ImageSQL

	result := db.imageDB.First(&image, "image_id = ?", imageID)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error retrieving image from Postgres")
		return result.Error
	}

	for _, tagname := range image.Tags {
		var tag models.TagSQL
		result := db.imageDB.First(&tag, "tagname = ?", tagname)
		if result.Error != nil {
			log.Error().Err(result.Error).Msg("Error retrieving tag from Postgres")
			return result.Error
		}

		tag.TotalEngagement -= image.Engagement
		db.imageDB.Save(&tagname)
	}

	result = db.imageDB.Delete(&models.ImageSQL{}, "image_id = ?", imageID)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error deleting image from Postgres")
		return result.Error
	}

	return nil
}

func (db *AnalyticsPostgresDB) DeleteAll() error {
	result := db.imageDB.Where("1 = 1").Delete(&models.ImageSQL{})
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error deleting images from Postgres")
		return result.Error
	}

	result = db.imageDB.Where("1 = 1").Delete(&models.TagSQL{})
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error deleting tags from Postgres")
		return result.Error
	}

	return nil
}
