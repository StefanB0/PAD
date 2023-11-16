package service

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"padimage/database"
	"padimage/models"
	"strconv"

	"github.com/rs/zerolog/log"
)

type ImageService struct {
	db *database.ImageMongoDB
	as *AnalyticsService
	ts *TokenService

	transactions map[string]int
}

func NewImageService(db *database.ImageMongoDB, as *AnalyticsService, ts *TokenService) *ImageService {
	return &ImageService{
		db:           db,
		as:           as,
		ts:           ts,
		transactions: make(map[string]int),
	}
}

func (s *ImageService) GetImage(imageID int) (*models.Image, error) {
	return s.db.GetImage(imageID)
}

func (s *ImageService) CreateImage(image models.Image, token string) (int, error) {
	_, err := s.ts.VerifyAccessToken(token)
	if err != nil {
		log.Error().Err(err).Msg("Error verifying access token")
		return 0, err
	}

	id, err := s.db.CreateImage(image)
	if err != nil {
		log.Error().Err(err).Msg("Error creating image")
		return 0, err
	}

	err = s.as.AddImage(id, image.Tags)
	if err != nil {
		s.db.DeleteImage(int64(id))
		log.Error().Err(err).Msg("Error adding image to analytics")
		return 0, err
	}

	return id, err
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

func (s *ImageService) AddSagaTransaction(sagaID string, imageID int) {
	s.transactions[sagaID] = imageID
}

func (s *ImageService) RevertSagaTransaction(sagaID string) error {
	imageID := s.transactions[sagaID]
	return s.db.DeleteImage(int64(imageID))
}

func (s *ImageService) ConfirmSagaTransaction(sagaID string) error {
	gateway := os.Getenv("GATEWAY_ADDRESS")
	url := fmt.Sprintf("http://%s/transaction/%s", gateway, sagaID)
	payload := []byte(`{"status": "success", "service": "image_service"}`)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Error().Err(err).Msg("Error creating SAGA CONFIRM request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Error sending SAGA CONFIRM request")
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Error().Err(err).Msg("Error sending SAGA CONFIRM request. Status code: " + strconv.Itoa(res.StatusCode))
		return errors.New("Error sending SAGA CONFIRM request. Status code: " + strconv.Itoa(res.StatusCode))
	}

	return nil
}

func (s *ImageService) CancelSagaTransaction(sagaID string) error {
	gateway := os.Getenv("GATEWAY_ADDRESS")
	url := fmt.Sprintf("http://%s/transaction/%s", gateway, sagaID)
	payload := []byte(`{"status": "failure", "service": "image_service"}`)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Error().Err(err).Msg("Error creating SAGA CANCEL request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Error sending SAGA CANCEL request")
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Error().Err(err).Msg("Error sending SAGA CANCEL request. Status code: " + strconv.Itoa(res.StatusCode))
		return errors.New("Error sending SAGA CANCEL request. Status code: " + strconv.Itoa(res.StatusCode))
	}

	return nil
}
