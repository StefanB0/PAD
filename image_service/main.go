package main

import (
	"padimage/controller"
	"padimage/database"
	"padimage/lib"
	"padimage/service"
)

func main() {
	flags := lib.GetFlags()

	database := database.NewMongoDB()
	tokenService := service.NewTokenService(true)
	discoveryService := service.NewDiscoveryService("http://"+flags["host"]+":8082", flags["service_discovery"])
	discoveryService.Subscribe()
	analyticsService := service.NewAnalyticsService(discoveryService)
	imageService := service.NewImageService(database, analyticsService, tokenService)
	controller := controller.NewImageController(imageService, tokenService)

	controller.Run(":8082")
}
