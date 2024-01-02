package db

import (
	"context"
	"log"
	"testing"

	"github.com/broemp/red_card/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	hashed_password, err := util.HashPassword(util.RandomString(12))
	if err != nil {
		log.Fatal("could not hash password:", err)
	}

	arg := CreateUserParams{
		Username:       util.RandomUsername(),
		HashedPassword: hashed_password,
		Name:           util.RandomString(6),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Name, user.Name)
	require.NotZero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)
	require.False(t, user.DeletedAt.Valid)
}
