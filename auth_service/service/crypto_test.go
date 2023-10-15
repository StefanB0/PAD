package service

import "testing"

func TestHashPassword(t *testing.T) {
	passwords := []string{
		"password",
		"password123",
		"password123456",
		"kentuckyfriedchicken",
		"americanexpress",
		"starbucks",
	}

	hashedPasswords := make([]string, len(passwords))

	for i, password := range passwords {
		hashedPasswords[i] = hashPassword(password)
	}

	for i, hash := range hashedPasswords {
		if !comparePasswords(hash, passwords[i]) {
			t.Errorf("Password %s was not hashed correctly", passwords[i])
		}
	}
}
