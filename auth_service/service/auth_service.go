package service

import (
	"errors"
	"padauth/database"
	"padauth/models"
	"strconv"

	"github.com/rs/zerolog/log"
)

type AuthService struct {
	ts *TokenService
	db *database.UserPostgresDB
}

func NewAuthService(db *database.UserPostgresDB, tservice *TokenService) *AuthService {
	return &AuthService{
		db: db,
		ts: tservice,
	}
}

func (s *AuthService) Login(username, password string) (accessToken string, refreshToken string, err error) {
	user, err := s.db.GetUser(username)
	if err != nil {
		return "", "", err
	}

	if !comparePasswords(user.Password, password) {
		return "", "", errors.New("passwords do not match")
	}

	accessToken = s.ts.NewAccessToken(user.ID)
	refreshToken = s.ts.NewRefreshToken(user.ID)

	err = s.cacheTokens(accessToken, refreshToken)
	if err != nil {
		log.Error().Err(err).Msg("Error caching tokens")
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Register(username, password string) error {
	user := models.User{
		ID:       0,
		Username: username,
		Password: hashPassword(password),
	}

	err := s.db.CreateUser(user)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user")
		return err
	}

	return nil
}

func (s *AuthService) Delete(username string) error {
	s.db.DeleteUser(username)
	return nil
}

func (s *AuthService) Refresh(pasetoToken string) (accessToken string, refreshToken string, err error) {
	token, err := s.ts.VerifyRefreshToken(pasetoToken)
	if err != nil {
		log.Error().Err(err).Msg("Error verifying refresh token")
		return "", "", err
	}

	var userStringID string
	err = token.Get("user-id", &userStringID)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing user-id")
		return "", "", err
	}

	userID, err := strconv.Atoi(userStringID)
	if err != nil {
		log.Error().Err(err).Msg("Error converting user-id to int")
		return "", "", err
	}

	accessToken = s.ts.NewAccessToken(userID)
	refreshToken = s.ts.NewRefreshToken(userID)

	err = s.cacheTokens(accessToken, refreshToken)
	if err != nil {
		log.Error().Err(err).Msg("Error caching tokens")
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) cacheTokens(accessToken string, refreshToken string) error {
	return nil
}
