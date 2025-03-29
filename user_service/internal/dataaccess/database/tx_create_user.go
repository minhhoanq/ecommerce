package database

import (
	"context"

	"github.com/minhhoanq/ecommerce/user_service/internal/util"
)

func (q *userDataAccessor) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var user *User
	var result CreateUserTxResult

	err := q.execTx(ctx, func(q UserDataAccessor) error {
		var err error

		user, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		result.User = *user

		payloadVerifyEmail := CreateVerifyEmailParams{
			UserId:     user.ID,
			Email:      user.Email,
			SecretCode: util.RandomString(32),
		}

		verifyEmail, err := q.CreateVerifyEmail(ctx, payloadVerifyEmail)
		if err != nil {
			return err
		}

		message := UserCreatedParams{
			UserID:        user.ID,
			Email:         user.Email,
			Username:      user.Username,
			SecretCode:    verifyEmail.SecretCode,
			VerifyEmailID: verifyEmail.ID,
		}

		return arg.AfterCreate(&message)
	})

	return result, err
}
