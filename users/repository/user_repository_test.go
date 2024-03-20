package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/victorsantoso/endeus/domain"
	"github.com/victorsantoso/endeus/entity"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	userRepository := NewUserRepository(db)

	type args struct {
		role         string
		email        string
		password     string
		name         string
		profileImage string
	}
	tests := []struct {
		name         string
		args         args
		testFunction func(t *testing.T, tt args)
	}{
		{
			name: "test create a new user",
			args: args{
				role:         "ADMIN",
				email:        "testtest@gmail.com",
				password:     "Test*123",
				name:         "Test User",
				profileImage: "googlecloudstorageurl", // for demo application purpose only
			},
			testFunction: func(t *testing.T, tt args) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tt.password), bcrypt.DefaultCost)
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				user := &entity.User{
					Role:         tt.role,
					Email:        tt.email,
					Password:     string(hashedPassword),
					Name:         tt.name,
					ProfileImage: tt.profileImage,
				}
				mock.ExpectBegin()                                                       // expect begin transaction
				rows := sqlmock.NewRows([]string{"role", "user_id"}).AddRow("ADMIN", 1)  // should return role "ADMIN" and user_id 1 since all datas are valid
				mock.ExpectQuery(regexp.QuoteMeta(CreateUserQuery)).WillReturnRows(rows) // expect query will return rows
				mock.ExpectCommit()                                                      // expect commit transaction because of no error
				role, userId, err := userRepository.Create(context.Background(), user)   // create a new user
				assert.NoError(t, err)                                                   // there should be no error on creating a new data
				assert.Equal(t, domain.ADMIN, role)                                      // test the role should be the same
				assert.Equal(t, int64(1), userId)                                        // test with userId 1 if created a new user
				assert.NoError(t, mock.ExpectationsWereMet())                            // expectations should be fulfilled
			},
		},
		{
			name: "test create a new user with invalid role",
			args: args{
				role:         "asdf", // here we'll do a user creation with invalid role
				email:        "testtest@gmail.com",
				password:     "Test*123",
				name:         "Test User",
				profileImage: "googlecloudstorageurl", // for demo application purpose only
			},
			testFunction: func(t *testing.T, tt args) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tt.password), bcrypt.DefaultCost)
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				user := &entity.User{
					Role:         tt.role,
					Email:        tt.email,
					Password:     string(hashedPassword),
					Name:         tt.name,
					ProfileImage: tt.profileImage,
				}
				mock.ExpectBegin()                                                       // expect begin transaction
				rows := sqlmock.NewRows([]string{"role", "user_id"})                     // add empty result
				mock.ExpectQuery(regexp.QuoteMeta(CreateUserQuery)).WillReturnRows(rows) // expect there will be no rows returned
				mock.ExpectRollback()                                                    // expect rollback due to error on role validation
				role, userId, err := userRepository.Create(context.Background(), user)   // here we attempted to create a new user with invalid data
				assert.Error(t, err)                                                     // there should be an error on creating user
				assert.Empty(t, role)                                                    // the role should be empty
				assert.Empty(t, userId)                                                  // the userId should be empty
				assert.NoError(t, mock.ExpectationsWereMet())                            // expectations should be fulfilled
			},
		},
		{
			name: "test create a new user with duplicate data",
			args: args{
				role:         "ADMIN",
				email:        "testtest@gmail.com",
				password:     "Test*123",
				name:         "Test User",
				profileImage: "googlecloudstorageurl", // for demo application purpose only
			},
			testFunction: func(t *testing.T, tt args) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tt.password), bcrypt.DefaultCost)
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				user := &entity.User{
					Role:         tt.role,
					Email:        tt.email,
					Password:     string(hashedPassword),
					Name:         tt.name,
					ProfileImage: tt.profileImage,
				}
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"role", "user_id"}).AddRow("ADMIN", 1)
				mock.ExpectQuery(regexp.QuoteMeta(CreateUserQuery)).WillReturnRows(rows) // should return role "ADMIN" and user_id 1
				mock.ExpectCommit()                                                      // here we expect commit due to no error
				role, userId, err := userRepository.Create(context.Background(), user)   // here we attempted to create a new user
				assert.NoError(t, err)                                                   // there should be no error on creating a new data
				assert.Equal(t, domain.ADMIN, role)                                      // test the role should be the same
				assert.Equal(t, int64(1), userId)                                        // test with userId 1 if created a new user
				assert.NoError(t, mock.ExpectationsWereMet())                            // first expectatios should met
				mock.ExpectBegin()                                                       // expect begin
				rows2 := sqlmock.NewRows([]string{"role", "user_id"})
				mock.ExpectQuery(regexp.QuoteMeta(CreateUserQuery)).WillReturnRows(rows2) // should return empty data because of conflicted unique email
				mock.ExpectRollback()                                                     // expect rollback due to error
				role2, userId2, err2 := userRepository.Create(context.Background(), user) // here we attempted to create the same user
				assert.Error(t, err2)                                                     // should return error
				assert.Empty(t, role2)                                                    // role2 should be empty
				assert.Empty(t, userId2)                                                  // userId2 should be empty
				assert.NoError(t, mock.ExpectationsWereMet())                             // expectations should be fulfilled
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunction(t, tt.args)
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	userRepository := NewUserRepository(db)

	type args struct {
		email string
	}
	tests := []struct {
		name         string
		args         args
		testFunction func(t *testing.T, tt args)
	}{
		{
			name: "test find user by email from empty table",
			args: args{
				email: "testtest@gmail.com",
			},
			testFunction: func(t *testing.T, tt args) {
				rows := sqlmock.NewRows([]string{"user_id", "role", "email", "password", "name", "profile_image", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta(FindByEmailQuery)).WillReturnRows(rows)
				user, err := userRepository.FindByEmail(context.Background(), tt.email)
				assert.Error(t, err)
				assert.Empty(t, user)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			name: "test find user by email from existing data in table",
			args: args{
				email: "testtest@gmail.com",
			},
			testFunction: func(t *testing.T, tt args) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Test*123"), bcrypt.DefaultCost)
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				testUser := &entity.User{
					UserId:       1,
					Role:         "ADMIN",
					Email:        tt.email,
					Password:     string(hashedPassword),
					Name:         "Test User",
					ProfileImage: "sample google cloud storage url from another service",
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				}
				rows := sqlmock.NewRows([]string{"user_id", "role", "email", "password", "name", "profile_image", "created_at", "updated_at"})
				rows.AddRow(testUser.UserId, testUser.Role, testUser.Email, hashedPassword, testUser.Name, testUser.ProfileImage, testUser.CreatedAt, testUser.UpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(FindByEmailQuery)).WillReturnRows(rows) // expect query will return rows
				user, err := userRepository.FindByEmail(context.Background(), tt.email)
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunction(t, tt.args)
		})
	}
}

