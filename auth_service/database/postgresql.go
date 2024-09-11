package database

import (
	"padauth/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserPostgresDB struct {
	userDB *gorm.DB
}

// Create new user repository
func NewPostgresDatabase(user string, password string, dbName string) *UserPostgresDB {
	dbURL := "postgres://" + user + ":" + password + "@localhost:5432/" + dbName

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Error().Err(err).Msg("Error connecting to Postgres")
		return nil
	}

	db.AutoMigrate(&models.User{})

	log.Info().Msg("Connected to Postgres")
	return &UserPostgresDB{
		userDB: db,
	}
}

func (db *UserPostgresDB) GetUser(username string) (*models.User, error) {
	var user models.User

	result := db.userDB.First(&user, "username = ?", username)
	// db.userDB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error retrieving user from Postgres")
		return nil, result.Error
	}

	return &user, nil
}

func (db *UserPostgresDB) CreateUser(user models.User) error {
	result := db.userDB.Create(&user)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error creating user in Postgres")
		return result.Error
	}

	return nil
}

func (db *UserPostgresDB) DeleteUser(username string) error {
	result := db.userDB.Delete(&models.User{}, "username = ?", username)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error deleting user from Postgres")
		return result.Error
	}

	return nil
}

func (db *UserPostgresDB) GetAll() ([]models.User, error) {
	var users []models.User

	result := db.userDB.Find(&users)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error retrieving users from Postgres")
		return nil, result.Error
	}

	return users, nil
}

func (db *UserPostgresDB) DeleteAll() error {
	result := db.userDB.Where("1 = 1").Delete(&models.User{})

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error deleting users from Postgres")
		return result.Error
	}

	return nil
}
