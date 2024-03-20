package usecase

import (
	"context"
	"database/sql"

	"github.com/apex/log"
	"github.com/victorsantoso/endeus/domain"
	"github.com/victorsantoso/endeus/entity"
	"github.com/victorsantoso/endeus/helper"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(userRepository domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (uu *userUsecase) Register(ctx context.Context, registerDTO *domain.RegisterDTO) (string, error) {
	// validation for register user
	validateUser, err := uu.userRepository.FindByEmail(ctx, registerDTO.Email)
	if validateUser != nil {
		log.Debugf("[user_usecase.Register] user already registered with: %+v", validateUser)
		return "", domain.ErrDuplicateUser
	}
	if err == sql.ErrNoRows {
		log.Debugf("[user_usecase.Register] error finding user by email, err: %v", err)
	}

	// hash password for security purposes
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("[user_usecase.Register] error generating password hash, err: %v", err)
		return "", err
	}

	// create a new user
	role, userId, err := uu.userRepository.Create(ctx, &entity.User{
		Role:         registerDTO.Role,
		Email:        registerDTO.Email,
		Password:     string(hashedPassword),
		Name:         registerDTO.Name,
		ProfileImage: registerDTO.ProfileImage,
	})
	if err != nil {
		if err == domain.ErrInvalidRole {
			log.Debugf("[user_usecase.Register] invalid role, err: %v", err)
		}
		if err == domain.ErrDuplicateUser {
			log.Debugf("[user_usecase.Register] duplicate entry, err: %v", err)
		}
		log.Errorf("[user_usecase.Register] error creating a new user, err: %v", err)
		return "", err
	}

	// generate access token
	accessToken, err := helper.GenerateJWT(role, userId)
	if err != nil {
		log.Errorf("[user_usecase.Register] error generating access token, err: %v", err)
		return "", err // Return the error here instead of nil
	}
	return accessToken, nil
}

func (uu *userUsecase) Login(ctx context.Context, loginDTO *domain.LoginDTO) (string, error) {
	// validate user existence
	user, err := uu.userRepository.FindByEmail(ctx, loginDTO.Email)
	if user == nil || err != nil {
		log.Debugf("[user_usecase.Login] failed to find user with email: %s, err: %v", loginDTO.Email, err)
		return "", err // Return the error here
	}
	// validate user password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDTO.Password))
	if err != nil {
		log.Debugf("[user_usecase.Login] failed on compare hash and password process, err: %v", err)
		return "", err
	}
	// generate jwt if the credential is valid
	accessToken, err := helper.GenerateJWT(user.Role, user.UserId)
	if err != nil {
		log.Errorf("[user_usecase.Login] failed on generating jwt process: %v", err)
		return "", err // Return the error here instead of nil
	}
	return accessToken, nil
}