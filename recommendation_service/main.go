package main

import (
	"padrecommendations/controller"
	"padrecommendations/database"
	"padrecommendations/service"
)

func main() {
	analyticsDB := database.NewPostgresDatabase()
	recommendationService := service.NewRecommendationService(analyticsDB)
	controller := controller.NewController(recommendationService)

	discoveryService := service.NewDiscoveryService("http://localhost:8083", "http://localhost:8500")
	discoveryService.Subscribe()

	controller.Run()
}
