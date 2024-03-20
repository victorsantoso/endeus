package usecase

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/victorsantoso/endeus/domain"
	"github.com/victorsantoso/endeus/entity"
	mocks "github.com/victorsantoso/endeus/mocks/domain"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase_Register(t *testing.T) {

	t.Run("test register success", func(t *testing.T) {
		mockUserRepository := new(mocks.UserRepository)
		userUsecase := NewUserUsecase(mockUserRepository)
		registerDTO := &domain.RegisterDTO{
			Role:         domain.ADMIN,
			Email:        "testtest@gmail.com",
			Password:     "Test*999",
			Name:         "Test User",
			ProfileImage: "Google Cloud Storage Sample Path",
		}
		mockUserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
		mockUserRepository.On("Create", mock.Anything, mock.Anything).Return(registerDTO.Role, int64(1), nil)
		accessToken, err := userUsecase.Register(context.Background(), registerDTO)
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		mockUserRepository.AssertCalled(t, "FindByEmail", mock.Anything, mock.Anything)
		mockUserRepository.AssertCalled(t, "Create", mock.Anything, mock.Anything)
		defer mockUserRepository.AssertExpectations(t)
	})

	t.Run("test register failed", func(t *testing.T) {
		mockUserRepository := new(mocks.UserRepository)
		userUsecase := NewUserUsecase(mockUserRepository)
		registerDTO := &domain.RegisterDTO{
			Role:         "random", // register with invalid role
			Email:        "testtest@gmail.com",
			Password:     "Test*999",
			Name:         "Test User",
			ProfileImage: "Google Cloud Storage Sample Path",
		}
		mockUserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
		mockUserRepository.On("Create", mock.Anything, mock.Anything).Return("", int64(0), domain.ErrInvalidRole)
		accessToken, err := userUsecase.Register(context.Background(), registerDTO)
		assert.Error(t, err)
		assert.Empty(t, accessToken)
		mockUserRepository.AssertCalled(t, "FindByEmail", mock.Anything, mock.Anything)
		mockUserRepository.AssertCalled(t, "Create", mock.Anything, mock.Anything)
		defer mockUserRepository.AssertExpectations(t)
	})

	t.Run("test register existing user", func(t *testing.T) {
		mockUserRepository := new(mocks.UserRepository)
		userUsecase := NewUserUsecase(mockUserRepository)
		registerDTO := &domain.RegisterDTO{
			Role:         domain.ADMIN,
			Email:        "testtest@gmail.com",
			Password:     "Test*999",
			Name:         "Test User",
			ProfileImage: "Google Cloud Storage Sample Path",
		}
		mockUserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&entity.User{
			UserId:       1,
			Role:         registerDTO.Role,
			Email:        registerDTO.Email,
			Password:     registerDTO.Password,
			Name:         registerDTO.Name,
			ProfileImage: registerDTO.ProfileImage,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		}, nil)
		accessToken, err := userUsecase.Register(context.Background(), registerDTO)
		assert.Error(t, err)
		assert.Empty(t, accessToken)
		mockUserRepository.AssertCalled(t, "FindByEmail", mock.Anything, mock.Anything) // will call find by email and find existing user
		mockUserRepository.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)   // will not be continued to Create, because existing user found
		defer mockUserRepository.AssertExpectations(t)
	})
}

func TestUserUsecase_Login(t *testing.T) {
	t.Run("test login wrong credential", func(t *testing.T) {
		mockUserRepository := new(mocks.UserRepository)
		userUsecase := NewUserUsecase(mockUserRepository)
		loginDTO := &domain.LoginDTO{
			Email:    "testtest@gmail.com",
			Password: "Test*999",
		}
		mockUserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
		accessToken, err := userUsecase.Login(context.Background(), loginDTO)
		assert.Error(t, err)
		assert.Empty(t, accessToken)
		mockUserRepository.AssertCalled(t, "FindByEmail", mock.Anything, mock.Anything)
		defer mockUserRepository.AssertExpectations(t)
	})
	t.Run("test login correct credential", func(t *testing.T) {
		mockUserRepository := new(mocks.UserRepository)
		userUsecase := NewUserUsecase(mockUserRepository)
		loginDTO := &domain.LoginDTO{
			Email:    "testtest@gmail.com",
			Password: "Test*999",
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginDTO.Password), bcrypt.DefaultCost)
		assert.NoError(t, err)
		mockUserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&entity.User{
			UserId:       1,
			Role:         domain.ADMIN,
			Email:        loginDTO.Email,
			Password:     string(hashedPassword),
			Name:         "Test User",
			ProfileImage: "Google Cloud Storage Path",
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		}, nil)
		accessToken, err := userUsecase.Login(context.Background(), loginDTO)
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		mockUserRepository.AssertCalled(t, "FindByEmail", mock.Anything, mock.Anything)
		defer mockUserRepository.AssertExpectations(t)
	})
}
