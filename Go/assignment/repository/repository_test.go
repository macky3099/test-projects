package repository

import (
	"codeassign/models"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestCreateUser_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewUserRepository(gormDB)
	user := &models.UserDetails{
		Name:   "John Doe",
		PAN:    "ABCDE1234F",
		Mobile: 9876543210,
		Email:  "john.doe@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_details`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "John Doe", "ABCDE1234F", 9876543210, "john.doe@example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateUser(user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewUserRepository(gormDB)

	user := &models.UserDetails{
		Name:   "John Doe",
		PAN:    "ABCDE1234F",
		Mobile: 9876543210,
		Email:  "john.doe@example.com",
	}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_details`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "John Doe", "ABCDE1234F", 9876543210, "john.doe@example.com").
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	err := repo.CreateUser(user)
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewUserRepository(t *testing.T) {
	gormDB, _ := setupMockDB(t)
	repo := NewUserRepository(gormDB)
	assert.NotNil(t, repo)
	_, ok := repo.(*GormUserRepository)
	assert.True(t, ok)
}
