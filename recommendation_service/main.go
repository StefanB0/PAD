package main

import (
	"padrecommendations/controller"
	"padrecommendations/database"
	"padrecommendations/lib"
	"padrecommendations/service"
)

func main() {
	port, sd := lib.GetFlags()

	analyticsDB := database.NewPostgresDatabase()
	recommendationService := service.NewRecommendationService(analyticsDB)
	controller := controller.NewController(recommendationService)

	discoveryService := service.NewDiscoveryService("http://localhost"+port, sd)
	discoveryService.Subscribe()

	controller.Run(port)
}
