package handler

import (
	"avito_test/internal/domain"
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Auth(ctx context.Context, req domain.AuthRequest) (*domain.AuthResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.AuthResponse), args.Error(1)
}

func TestAuthHandler(t *testing.T) {
	mockService := new(MockAuthService)

	handler := NewAuth(mockService)

	app := fiber.New()
	app.Post("/auth", handler.Auth())

	tests := []struct {
		name           string
		body           string
		mock           func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			body: `{"username": "validUser", "password": "ValidPass123!"}`,
			mock: func() {
				mockService.On("Auth", mock.Anything, domain.AuthRequest{
					Username: "validUser",
					Password: "ValidPass123!",
				}).Return(&domain.AuthResponse{Token: "12345"}, nil)
			},
			expectedStatus: fiber.StatusOK,
			expectedBody:   `{"token":"12345"}`,
		},
		{
			name: "Invalid Credentials",
			body: `{"username": "invalidUser", "password": "invalidPass"}`,
			mock: func() {
				mockService.On("Auth", mock.Anything, domain.AuthRequest{
					Username: "invalidUser",
					Password: "invalidPass",
				}).Return((*domain.AuthResponse)(nil), domain.ErrInvalidCredentials)
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody:   `{"errors":"invalid credentials"}`,
		},
		{
			name: "Unauthorized",
			body: `{"username": "unauthUser", "password": "unauthPass"}`,
			mock: func() {
				mockService.On("Auth", mock.Anything, domain.AuthRequest{
					Username: "unauthUser",
					Password: "unauthPass",
				}).Return((*domain.AuthResponse)(nil), domain.ErrUnauthorized)
			},
			expectedStatus: fiber.StatusUnauthorized,
			expectedBody:   `{"errors":"unauthorized"}`,
		},
		{
			name: "Internal Server Error",
			body: `{"username": "errorUser", "password": "errorPass"}`,
			mock: func() {
				mockService.On("Auth", mock.Anything, domain.AuthRequest{
					Username: "errorUser",
					Password: "errorPass",
				}).Return((*domain.AuthResponse)(nil), errors.New("internal error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody:   `{"errors":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			req := httptest.NewRequest("POST", "/auth", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expectedBody, string(body))

			mockService.AssertExpectations(t)
		})
	}
}
