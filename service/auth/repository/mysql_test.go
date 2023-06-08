package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/miiy/goc/service/auth/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func newMockDb() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := gorm.Open(mysql.Dialector{Config: &mysql.Config{DriverName: "mysql", Conn: db, SkipInitializeWithVersion: true}}, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, err
}

func TestMysqlAuthRepository_Create(t *testing.T) {
	db, mock, err := newMockDb()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewAuthRepository(db)
	err = repo.Create(context.Background(), &entity.User{
		Username:          "test",
		Password:          "123456",
		Email:             "test@test.com",
		EmailVerifiedTime: "",
		Phone:             "",
		Status:            0,
	})
	if err != nil {
		t.Error(err)
	}

}
