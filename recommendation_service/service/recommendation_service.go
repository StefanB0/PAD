package service

import (
	"fmt"
	"math/rand"
	"padrecommendations/database"
	"padrecommendations/models"
)

type RecommendationService struct {
	analyticsDB *database.AnalyticsPostgresDB
	// analyticsDB *database.AnalyticsPostgresDB
}

func NewRecommendationService(analyticsDB *database.AnalyticsPostgresDB) *RecommendationService {
	return &RecommendationService{
		analyticsDB: analyticsDB,
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

func (s *RecommendationService) AddImage(image models.Image) {
	image.Engagement = image.Views + image.Likes
	image.Engagement += 100
	s.analyticsDB.AddImage(image)
}

// func (service *RecommendationService) UpdateImage(image models.Image) {
// 	service.analyticsDB.UpdateImage(image)
// }

// func (s *RecommendationService) GetImage(imageID int) (*models.Image, error) {
// 	return s.analyticsDB.GetImage(imageID)
// }

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
