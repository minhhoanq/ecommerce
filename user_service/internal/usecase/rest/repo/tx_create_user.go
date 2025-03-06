package repo

import (
	"context"

	"github.com/minhhoanq/ecommerce/user_service/internal/entity"
)

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user *entity.User) error
}

type CreateUserTxResult struct {
	User *entity.User
}

func (q *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := q.execTx(ctx, func(q Querier) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
