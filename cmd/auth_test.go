package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"go/link-shorter/internal/auth"
	"go/link-shorter/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "test@test.com",
		Password: "$2a$10$L.4G92r4o6jNZBvqJTXJEuwBo1JJTVxPGQu/sQTZDbRQ3qqtOFrq2",
		Name:     "Vasya",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().Where("email =?", "test@test.com").Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	// подготовка
	db := initDb()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@test.com",
		Password: "password",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("got %d, want %d", res.StatusCode, http.StatusOK)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var resp auth.LoginResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Token == "" {
		t.Errorf("got empty token")
	}
	removeData(db)
}

func TestLoginFail(t *testing.T) {
	// подготовка
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@example.com",
		Password: "noPassword",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", res.StatusCode, http.StatusUnauthorized)
	}
	removeData(db)
}
