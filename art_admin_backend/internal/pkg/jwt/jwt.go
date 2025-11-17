package jwt

import (
	"errors"
	"time"

	"art_admin_backend/internal/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 声明
type Claims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateToken 生成访问令牌
func GenerateToken(userID int64) (string, error) {
	cfg := config.GlobalConfig.JWT

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.GetAccessTokenDuration())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID int64) (string, error) {
	cfg := config.GlobalConfig.JWT

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.GetRefreshTokenDuration())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析令牌
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GlobalConfig.JWT

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// RefreshAccessToken 刷新访问令牌
func RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 生成新的访问令牌
	return GenerateToken(claims.UserID)
}

