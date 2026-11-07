package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims menyimpan payload klaim kustom untuk token JWT.
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken membuat token JWT baru yang ditandatangani dengan algoritma HS256.
func GenerateToken(userID, email, secret string, expirationHours int) (string, error) {
	if secret == "" {
		return "", errors.New("JWT secret tidak boleh kosong")
	}
	if expirationHours <= 0 {
		expirationHours = 24
	}

	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expirationHours) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken memverifikasi signature dan masa berlaku token JWT, lalu mengembalikan klaimnya.
func ValidateToken(tokenString, secret string) (*JWTClaims, error) {
	if secret == "" {
		return nil, errors.New("JWT secret tidak boleh kosong")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing token tidak valid")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token JWT tidak valid")
	}

	return claims, nil
}
