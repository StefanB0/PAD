package main

import (
	"padrecommendations/controller"
	"padrecommendations/database"
	"padrecommendations/service"
)

func main() {
	analyticsDB := database.NewAnalyticsPostgresDB()
	recommendationService := service.NewRecommendationService(analyticsDB)
	controller := controller.NewController(recommendationService)

	controller.Run()
}
