package handler

import (
	"avito_test/internal/config"
	"avito_test/internal/domain"
	"avito_test/internal/jwt"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) Buy(ctx context.Context, userIDStr string, itemType string) error {
	args := m.Called(ctx, userIDStr, itemType)
	return args.Error(0)
}

func (m *MockTransactionService) Send(ctx context.Context, userIDStr string, req domain.SendCoinRequest) error {
	args := m.Called(ctx, userIDStr, req)
	return args.Error(0)
}

func (m *MockTransactionService) Info(ctx context.Context, userIDStr string) (*domain.InfoResponse, error) {
	args := m.Called(ctx, userIDStr)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.InfoResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestTransactionHandler_Buy(t *testing.T) {
	mockService := new(MockTransactionService)

	handler := NewTransaction(mockService)
	app := fiber.New()

	cfg := &config.Config{
		Auth: struct {
			Secret string `json:"secret"`
		}{
			Secret: "secret",
		},
	}

	jwtService := jwt.NewJWTService(cfg)

	validUserID := uuid.New().String()
	invalidUserID := "invalid-uuid"

	token, err := jwtService.GenerateJWT(jwt.Claims{
		ID:       validUserID,
		Username: "testuser",
	})

	require.NoError(t, err)

	app.Post("/buy/:itemType", handler.Buy())

	tests := []struct {
		name           string
		userID         string
		itemType       string
		mock           func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "Success",
			userID:   validUserID,
			itemType: "t-shirt",
			mock: func() {
				mockService.On("Buy", mock.Anything, validUserID, "t-shirt").Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			expectedBody:   `{}`,
		},
		{
			name:           "Invalid UUID",
			userID:         invalidUserID,
			itemType:       "cup",
			mock:           func() {},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody:   `{"errors":"invalid user ID format"}`,
		},
		{
			name:           "Invalid Item Type",
			userID:         validUserID,
			itemType:       "!!!",
			mock:           func() {},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody:   `{"errors":"invalid item type format"}`,
		},
		{
			name:     "Internal Server Error",
			userID:   validUserID,
			itemType: "hoody",
			mock: func() {
				mockService.On("Buy", mock.Anything, validUserID, "hoody").Return(errors.New("db error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody:   `{"errors":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			req := httptest.NewRequest(http.MethodPost, "/buy/"+tt.itemType, bytes.NewBuffer(nil))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			resp, _ := app.Test(req)

			require.Equal(t, tt.expectedStatus, resp.StatusCode)

			var body map[string]string
			_ = json.NewDecoder(resp.Body).Decode(&body)
			require.Equal(t, tt.expectedBody, body["errors"])
		})
	}
}

func TestTransactionHandler_Send(t *testing.T) {
	mockService := new(MockTransactionService)

	handler := NewTransaction(mockService)
	app := fiber.New()

	cfg := &config.Config{
		Auth: struct {
			Secret string `json:"secret"`
		}{
			Secret: "secret",
		},
	}

	jwtService := jwt.NewJWTService(cfg)

	validUserID := uuid.New().String()
	token, err := jwtService.GenerateJWT(jwt.Claims{
		ID:       validUserID,
		Username: "testuser",
	})
	require.NoError(t, err)

	app.Post("/send", handler.Send())

	tests := []struct {
		name           string
		userID         string
		requestBody    domain.SendCoinRequest
		mock           func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Success",
			userID: validUserID,
			requestBody: domain.SendCoinRequest{
				ToUser: "test_user",
				Amount: 10,
			},
			mock: func() {
				mockService.On("Send", mock.Anything, validUserID, mock.Anything).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			expectedBody:   `{}`,
		},
		{
			name:   "Invalid request body",
			userID: validUserID,
			requestBody: domain.SendCoinRequest{
				ToUser: "test_user",
				Amount: -5,
			},
			mock:           func() {},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody:   `{"errors":"invalid request body"}`,
		},
		{
			name:   "Internal Server Error",
			userID: validUserID,
			requestBody: domain.SendCoinRequest{
				ToUser: "test_user",
				Amount: 15,
			},
			mock: func() {
				mockService.On("Send", mock.Anything, validUserID, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody:   `{"errors":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			resp, _ := app.Test(req)

			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTransactionHandler_Info(t *testing.T) {
	mockService := new(MockTransactionService)

	handler := NewTransaction(mockService)
	app := fiber.New()

	cfg := &config.Config{
		Auth: struct {
			Secret string `json:"secret"`
		}{
			Secret: "secret",
		},
	}

	jwtService := jwt.NewJWTService(cfg)

	validUserID := uuid.New().String()
	token, err := jwtService.GenerateJWT(jwt.Claims{
		ID:       validUserID,
		Username: "testuser",
	})
	require.NoError(t, err)

	app.Get("/info", handler.Info())

	tests := []struct {
		name           string
		userID         string
		mock           func()
		expectedStatus int
	}{
		{
			name:   "Success",
			userID: validUserID,
			mock: func() {
				mockService.On("Info", mock.Anything, validUserID).Return(&domain.InfoResponse{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:   "Internal Server Error",
			userID: validUserID,
			mock: func() {
				mockService.On("Info", mock.Anything, validUserID).Return(nil, errors.New("db error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			req := httptest.NewRequest(http.MethodGet, "/info", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, _ := app.Test(req)

			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
