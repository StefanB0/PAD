package service

import (
	"testing"

	"aidanwoods.dev/go-paseto"
	"github.com/joho/godotenv"
)

func TestNewTokenService(t *testing.T) {
	godotenv.Load("../.env")

	tokenService := NewTokenService()

	expectedKey := "173d492b114a7816f70cada3e0e52648723edc1894247143781a0cc51553f802"
	if tokenService.key.ExportHex() != expectedKey {
		t.Errorf("Expected key to be %s, got %s", expectedKey, tokenService.key.ExportHex())
	}

	expectedAccessTokenDuration := 900
	if int(tokenService.accessTokenDuration.Seconds()) != expectedAccessTokenDuration {
		t.Errorf("Expected accessTokenDuration to be %d, got %d", expectedAccessTokenDuration, int(tokenService.accessTokenDuration.Seconds()))
	}

	expectedRefreshTokenDuration := 7200
	if int(tokenService.refreshTokenDuration.Seconds()) != expectedRefreshTokenDuration {
		t.Errorf("Expected refreshTokenDuration to be %d, got %d", expectedRefreshTokenDuration, int(tokenService.refreshTokenDuration.Seconds()))
	}
}

func TestNewTokens(t *testing.T) {
	godotenv.Load("../.env")

	tokenService := NewTokenService()

	accessToken := tokenService.NewAccessToken(1)

	if accessToken == "" {
		t.Errorf("Expected accessToken to not be empty")
	}

	refreshToken := tokenService.NewRefreshToken(1)

	if refreshToken == "" {
		t.Errorf("Expected refreshToken to not be empty")
	}
}

func TestVerifyAccessToken(t *testing.T) {
	godotenv.Load("../.env")

	tokenService := NewTokenService()

	accessToken := tokenService.NewAccessToken(1)

	pasetoAccessToken, err := tokenService.VerifyAccessToken(accessToken)
	if err != nil {
		t.Errorf("Error parsing access token %s", err.Error())
	}

	var userID string
	err = pasetoAccessToken.Get("user-id", &userID)
	if err != nil {
		t.Errorf("Error parsing user-id %s", err.Error())
	}

	if userID != "1" {
		t.Errorf("Expected userID to be 1, got %s", userID)
	}	
}

func TestVerifyRefreshToken(t *testing.T) {
	godotenv.Load("../.env")

	tokenService := NewTokenService()

	refreshToken := tokenService.NewRefreshToken(1)

	pasetoRefreshToken, err := tokenService.VerifyRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("Error parsing refresh token %s", err.Error())
	}

	var userID string
	err = pasetoRefreshToken.Get("user-id", &userID)
	if err != nil {
		t.Errorf("Error parsing user-id %s", err.Error())
	}

	if userID != "1" {
		t.Errorf("Expected userID to be 1, got %s", userID)
	}	
}

func TestExportHex(t *testing.T) {
	godotenv.Load("../.env")

	tokenService := NewTokenService()

	key := tokenService.key.ExportHex()

	expectedKey := "173d492b114a7816f70cada3e0e52648723edc1894247143781a0cc51553f802"
	if key != expectedKey {
		t.Errorf("Expected key to be %s, got %s", expectedKey, key)
	}
}

func TestImportHex(t *testing.T) {
	godotenv.Load("../.env")

	tokenService := NewTokenService()

	newkey := paseto.NewV4SymmetricKey().ExportHex()

	tokenService.ImportKeyHex(newkey)

	key := tokenService.key.ExportHex()

	if key != newkey {
		t.Errorf("Expected key to be %s, got %s", newkey, key)
	}
}
