package service

import (
	"math/rand"
	"padrecommendations/database"
	"padrecommendations/models"
)

type RecommendationService struct {
	analyticsDB *database.AnalyticsPostgresDB
}

func NewRecommendationService(analyticsDB *database.AnalyticsPostgresDB) *RecommendationService {
	return &RecommendationService{
		analyticsDB: analyticsDB,
	}
}

func (service *RecommendationService) GetTags() []string {
	return service.analyticsDB.GetTags()
}

func (service *RecommendationService) GetRecommendations(tagname string) (int, error) {
	taglist, err := service.analyticsDB.GetTaglist(tagname)
	if err != nil {
		return 0, err
	}

	randNum := rand.Intn(len(taglist.ImageList)) + 1
	randImage := 1

	var counter int

	for _, image := range taglist.ImageList {
		counter += image.Engagement
		if counter >= randNum {
			randImage = image.ImageID
			break
		}
	}

	return randImage, nil
}

func (service *RecommendationService) AddImage(image models.Image) {
	service.analyticsDB.AddImage(image)
}

func (service *RecommendationService) UpdateImage(image models.Image) {
	service.analyticsDB.UpdateImage(image)
}

func (service *RecommendationService) GetImage(imageID int) (models.Image, error) {
	return service.analyticsDB.GetImage(imageID)
}

func (service *RecommendationService) AddView(imageID int, views int) {
	_, err := service.analyticsDB.GetImage(imageID)
	if err != nil {
		return
	}

	service.analyticsDB.AddViews(imageID, views)
}

func (service *RecommendationService) AddLike(imageID int, likes int) {
	_, err := service.analyticsDB.GetImage(imageID)
	if err != nil {
		return
	}

	service.analyticsDB.AddLikes(imageID, likes)
}
