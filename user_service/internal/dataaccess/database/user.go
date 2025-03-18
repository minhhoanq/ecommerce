package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/common/logger"
)

type UserDataAccessor interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (*User, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (*VerifyEmail, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (*Session, error)
	GetSessionByUserId(ctx context.Context, user_id uuid.UUID) (*Session, error)
}

type userDataAccessor struct {
	database Database
	l        logger.Interface
}

func NewUserDataAccessor(database Database, l logger.Interface) UserDataAccessor {
	return &userDataAccessor{
		database: database,
		l:        l,
	}
}

func (u *userDataAccessor) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	fmt.Println("[INFO] UserRepo - GetUser - User ID", id, time.Now())
	var user *User
	result := u.database.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		fmt.Println("[ERROR] UserRepo - GetUser - User ID", result.Error)
		return nil, result.Error
	}

	return user, nil
}

func (u *userDataAccessor) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	fmt.Println("[INFO] UserRepo - GetUser - UserName", username)
	var user *User
	result := u.database.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u *userDataAccessor) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	fmt.Println("[INFO] UserRepo - CreateUser - User", arg)

	var user *User = &User{
		Username: arg.Username,
		Email:    arg.Email,
		Password: arg.Password,
		RoleId:   arg.RoleId,
	}

	result := u.database.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u *userDataAccessor) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (*VerifyEmail, error) {
	var verifyEmail *VerifyEmail = &VerifyEmail{
		UserId:     arg.UserId,
		Email:      arg.Email,
		SecretCode: arg.SecretCode,
	}
	result := u.database.WithContext(ctx).Create(&verifyEmail)
	if result.Error != nil {
		return nil, result.Error
	}

	return verifyEmail, nil
}

func (u *userDataAccessor) CreateSession(ctx context.Context, arg CreateSessionParams) (*Session, error) {
	var session *Session = &Session{
		UserId:       arg.UserId,
		RefreshToken: arg.RefreshToken,
		UserAgent:    arg.UserAgent,
		ClientIp:     arg.ClientIp,
		IsBlocked:    arg.IsBlocked,
		ExpiredAt:    arg.ExpiredAt,
	}

	result := u.database.WithContext(ctx).Create(&session)

	if result.Error != nil {
		return nil, result.Error
	}

	return session, nil
}

func (u *userDataAccessor) GetSessionByUserId(ctx context.Context, user_id uuid.UUID) (*Session, error) {
	var session *Session
	result := u.database.WithContext(ctx).Where("user_id = ?", user_id).First(&session)
	fmt.Println("result: ", session.UserId)
	if result.Error != nil {
		return nil, result.Error
	}

	return session, nil
}

func (u *userDataAccessor) execTx(ctx context.Context, fn func(UserDataAccessor) error) error {
	// Bắt đầu transaction với GORM
	tx := u.database.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// Thực thi function với transaction
	err := fn(u)
	if err != nil {
		// Nếu có lỗi, rollback transaction
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}
