package main

import (
	"padrecommendations/controller"
	"padrecommendations/database"
	"padrecommendations/lib"
	"padrecommendations/service"
)

func main() {
	flags := lib.GetFlags()

	analyticsDB := database.NewPostgresDatabase()
	recommendationService := service.NewRecommendationService(analyticsDB)
	controller := controller.NewController(recommendationService)

	discoveryService := service.NewDiscoveryService("http://"+flags["host"]+":8083", flags["service_discovery"])
	discoveryService.Subscribe()

	controller.Run(":8083")
}
