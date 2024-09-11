package service

import (
	"crypto/rand"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)



func hashPassword(password string) string {
	passwordBytes := []byte(password)

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
	}

	return string(hashedPasswordBytes[:])
}

func comparePasswords(hashedPassword string, password string) bool {
	hashedPasswordBytes := []byte(hashedPassword)
	passwordBytes := []byte(password)

	err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes)

	return err == nil
}

func generateSalt(saltSize int) []byte {
	salt := make([]byte, saltSize)

	_, err := rand.Read(salt)
	if err != nil {
		log.Error().Err(err).Msg("Error generating salt")
	}

	return salt
}