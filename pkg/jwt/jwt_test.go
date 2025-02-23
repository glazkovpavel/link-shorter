package jwt_test

import (
	"go/link-shorter/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {
	const email = "test@example.com"
	jwtService := jwt.NewJWT("NMLHybLCGsFhNPWUR5hN84AFlGt0JaKz-Epk3vVAg2YVBEPNp-9GRlwpukKG6J230UfrRW6VMw3y87xfkgs7Bg")
	token, err := jwtService.CreateToken(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	isValid, data := jwtService.ParseToken(token)
	if !isValid {
		t.Fatal("token is not valid")
	}
	if data.Email != email {
		t.Fatal("email is not valid")
	}

}
