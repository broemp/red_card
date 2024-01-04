package token

import (
	"fmt"
	"testing"
	"time"

	db "github.com/broemp/red_card/db/sqlc"
	"github.com/broemp/red_card/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	testUser := createTestUser()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, claims, err := maker.CreateToken(testUser, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, claims)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, testUser.Username, payload.Username)
	require.Equal(t, fmt.Sprint(testUser.ID), payload.Subject)
	require.WithinDuration(t, issuedAt, payload.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt.Time, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, claims, err := maker.CreateToken(createTestUser(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, claims)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	claims := Claims{
		Username: util.RandomUsername(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}

func createTestUser() db.User {
	return db.User{
		ID:       util.RandomInt(0, 25680),
		Username: util.RandomUsername(),
	}
}
