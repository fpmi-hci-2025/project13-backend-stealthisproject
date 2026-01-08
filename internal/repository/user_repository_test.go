package repository

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/project13/backend-stealthisproject/internal/models"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Use test database URL or skip if not available
	dbURL := "postgres://postgres:postgres@localhost:5432/railway_tickets_test?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
		return nil
	}

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewUserRepository(db)

	user := &models.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "PASSENGER",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("User ID should be set after creation")
	}
}

