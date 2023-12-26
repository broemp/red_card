package token

import (
	"errors"
	"fmt"
	"time"

	db "github.com/broemp/red_card/db/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("token is invalid")

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

type Claims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (m *JWTMaker) CreateToken(user db.User, duration time.Duration) (string, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	claims := Claims{
		Username: user.Username,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			ID:        tokenID.String(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString([]byte(m.secretKey))
}

func (m *JWTMaker) VerifyToken(token string) (*Claims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(m.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
