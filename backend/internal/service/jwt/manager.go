package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims define los claims personalizados del JWT
type CustomClaims struct {
	UserID       string   `json:"sub"`
	Email        string   `json:"email"`
	RoleID       string   `json:"role_id"`
	Permissions  []string `json:"permissions"`
	IsSuperAdmin bool     `json:"is_super_admin"`
	JTI          string   `json:"jti"` // JWT ID para revocación
	jwt.RegisteredClaims
}

// Manager gestiona operaciones JWT
type Manager struct {
	secret            string
	expiration        int
	refreshExpiration int
}

// NewManager crea una nueva instancia del manager JWT
func NewManager(secret string, expiration, refreshExpiration int) *Manager {
	return &Manager{
		secret:            secret,
		expiration:        expiration,
		refreshExpiration: refreshExpiration,
	}
}

// GenerateToken genera un access token
func (m *Manager) GenerateToken(userID, email, roleID string, permissions []string, isSuperAdmin bool) (string, int, error) {
	now := time.Now()
	expirationTime := now.Add(time.Duration(m.expiration) * time.Second)

	claims := CustomClaims{
		UserID:       userID,
		Email:        email,
		RoleID:       roleID,
		Permissions:  permissions,
		IsSuperAdmin: isSuperAdmin,
		JTI:          generateJTI(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "psicosst-cloud",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", 0, fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, m.expiration, nil
}

// GenerateRefreshToken genera un refresh token
func (m *Manager) GenerateRefreshToken(userID string) (string, int, error) {
	now := time.Now()
	expirationTime := now.Add(time.Duration(m.refreshExpiration) * time.Second)

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "psicosst-cloud",
		ID:        generateJTI(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", 0, fmt.Errorf("error signing refresh token: %w", err)
	}

	return tokenString, m.refreshExpiration, nil
}

// ValidateToken valida y parsea un token
func (m *Manager) ValidateToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// generateJTI genera un ID único para el token (para revocación)
func generateJTI() string {
	// En producción, usar UUID
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
