package service_test

import (
	"context"
	"testing"

	"github.com/leonardonicola/golerplate/internal/domain/entity"
	"github.com/leonardonicola/golerplate/internal/domain/service"
	"github.com/leonardonicola/golerplate/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByCPF(ctx context.Context, cpf string) (*entity.User, error) {
	args := m.Called(ctx, cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	testCases := []struct {
		name        string
		user        dto.RegisterUserDTO
		setupMock   func()
		want        *entity.User
		expectError bool
	}{
		{
			name: "Success",
			user: dto.RegisterUserDTO{
				Email:    "test@gmail.com",
				CPF:      "15245901854",
				Age:      20,
				FullName: "Test",
				Password: "aosdaosdoa",
			},
			setupMock: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(&entity.User{
					ID:       "generated-id",
					Email:    "test@gmail.com",
					CPF:      "15245901854",
					Age:      20,
					FullName: "Test",
					Password: "@12msa38sjdasdji123",
				}, nil)
			},
			want: &entity.User{
				ID:       "generated-id",
				Email:    "test@gmail.com",
				CPF:      "15245901854",
				Age:      20,
				FullName: "Test",
				Password: "@12msa38sjdasdji123",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			got, err := userService.Create(context.Background(), tc.user)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NoError(t, err)
				assert.Equal(t, tc.want.ID, got.ID)
				assert.Equal(t, tc.want.FullName, got.FullName)
				assert.Equal(t, tc.want.Email, got.Email)
				assert.Equal(t, tc.want.CPF, got.CPF)
				assert.Equal(t, tc.want.Age, got.Age)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
