package service

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"padrecommendations/database"
	"padrecommendations/models"
	"strconv"

	"github.com/rs/zerolog/log"
)

type RecommendationService struct {
	analyticsDB *database.AnalyticsPostgresDB

	transactionList map[string]int
}

func NewRecommendationService(analyticsDB *database.AnalyticsPostgresDB) *RecommendationService {
	return &RecommendationService{
		analyticsDB:     analyticsDB,
		transactionList: make(map[string]int),
	}
}

func (s *RecommendationService) GetTags() []string {
	return s.analyticsDB.GetTags()
}

func (s *RecommendationService) GetRecommendations(tagname string) (int, error) {
	imagelist, err := s.analyticsDB.GetImageList(tagname)
	if err != nil {
		return 0, err
	}

	totalEngagement, err := s.analyticsDB.GetTagEngagement(tagname)
	if err != nil {
		return 0, err
	}

	randNum := rand.Intn(totalEngagement) + 1
	randImage := 0

	var counter int

	for _, image := range imagelist {
		counter += image.Engagement
		if counter >= randNum {
			randImage = image.ImageID
			break
		}
	}

	fmt.Println("randImage: ", randImage)
	fmt.Println("totalEngagement: ", totalEngagement)
	for _, image := range imagelist {
		fmt.Println("image: ", image.ImageID, "engagement: ", image.Engagement)
	}

	return randImage, nil
}

func (s *RecommendationService) AddImage(image models.Image) error {
	image.Engagement = image.Views + image.Likes
	image.Engagement += 100
	err := s.analyticsDB.AddImage(image)
	if err != nil {
		log.Err(err).Msg("Error adding image")
		return err
	}

	return nil
}

func (service *RecommendationService) AddView(imageID int, views int) {
	_, err := service.analyticsDB.GetImage(imageID)
	if err != nil {
		return
	}

	service.analyticsDB.AddViews(imageID, views, views)
}

func (s *RecommendationService) AddLike(imageID int, likes int) {
	_, err := s.analyticsDB.GetImage(imageID)
	if err != nil {
		return
	}

	s.analyticsDB.AddLikes(imageID, likes, likes*50)
}

func (s *RecommendationService) DeleteAll() {
	s.analyticsDB.DeleteAll()
}

func (s *RecommendationService) RevertSagaTransaction(id string) error {
	imgageID := s.transactionList[id]
	err := s.analyticsDB.DeleteImage(imgageID)
	if err != nil {
		log.Err(err).Msg("Error reverting transaction")
		return err
	}

	return nil
}

func (s *RecommendationService) AddTransaction(id string, imageID int) {
	s.transactionList[id] = imageID
}

func (s *RecommendationService) ConfirmSagaTransaction(id string) error {
	return s.PutSagaTransaction(id, "success")
}

func (s *RecommendationService) CancelSagaTransaction(id string) error {
	return s.PutSagaTransaction(id, "failure")
}

func (s *RecommendationService) PutSagaTransaction(id string, status string) error {
	gateway := os.Getenv("GATEWAY_ADDRESS")
	url := fmt.Sprintf("http://%s/transaction/%s", gateway, id)
	payload := []byte(fmt.Sprintf(`{"status": "%s", "service": "analytics_service"}`, status))

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Err(err).Msg("Error creating request")
		return err
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		log.Err(err).Msg("Error sending request")
		return err
	}

	defer req.Body.Close()

	if res.StatusCode != 200 {
		log.Error().Err(err).Msg("Error sending SAGA CANCEL request. Status code: " + strconv.Itoa(res.StatusCode))
		return errors.New("Error sending SAGA CANCEL request. Status code: " + strconv.Itoa(res.StatusCode))
	}

	return nil
}