func TestUserRepository_FindById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	userRepository := NewUserRepository(db)

	type args struct {
		userId int64
	}
	tests := []struct {
		name         string
		args         args
		testFunction func(t *testing.T, tt args)
	}{
		{
			name: "test find user by id from empty table",
			args: args{
				userId: 1,
			},
			testFunction: func(t *testing.T, tt args) {
				rows := sqlmock.NewRows([]string{"user_id", "role", "email", "password", "name", "profile_image", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta(FindByIdQuery)).WillReturnRows(rows)
				user, err := userRepository.FindById(context.Background(), tt.userId)
				assert.Error(t, err)
				assert.Empty(t, user)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			name: "test find user by id from existing data in table",
			args: args{
				userId: 1,
			},
			testFunction: func(t *testing.T, tt args) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Test*123"), bcrypt.DefaultCost)
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				testUser := &entity.User{
					UserId:       1,
					Role:         "ADMIN",
					Email:        "testtest@gmail.com",
					Password:     string(hashedPassword),
					Name:         "Test User",
					ProfileImage: "sample google cloud storage url from another service",
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				}
				rows := sqlmock.NewRows([]string{"user_id", "role", "email", "password", "name", "profile_image", "created_at", "updated_at"})
				rows.AddRow(testUser.UserId, testUser.Role, testUser.Email, hashedPassword, testUser.Name, testUser.ProfileImage, testUser.CreatedAt, testUser.UpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(FindByIdQuery)).WillReturnRows(rows) // expect query will return rows
				user, err := userRepository.FindById(context.Background(), tt.userId)
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunction(t, tt.args)
		})
	}
}
