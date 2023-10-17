package service

import "padrecommendations/database"

type RecommendationService struct {
	analyticsDB *database.AnalyticsPostgresDB
}

func NewRecommendationService(analyticsDB *database.AnalyticsPostgresDB) *RecommendationService {
	return &RecommendationService{
		analyticsDB: analyticsDB,
	}
}


