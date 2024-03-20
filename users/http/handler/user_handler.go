package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorsantoso/endeus/domain"
)

type userHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(g *gin.Engine, userUsecase domain.UserUsecase) {
	userHandler := &userHandler{
		userUsecase: userUsecase,
	}

	userGroup := g.Group("/api/v1")
	userGroup.POST("/register", userHandler.Register)
	userGroup.POST("/login", userHandler.Login)
}

func (uh *userHandler) Register(c *gin.Context) {
	registerDTO := &domain.RegisterDTO{}
	err := c.ShouldBindJSON(registerDTO)
	// handle validation error
	if err != nil {
		c.JSON(http.StatusBadRequest, &domain.RegisterResponse{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// handle role validation
	if registerDTO.Role != domain.ADMIN && registerDTO.Role != domain.READER {
		c.JSON(http.StatusBadRequest, &domain.RegisterResponse{
			Message: domain.ErrBadRequest.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// process registration and return access token
	accessToken, err := uh.userUsecase.Register(context.Background(), registerDTO)
	if err != nil {
		if err == domain.ErrDuplicateUser {
			c.JSON(http.StatusConflict, &domain.RegisterResponse{
				Message: domain.ErrDuplicateUser.Error(),
				Code:    http.StatusConflict,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &domain.RegisterResponse{
			Message: domain.ErrInternalServerError.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// return access token upon successful registration process
	c.JSON(http.StatusOK, &domain.RegisterResponse{
		AccessToken: accessToken,
		Message:     "successfully registered a new user",
		Code:        http.StatusOK,
	})
}

func (uh *userHandler) Login(c *gin.Context) {
	loginDTO := &domain.LoginDTO{}
	// handle validation error
	if err := c.ShouldBindJSON(loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, &domain.LoginResponse{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// validate user login process
	accessToken, err := uh.userUsecase.Login(context.Background(), loginDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, &domain.LoginResponse{
			Message: domain.ErrInvalidCredential.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// return access token upon successful login process
	c.JSON(http.StatusOK, &domain.LoginResponse{
		AccessToken: accessToken,
		Message:     "successfully logged in.",
		Code:        http.StatusOK,
	})
}
