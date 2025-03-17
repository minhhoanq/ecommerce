package consumers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/user_service/internal/email"
	"github.com/minhhoanq/ecommerce/user_service/internal/handler/producers"
	"github.com/minhhoanq/ecommerce/user_service/internal/util"
	"go.uber.org/zap"
)

type UserSignupHandler interface {
	Handle(ctx context.Context, message producers.UserSignupParams) error
}

type userSignupHandler struct {
	mailer           email.EmailSender
	userDataAccessor database.UserDataAccessor
	l                logger.Interface
}

func NewUserSignupHandler(mailer email.EmailSender, userDataAccessor database.UserDataAccessor, l logger.Interface) UserSignupHandler {
	return &userSignupHandler{
		mailer:           mailer,
		userDataAccessor: userDataAccessor,
		l:                l,
	}
}

func (u userSignupHandler) Handle(ctx context.Context, message producers.UserSignupParams) error {
	u.l.Info("Process user signup event")

	user, err := u.userDataAccessor.GetUserByID(ctx, message.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found %w", asynq.SkipRetry)
		}

		return fmt.Errorf("failed to get user: %w", asynq.SkipRetry)
	}

	verifyEmail, err := u.userDataAccessor.CreateVerifyEmail(ctx, database.CreateVerifyEmailParams{
		UserId:     user.ID,
		Email:      user.Email,
		SecretCode: util.RamdomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to LIFEAT"
	verifyUrl := fmt.Sprintf("http://localhost:8000/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s, <br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">Click here<a> to verify your email address.<br/>`, user.Username, verifyUrl)
	to := []string{user.Email}

	err = u.mailer.SendMail(subject, content, to, nil, nil, nil)

	u.l.Info("process task send verify email kafka", zap.String("email", user.Email))

	return nil
}
