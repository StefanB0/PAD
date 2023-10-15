package main

import (
	"os"
	"padauth/controller"
	"padauth/database"
	"padauth/service"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	loadConfig()
	startLogger()

	postgresDB := database.NewPostgresDatabase(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	tokenService := service.NewTokenService()
	authService := service.NewAuthService(postgresDB, tokenService)

	controller := controller.NewUserController(authService)
	controller.Run()
}

func startLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
}
