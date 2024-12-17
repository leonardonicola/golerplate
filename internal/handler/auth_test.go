package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/leonardonicola/golerplate/internal/domain/entity"
	"github.com/leonardonicola/golerplate/internal/dto"
	"github.com/leonardonicola/golerplate/internal/handler"
	"github.com/leonardonicola/golerplate/pkg/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, dto dto.RegisterUserDTO) (*entity.User, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) GetByCPF(ctx context.Context, cpf string) (*entity.User, error) {
	args := m.Called(ctx, cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) Authenticate(ctx context.Context, email, password string) (*entity.User, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GenerateToken(user *entity.User) (*dto.TokenResponseDTO, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TokenResponseDTO), args.Error(1)
}

func (m *MockAuthService) RefreshToken(token string) (*dto.TokenResponseDTO, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TokenResponseDTO), args.Error(1)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    dto.LoginRequestDTO
		setupMock      func(*MockUserService, *MockAuthService)
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Successful Login",
			requestBody: dto.LoginRequestDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(mus *MockUserService, mas *MockAuthService) {
				user := &entity.User{
					ID:    "192391239",
					Email: "test@example.com",
				}

				token := &dto.TokenResponseDTO{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				}

				mus.On("Authenticate", mock.Anything, "test@example.com", "password123").Return(user, nil)
				mas.On("GenerateToken", user).Return(token, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.TokenResponseDTO{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			},
		},
		{
			name: "Invalid credentials",
			requestBody: dto.LoginRequestDTO{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(us *MockUserService, as *MockAuthService) {
				us.On("Authenticate", mock.Anything, "test@example.com", "wrongpassword").
					Return(nil, errors.New(constants.ErrMsgInvalidCredentials))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   gin.H{"message": constants.ErrMsgInvalidCredentials},
		},
		{
			name: "User not found",
			requestBody: dto.LoginRequestDTO{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(us *MockUserService, as *MockAuthService) {
				us.On("Authenticate", mock.Anything, "nonexistent@example.com", "password123").
					Return(nil, errors.New(constants.ErrMsgUserNotFound))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   gin.H{"message": constants.ErrMsgUserNotFound},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Instantiate the mocks.
			userService := new(MockUserService)
			authService := new(MockAuthService)
			handler := handler.NewAuthHandler(userService, authService)

			tt.setupMock(userService, authService)

			// Arrange.
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			// Act.
			handler.Login(c)

			// Assert.
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			expectedJson, _ := json.Marshal(tt.expectedBody)
			actualJson, _ := json.Marshal(response)
			assert.JSONEq(t, string(expectedJson), string(actualJson))

			userService.AssertExpectations(t)
			authService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Refresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    dto.RefreshRequestDTO
		setupMocks     func(*MockAuthService)
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Successful token refresh",
			requestBody: dto.RefreshRequestDTO{
				RefreshToken: "valid-refresh-token",
			},
			setupMocks: func(as *MockAuthService) {
				token := &dto.TokenResponseDTO{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token",
				}
				as.On("RefreshToken", "valid-refresh-token").Return(token, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.TokenResponseDTO{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
			},
		},
		{
			name: "Invalid refresh token",
			requestBody: dto.RefreshRequestDTO{
				RefreshToken: "invalid-refresh-token",
			},
			setupMocks: func(as *MockAuthService) {
				as.On("RefreshToken", "invalid-refresh-token").
					Return(nil, errors.New(constants.ErrMsgInvalidToken))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   gin.H{"message": constants.ErrMsgInvalidToken},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			userService := new(MockUserService)
			authService := new(MockAuthService)
			handler := handler.NewAuthHandler(userService, authService)

			tt.setupMocks(authService)

			// Create request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			handler.Refresh(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			expectedJSON, _ := json.Marshal(tt.expectedBody)
			actualJSON, _ := json.Marshal(response)
			assert.JSONEq(t, string(expectedJSON), string(actualJSON))

			authService.AssertExpectations(t)
		})
	}
}
