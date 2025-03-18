package database

import (
	"context"
	"testing"
	"time"

	"github.com/minhhoanq/ecommerce/user_service/internal/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
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

func TestGetUserByID(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUserByID(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)

	require.WithinDuration(t, user1.PasswordChangeAt, user2.PasswordChangeAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestGetUserByUsername(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)

	require.WithinDuration(t, user1.PasswordChangeAt, user2.PasswordChangeAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestCreateVerifyEmail(t *testing.T) {
	user := createRandomUser(t)

	arg := CreateVerifyEmailParams{
		UserId:     user.ID,
		Email:      user.Email,
		SecretCode: util.RandomString(8),
	}

	verifyEmail, err := testQueries.CreateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, verifyEmail)

	require.Equal(t, verifyEmail.SecretCode, arg.SecretCode)
}

func TestCreateSession(t *testing.T) {
	user := createRandomUser(t)
	arg := CreateSessionParams{
		UserId:       user.ID,
		RefreshToken: util.RandomString(10),
		UserAgent:    util.RandomString(6),
		ClientIp:     util.RandomString(6),
		IsBlocked:    false,
		ExpiredAt:    time.Now(),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, session.UserId, arg.UserId)
	require.Equal(t, session.RefreshToken, arg.RefreshToken)
	require.Equal(t, session.UserAgent, arg.UserAgent)
	require.Equal(t, session.ClientIp, arg.ClientIp)
	require.Equal(t, session.IsBlocked, arg.IsBlocked)
	require.WithinDuration(t, session.ExpiredAt, arg.ExpiredAt, time.Second)
}

func createRandomSession(t *testing.T) (*Session, *User) {
	user := createRandomUser(t)
	arg := CreateSessionParams{
		UserId:       user.ID,
		RefreshToken: util.RandomString(10),
		UserAgent:    util.RandomString(6),
		ClientIp:     util.RandomString(6),
		IsBlocked:    false,
		ExpiredAt:    time.Now(),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, session.UserId, arg.UserId)
	require.Equal(t, session.RefreshToken, arg.RefreshToken)
	require.Equal(t, session.UserAgent, arg.UserAgent)
	require.Equal(t, session.ClientIp, arg.ClientIp)
	require.Equal(t, session.IsBlocked, arg.IsBlocked)
	require.WithinDuration(t, session.ExpiredAt, arg.ExpiredAt, time.Second)

	return session, user
}

func TestGetSessionByUserId(t *testing.T) {
	session, user := createRandomSession(t)

	session2, err := testQueries.GetSessionByUserId(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session.ID, session2.ID)
	require.Equal(t, session.UserId, session2.UserId)
	require.Equal(t, session.RefreshToken, session2.RefreshToken)
	require.Equal(t, session.UserAgent, session2.UserAgent)
	require.Equal(t, session.ClientIp, session2.ClientIp)
	require.Equal(t, session.IsBlocked, session2.IsBlocked)
	require.WithinDuration(t, session.ExpiredAt, session2.ExpiredAt, time.Second)
}
