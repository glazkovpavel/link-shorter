package auth_test

import (
	"go/link-shorter/internal/auth"
	"go/link-shorter/internal/user"
	"testing"
)

type MockUserRepository struct{}

func (m *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "user@example.com",
	}, nil
}

func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const initialEmail = "user@example.com"
	authService := auth.NewAuthService(&MockUserRepository{})
	email, err := authService.Register(initialEmail, "1", "Dasha")
	if err != nil {
		t.Fatal(err)
	}
	if email != initialEmail {
		t.Fatal("Email does not match")
	}

}
