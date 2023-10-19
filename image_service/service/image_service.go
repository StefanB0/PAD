package service

import (
	"padimage/database"
	"padimage/models"

	"github.com/rs/zerolog/log"
)

type ImageService struct {
	db *database.ImageMongoDB
	as *AnalyticsService
	ts *TokenService
}

func NewImageService(db *database.ImageMongoDB, as *AnalyticsService, ts *TokenService) *ImageService {
	return &ImageService{
		db: db,
		as: as,
		ts: ts,
	}
}

func (s *ImageService) GetImage(imageID int) (*models.Image, error) {
	return s.db.GetImage(imageID)
}

func (s *ImageService) CreateImage(image models.Image, token string) error {
	_, err := s.ts.VerifyAccessToken(token)
	if err != nil {
		log.Error().Err(err).Msg("Error verifying access token")
		return err
	}

	s.as.AddImage(image.ImageID, image.Tags)

	return s.db.CreateImage(image)
}

func (s *ImageService) UpdateImage(image models.Image, token string) error {
	_, err := s.ts.VerifyAccessToken(token)
	if err != nil {
		log.Error().Err(err).Msg("Error verifying access token")
		return err
	}

	return s.db.ModifyImage(image.ImageID, image)
}

func (s *ImageService) DeleteImage(imageID int64, token string) error {
	_, err := s.ts.VerifyAccessToken(token)
	if err != nil {
		log.Error().Err(err).Msg("Error verifying access token")
		return err
	}

	return s.db.DeleteImage(imageID)
}

func (s *ImageService) AddViews(imageId, views int) error {
	s.as.AddEngagement(imageId, views, 0)
	return s.db.AddViews(imageId, views)
}

func (s *ImageService) AddLikes(imageId, likes int) error {
	s.as.AddEngagement(imageId, 0, likes)
	return s.db.AddLikes(imageId, likes)
}

func (s *ImageService) GetImagesByAuthor(author string) ([]int, error) {
	return s.db.GetAuthorImages(author)
}

func (s *ImageService) GetImagesByTag(tag string) ([]int, error) {
	return s.db.GetTagImages(tag)
}

func (s *ImageService) DeleteAllImages() error {
	return s.db.DeleteAll()
}
