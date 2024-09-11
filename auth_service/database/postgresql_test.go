package database

import (
	"os"
	"padauth/models"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

func postgresSetup(t *testing.T) *UserPostgresDB {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	godotenv.Load("../.env")

	db := NewPostgresDatabase(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db.userDB.Logger.LogMode(logger.Silent)

	if db == nil {
		t.Error("PostgresDatabase was not created")
	}
	return db
}

func TestCreateUser(t *testing.T) {
	db := postgresSetup(t)
	defer db.DeleteAll()

	if db == nil {
		return
	}

	user := models.User{
		Username: "testuser19",
		Password: "testpassword",
	}

	err := db.CreateUser(user)

	if err != nil {
		t.Errorf("Error creating user: %s", err)
		return
	}
}

func TestGetUser(t *testing.T) {
	db := postgresSetup(t)
	defer db.DeleteAll()

	if db == nil {
		return
	}

	user := &models.User{
		Username: "testuser11",
		Password: "testpassword",
	}

	err := db.CreateUser(*user)

	if err != nil {
		t.Errorf("Error creating user: %s", err)
		return
	}

	user, err = db.GetUser("testuser11")

	if err != nil {
		t.Errorf("Error retrieving user: %s", err)
		return
	}

	if user.Username != "testuser11" {
		t.Errorf("User was not retrieved correctly")
		return
	}
}

func TestDeleteUser(t *testing.T) {
	db := postgresSetup(t)
	defer db.DeleteAll()

	if db == nil {
		return
	}

	user := &models.User{
		Username: "testuser32",
		Password: "testpassword",
	}

	err := db.CreateUser(*user)

	if err != nil {
		t.Errorf("Error creating user: %s", err)
		return
	}

	err = db.DeleteUser("testuser32")

	if err != nil {
		t.Errorf("Error deleting user: %s", err)
		return
	}

	user, err = db.GetUser("testuser32")

	if err == nil {
		t.Errorf("User was not deleted")
		return
	}
}

func TestGetAll(t *testing.T) {
	db := postgresSetup(t)

	if db == nil {
		return
	}

	user := models.User{
		Username: "testuser44",
		Password: "testpassword",
	}

	err := db.CreateUser(user)

	if err != nil {
		t.Errorf("Error creating user: %s", err)
		return
	}

	users, err := db.GetAll()

	if err != nil {
		t.Errorf("Error retrieving users: %s", err)
		return
	}

	if len(users) == 0 {
		t.Errorf("Users were not retrieved correctly")
		return
	}
}
