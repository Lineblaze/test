package repository_test

import (
	"avito_test/internal/config"
	"avito_test/internal/entity"
	"avito_test/internal/repository"
	"avito_test/pkg/storage/postgres"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) postgres.Postgres {
	config := &config.Config{
		Postgres: struct {
			Host     string `json:"host"`
			Port     string `json:"port"`
			User     string `json:"user"`
			Password string `json:"password"`
			DBName   string `json:"dbName"`
			SSLMode  string `json:"sslMode"`
			PgDriver string `json:"pgDriver"`
		}{
			Host:     "localhost",
			Port:     "5432",
			User:     "test_user",
			Password: "test_password",
			DBName:   "test_db",
			SSLMode:  "disable",
		},
	}

	db, err := postgres.InitPsqlDB(config)
	require.NoError(t, err, "Failed to connect to test database")

	_, err = db.Exec(context.Background(), `DELETE FROM auth`)
	require.NoError(t, err)

	_, err = db.Exec(context.Background(), `DELETE FROM users`)
	require.NoError(t, err)

	return db
}

func TestAuth(t *testing.T) {
	db := setupTestDB(t)
	authRepo := repository.NewAuth(db)

	ctx := context.Background()
	testUser := entity.Auth{
		Username: "testuser",
		Password: "password123",
	}

	createdAuth, err := authRepo.Auth(ctx, testUser)
	require.NoError(t, err, "Auth method should not return an error")
	require.NotNil(t, createdAuth, "Auth method should return a valid entity")

	var username string
	err = db.QueryRow(ctx, `SELECT username FROM auth WHERE id = $1`, createdAuth.Id).Scan(&username)
	require.NoError(t, err, "Failed to retrieve user from DB")
	assert.Equal(t, testUser.Username, username, "Stored username should match input username")

	var createdAt time.Time
	err = db.QueryRow(ctx, `SELECT created_at FROM users WHERE id = $1`, createdAuth.Id).Scan(&createdAt)
	require.NoError(t, err, "Failed to retrieve user from users table")
	assert.WithinDuration(t, time.Now(), createdAt, 2*time.Second, "Created timestamp should be recent")
}
