package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/minhhoanq/ecommerce/user_service/internal/entity"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// GetUserByID -.
func (q *SQLStore) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	fmt.Println("[INFO] UserRepo - GetUser - User ID", id, time.Now())
	var user *entity.User
	result := q.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		fmt.Println("[ERROR] UserRepo - GetUser - User ID", result.Error)
		return nil, result.Error
	}

	return user, nil
}

// GetUserByUsername -.
func (q *SQLStore) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	fmt.Println("[INFO] UserRepo - GetUser - UserName", username)
	var user *entity.User
	result := q.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	RoleId   int    `json:"role_id"`
}

// CreateUser -.
func (q *SQLStore) CreateUser(ctx context.Context, arg CreateUserParams) (*entity.User, error) {
	fmt.Println("[INFO] UserRepo - CreateUser - User", arg)

	var user *entity.User = &entity.User{
		Username: arg.Username,
		Email:    arg.Email,
		Password: arg.Password,
		RoleId:   arg.RoleId,
	}

	result := q.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
