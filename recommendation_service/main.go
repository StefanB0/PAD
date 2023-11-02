package main

import (
	"fmt"
	"padrecommendations/controller"
	"padrecommendations/database"
	"padrecommendations/lib"
	"padrecommendations/service"

	"github.com/rs/zerolog/log"
)

func main() {
	flags := lib.GetFlags()

	log.Info().Msg(fmt.Sprintf("Service discovery: %s", flags["service_discovery"]))
	log.Info().Msg(fmt.Sprintf("Self address: %s", flags["host"]))

	analyticsDB := database.NewPostgresDatabase()
	recommendationService := service.NewRecommendationService(analyticsDB)
	controller := controller.NewController(recommendationService)

	discoveryService := service.NewDiscoveryService("http://" + flags["host"] + ":8083", flags["service_discovery"])
	discoveryService.Subscribe()

	controller.Run(":8083")
}
