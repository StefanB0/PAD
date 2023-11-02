package main

import (
	"fmt"
	"padimage/controller"
	"padimage/database"
	"padimage/lib"
	"padimage/service"

	"github.com/rs/zerolog/log"
)

func main() {
	flags := lib.GetFlags()

	log.Info().Msg(fmt.Sprintf("Service discovery: %s", flags["service_discovery"]))
	log.Info().Msg(fmt.Sprintf("Self address: %s", flags["host"]))

	database := database.NewMongoDB()
	tokenService := service.NewTokenService(true)
	discoveryService := service.NewDiscoveryService("http://" + flags["host"] + ":8082", flags["service_discovery"])
	discoveryService.Subscribe()
	analyticsService := service.NewAnalyticsService(discoveryService)
	imageService := service.NewImageService(database, analyticsService, tokenService)
	controller := controller.NewImageController(imageService, tokenService)

	controller.Run(":8082")
}
