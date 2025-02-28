// package auth

// import (
// 	"testing"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/google/uuid"
// )

// var claims *jwt.RegisteredClaims
// var tokenSecret string
// var userID uuid.UUID

// func TestMain(m *testing.M) {

// 	testCases := []struct{
// 		claims 			*jwt.RegisteredClaims
// 		tokenSecret		string
// 		userID			uuid.UUID
// 	} {
// 		{claims, "test123", userID},
// 	}

// 	for _, tC := range testCases {
// 		duration := time.Until(time.Now().Add(5000))
// 		generateToken, _ := MakeJWT(tC.userID, tC.tokenSecret, duration)
// 		ValidateJWT(generateToken, tC.tokenSecret)
// 	}
// 	m.Run()
// }

// func TestSetupClaims(t *testing.T) {
// 	claims = &jwt.RegisteredClaims{
// 		Issuer: "chirpy",
// 		IssuedAt: jwt.NewNumericDate(time.Now()),
// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(50000)),
// 		Subject: userID.String(),
// 	}

// }

// func TestValidateJWT(t *testing.T) {
// 	userID = uuid.New()
// 	validToken, _ := MakeJWT(userID, "testing", time.Hour)

// 	tests := []struct{
// 		name		string
// 		tokenString	string
// 		tokenSecret	string
// 		wantUserID	uuid.UUID
// 		wantErr		bool
// 	} {
// 		{
// 			name: "Valid JWT Token",
// 			tokenString: validToken,
// 			tokenSecret: "testing",
// 			wantUserID: userID,
// 			wantErr: false,
// 		},
// 		{
// 			name: "Invalid token",
// 			tokenString: "Invalid.token.string",
// 			tokenSecret: "testing",
// 			wantUserID: userID,
// 			wantErr: true,
// 		},
// 		{
// 			name: "Invalid secret",
// 			tokenString: validToken,
// 			tokenSecret: "blah",
// 			wantUserID: uuid.Nil,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ValidateJWT() error = %v, wantErr = %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotUserID != tt.wantUserID {
// 				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
// 			}
// 		})
// 	}

// }

// func TestCheckPasswordHash(t *testing.T) {
// 	password1 := "test123#"
// 	password2 := "wrongPassword@"
// 	hash1, _ := HashPassword(password1)
// 	hash2, _ := HashPassword(password2)

// 	tests := []struct{
// 		name		string
// 		password	string
// 		hash		string
// 		wantErr		bool
// 	} {
// 		{
// 			name: 		"Correct password",
// 			password: 	password1,
// 			hash: 		hash1,
// 			wantErr: 	false,
// 		},
// 		{
// 			name: 		"Incorrect password",
// 			password: 	"wrongPW",
// 			hash: 		hash1,
// 			wantErr: 	true,
// 		},
// 		{
// 			name: 		"Password doesn't match hash",
// 			password: 	password1,
// 			hash: 		hash2,
// 			wantErr: 	true,
// 		},
// 		{
// 			name: 		"Empty password",
// 			password: 	"",
// 			hash: 		hash1,
// 			wantErr: 	true,
// 		},
// 		{
// 			name: 		"Invalid hash",
// 			password: 	password1,
// 			hash: 		hash1,
// 			wantErr: 	true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := CheckPasswordHash(tt.password, tt.hash)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("CheckPasswordHash() error = %v, wantErr =  %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Logf("Test Name: %v", tt.name)
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}