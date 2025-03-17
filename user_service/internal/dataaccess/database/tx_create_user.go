package database

import (
	"context"
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

		return arg.AfterCreate(&result.User)
	})

	return result, err
}
