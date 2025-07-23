package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey []byte
}

type Claims struct {
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	RoleID    int    `json:"role_id"`
	EmpresaID int    `json:"empresa_id"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken gera um JWT token
func (s *JWTService) GenerateToken(userID int, email string, roleID, empresaID int) (string, error) {
	claims := Claims{
		UserID:    userID,
		Email:     email,
		RoleID:    roleID,
		EmpresaID: empresaID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "trabiju-telemetria",
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// GenerateRefreshToken gera um refresh token (válido por 7 dias)
func (s *JWTService) GenerateRefreshToken(userID int, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 dias
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "trabiju-telemetria-refresh",
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken valida e retorna as claims do token
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inválido")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}

// RefreshToken gera um novo token a partir de um refresh token válido
func (s *JWTService) RefreshToken(refreshTokenString string, roleID, empresaID int) (string, error) {
	claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	// Verificar se é um refresh token
	if claims.Issuer != "trabiju-telemetria-refresh" {
		return "", errors.New("token não é um refresh token")
	}

	// Gerar novo access token
	return s.GenerateToken(claims.UserID, claims.Email, roleID, empresaID)
}
