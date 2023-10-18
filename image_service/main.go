package main

import (
	"padimage/controller"
	"padimage/database"
	"padimage/service"
)

func main() {
	database := database.NewMongoDB()
	tokenService := service.NewTokenService()
	imageService := service.NewImageService(database, tokenService)
	controller := controller.NewImageController(imageService, tokenService)

	controller.Run()
}
