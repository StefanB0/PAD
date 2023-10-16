package service

import (
	"errors"
	"os"
	"strconv"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/rs/zerolog/log"
)

type TokenService struct {
	key paseto.V4SymmetricKey
	parser *paseto.Parser

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

const (
	AccessTokenType = "accessToken"
	RefreshTokenType = "refreshToken"
)

func NewTokenService() *TokenService {
	accessTokenDurationString := os.Getenv("ACCESS_TOKEN_DURATION")
	accessTokenDuration, err := strconv.Atoi(accessTokenDurationString)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing ACCESS_TOKEN_DURATION")
	}

	refreshTokenDurationString := os.Getenv("REFRESH_TOKEN_DURATION")
	refreshTokenDuration, err := strconv.Atoi(refreshTokenDurationString)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing REFRESH_TOKEN_DURATION")
	}

	key, err := paseto.V4SymmetricKeyFromHex(os.Getenv("PASETO_KEY"))
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing PASETO_KEY")
	}

	parser := paseto.NewParser()

	return &TokenService{
		key:                  key,
		parser:               &parser,
		accessTokenDuration:  time.Duration(accessTokenDuration) * time.Second,
		refreshTokenDuration: time.Duration(refreshTokenDuration) * time.Second,
	}
}

func (ts *TokenService) NewAccessToken(userID int, userName string) string {
	return ts.NewUserToken(userID, userName, AccessTokenType, ts.accessTokenDuration)
}

func (ts *TokenService) NewRefreshToken(userID int, userName string) string {
	return ts.NewUserToken(userID, userName, RefreshTokenType, ts.refreshTokenDuration)
}

func (ts *TokenService) NewUserToken(userID int, userName string, tokenType string, duration time.Duration) string {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(duration))

	token.SetString("user-id", strconv.Itoa(userID))
	token.SetString("token-type", tokenType)
	token.SetString("user-name", userName)
	encryptedToken := token.V4Encrypt(ts.key, nil)

	return encryptedToken
}

func (ts *TokenService) VerifyAccessToken(token string) (*paseto.Token, error) {
	return ts.VerifyToken(token, AccessTokenType)
}

func (ts *TokenService) VerifyRefreshToken(token string) (*paseto.Token, error) {
	return ts.VerifyToken(token, RefreshTokenType)
}

func (ts *TokenService) VerifyToken(token string, tokenType string) (*paseto.Token, error) {
	pasetoToken, err := ts.parser.ParseV4Local(ts.key, token, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing token")
		return nil, err
	}

	if pasetoToken.Claims()["token-type"] != tokenType {
		log.Error().Msg("Token type mismatch")
		return nil, errors.New("Token type mismatch")
	}

	return pasetoToken, nil
}

func (ts *TokenService) ExportKeyHex() string {
	return ts.key.ExportHex()
}

func (ts *TokenService) ImportKeyHex(hex string) error {
	key, err := paseto.V4SymmetricKeyFromHex(hex)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing key")
		return err
	}

	ts.key = key

	return nil
}
