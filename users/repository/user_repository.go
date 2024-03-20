package repository

import (
	"context"
	"database/sql"

	"github.com/apex/log"
	"github.com/lib/pq"
	"github.com/victorsantoso/endeus/domain"
	"github.com/victorsantoso/endeus/entity"
)

type userRepository struct {
	dbConn *sql.DB
}

func NewUserRepository(dbConn *sql.DB) domain.UserRepository {
	return &userRepository{
		dbConn: dbConn,
	}
}

const (
	// query with variables handling to escape sql injections
	CreateUserQuery = `
		INSERT INTO users(role, email, password, name, profile_image, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, now()::timestamptz, now()::timestamptz)
		RETURNING role, user_id;
	`
	FindByEmailQuery = `
		SELECT user_id, role, email, password, name, profile_image, created_at, updated_at FROM users WHERE email = $1;
	`
	FindByIdQuery = `
		SELECT user_id, role, email, password, name, profile_image, created_at, updated_at FROM users WHERE user_id = $1;
	`
)

func (ur *userRepository) Create(ctx context.Context, user *entity.User) (string, int64, error) {
    var userId int64
    var role string
    // begin transaction
    tx, err := ur.dbConn.Begin()
    if err != nil {
        return "", 0, err // Return the error immediately
    }
    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
        tx.Commit()
    }()
    if err := tx.QueryRowContext(ctx, CreateUserQuery, user.Role, user.Email, user.Password, user.Name, user.ProfileImage).Scan(&role, &userId); err != nil {
        if pqError, ok := err.(*pq.Error); ok {
            if pqError.Code == "23505" { // unique violation code on email
                return role, userId, domain.ErrDuplicateUser
            }
        }
        return role, userId, domain.ErrInvalidRole
    }
    log.Debugf("[user_repository] user with user_id: %d, role: %s created", userId, role)
    return role, userId, nil // Return nil error here
}

func (ur *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    var user entity.User
    row := ur.dbConn.QueryRowContext(ctx, FindByEmailQuery, email)
    if err := row.Scan(&user.UserId, &user.Role, &user.Email, &user.Password, &user.Name, &user.ProfileImage, &user.CreatedAt, &user.UpdatedAt); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // Return nil, nil if no rows found
        }
        return nil, err
    }
    return &user, nil
}

func (ur *userRepository) FindById(ctx context.Context, userId int64) (*entity.User, error) {
    var user entity.User
    row := ur.dbConn.QueryRowContext(ctx, FindByIdQuery, userId)
    if err := row.Scan(&user.UserId, &user.Role, &user.Email, &user.Password, &user.Name, &user.ProfileImage, &user.CreatedAt, &user.UpdatedAt); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // Return nil, nil if no rows found
        }
        return nil, err
    }
    return &user, nil
}