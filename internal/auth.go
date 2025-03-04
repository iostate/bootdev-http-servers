package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const TokenTypeAccess TokenType = "chirpy"

// HashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hashedPw), nil
}

// CheckPasswordHash compares a plain password with its hashed version.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// these are standardized JWT claims RFC 7519
	// RFC7519 link: https://datatracker.ietf.org/doc/html/rfc7519#section-4.1
	claims := &jwt.RegisteredClaims{
		Issuer: string(TokenTypeAccess),
		// Use UTC in APIs & distributed systems,
		// to avoid ambiguity with interpreting time
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// do not assume this will convert to byte automatically even though in jwt/v5 it supposedly does.
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	// Parse the claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Error getting Authorization header. Authorization header is empty")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("Unexpected authorization format")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return "", fmt.Errorf("empty token in authorization header")
	}

	return tokenString, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Error getting Authorization header. Authorization header is empty")
	}

	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", fmt.Errorf("Unexpected authorization format")
	}

	tokenString := strings.TrimPrefix(authHeader, "ApiKey ")
	if tokenString == "" {
		return "", fmt.Errorf("empty token in authorization header")
	}

	return tokenString, nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("Error reading random bytes into token")
	}
	return hex.EncodeToString(token), nil
}
