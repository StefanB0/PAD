package service

import (
	"os"
	"padauth/database"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func setup(t *testing.T) *AuthService {
	godotenv.Load("../.env")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")

	db := database.NewPostgresDatabase(postgresUser, postgresPassword, postgresDB)
	tokenService := NewTokenService()
	authService := NewAuthService(db, tokenService)

	if authService == nil {
		t.Errorf("Expected authService to not be nil")
	}

	return authService
}

func TestRegister(t *testing.T) {
	authService := setup(t)
	defer authService.db.DeleteAll()

	if authService == nil {
		return
	}

	err := authService.Register("testuser", "testpassword")
	if err != nil {
		t.Errorf("Error registering user: %s", err)
		return
	}
}

func TestLogin(t *testing.T) {
	authService := setup(t)
	defer authService.db.DeleteAll()

	if authService == nil {
		return
	}

	err := authService.Register("testuser", "testpassword")
	if err != nil {
		t.Errorf("Error registering user: %s", err)
		return
	}

	_, _, err = authService.Login("testuser", "testpassword")
	if err != nil {
		t.Errorf("Error logging in: %s", err)
		return
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	authService := setup(t)
	defer authService.db.DeleteAll()

	if authService == nil {
		return
	}

	err := authService.Register("testuser", "testpassword")
	if err != nil {
		t.Errorf("Error registering user: %s", err)
		return
	}

	_, _, err = authService.Login("testuser", "wrongpassword")
	if err == nil {
		t.Errorf("Expected error logging in with invalid password")
		return
	}
}

func TestLoginInvalidUsername(t *testing.T) {
	authService := setup(t)
	defer authService.db.DeleteAll()

	if authService == nil {
		return
	}

	err := authService.Register("testuser", "testpassword")
	if err != nil {
		t.Errorf("Error registering user: %s", err)
		return
	}

	_, _, err = authService.Login("wronguser", "testpassword")
	if err == nil {
		t.Errorf("Expected error logging in with invalid username")
		return
	}
}
