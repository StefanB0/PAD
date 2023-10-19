package main

import (
	"padimage/controller"
	"padimage/database"
	"padimage/lib"
	"padimage/service"
)

func main() {
	port, sd := lib.GetFlags()

	database := database.NewMongoDB()
	tokenService := service.NewTokenService(true)
	discoveryService := service.NewDiscoveryService("http://localhost"+port, sd)
	discoveryService.Subscribe()
	analyticsService := service.NewAnalyticsService(discoveryService)
	imageService := service.NewImageService(database, analyticsService, tokenService)
	controller := controller.NewImageController(imageService, tokenService)

	controller.Run(port)
}
