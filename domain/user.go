package domain

import (
	"context"

	"github.com/victorsantoso/endeus/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (string, int64, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindById(ctx context.Context, userId int64) (*entity.User, error)
}

type UserUsecase interface {
	Register(ctx context.Context, registerDTO *RegisterDTO) (string, error)
	Login(ctx context.Context, loginDTO *LoginDTO) (string, error)
}

const (
	ADMIN  string = "ADMIN"
	READER string = "READER"
)

type RegisterDTO struct {
	Role         string `json:"role" binding:"required,oneof=ADMIN READER"`
	Email        string `json:"email" binding:"required,email,min=9,max=60"`
	Password     string `json:"password" binding:"required,min=6,max=20"`
	Name         string `json:"name" binding:"required,min=3,max=60"`
	ProfileImage string `json:"profile_image"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	Message     string `json:"message"`
	Code        int    `json:"code"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	Message     string `json:"message"`
	Code        int    `json:"code"`
}
