package database

import (
	"context"
	"testing"

	"github.com/minhhoanq/ecommerce/user_service/internal/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *User {
	hashedPassword, err := util.HashPassword(util.RamdomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username: util.RandomUsername(),
		Email:    util.RandomEmail(),
		Password: hashedPassword,
		RoleId:   1,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.RoleId, user.RoleId)

	require.True(t, user.PasswordChangeAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// func TestGetUser
