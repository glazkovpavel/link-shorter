package auth_test

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"go/link-shorter/configs"
	"go/link-shorter/internal/auth"
	"go/link-shorter/internal/user"
	"go/link-shorter/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}

	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})

	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
	return &handler, mock, nil
}

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password", "name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	if err != nil {
		t.Fatal(err)
		return
	}
	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "test@test.com",
		Password: "password",
		Name:     "Vasya",
	})

	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", reader)
	handler.Register()(wr, req)
	if wr.Code != http.StatusCreated {
		t.Errorf("login: expected status %d, got %d", http.StatusCreated, wr.Code)
	}
}

func TestLoginHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("test@test.com", "$2a$10$L.4G92r4o6jNZBvqJTXJEuwBo1JJTVxPGQu/sQTZDbRQ3qqtOFrq2")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	if err != nil {
		t.Fatal(err)
		return
	}
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@test.com",
		Password: "password",
	})

	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", reader)
	handler.Login()(wr, req)
	if wr.Code != http.StatusOK {
		t.Errorf("login: expected status %d, got %d", http.StatusOK, wr.Code)
	}
}
