package main

import (
	"padimage/controller"
	"padimage/database"
	"padimage/service"
)

func main() {
	database := database.NewMongoDB()
	tokenService := service.NewTokenService(true)
	discoveryService := service.NewDiscoveryService("http://localhost:8081", "http://localhost:8500")
	analyticsService := service.NewAnalyticsService(discoveryService)
	imageService := service.NewImageService(database, analyticsService, tokenService)
	controller := controller.NewImageController(imageService, tokenService)

	discovery := service.NewDiscoveryService("http://localhost:8081", "http://localhost:8500")
	discovery.Subscribe()

	controller.Run()
}
